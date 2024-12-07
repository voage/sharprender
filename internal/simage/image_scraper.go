package simage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

const (
	defaultTimeout = 10 * time.Second
	evalScript     = `
		Array.from(document.images).map(img => ({
			src: img.src,
			alt: img.alt,
			width: img.width,
			height: img.height
		}))
	`
	scrollScript = `
		async function smoothScroll() {
					const height = document.documentElement.scrollHeight;
					const scrollStep = Math.floor(height / 20);
					
					for (let i = 0; i <= height; i += scrollStep) {
						window.scrollTo({
							top: i,
							behavior: 'smooth'
						});
						await new Promise(resolve => setTimeout(resolve, 300));
					}
				}
				smoothScroll();
	`
	resourceTimingScript = `
	(() => {
		return performance.getEntriesByType('resource')
			.filter(e => e.initiatorType === 'img')
			.map(e => ({
				name: e.name,
				domainLookupStart: e.domainLookupStart,
				domainLookupEnd: e.domainLookupEnd,
				connectStart: e.connectStart,
				connectEnd: e.connectEnd,
				secureConnectionStart: e.secureConnectionStart,
				requestStart: e.requestStart,
				responseStart: e.responseStart,
				responseEnd: e.responseEnd,
				transferSize: e.transferSize,
				encodedBodySize: e.encodedBodySize,
				decodedBodySize: e.decodedBodySize
			}));
	})();
	`
)

type ImageScraper struct {
	timeout          time.Duration
	networkCondition *network.EmulateNetworkConditionsParams
	headless         bool
}

func NewImageScraper() *ImageScraper {
	return &ImageScraper{
		timeout:  defaultTimeout,
		headless: false,
	}
}

func (s *ImageScraper) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *ImageScraper) SetHeadless(headless bool) {
	s.headless = headless
}

func (s *ImageScraper) SetNetworkProfile(profile string) error {
	networkProfiles := getNetworkProfiles()
	p, exists := networkProfiles[profile]
	if !exists {
		return fmt.Errorf("network profile %q not found", profile)
	}

	s.networkCondition = &network.EmulateNetworkConditionsParams{
		Latency:            p.Latency,
		DownloadThroughput: p.Download,
		UploadThroughput:   p.Upload,
		Offline:            false,
	}

	return nil
}

func (s *ImageScraper) ScrapeImages(ctx context.Context, targetURL string) ([]Image, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", s.headless),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("window-size", "1920,1080"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	imagesByRequestID := make(map[network.RequestID]Image)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		handleImageEvents(ev, imagesByRequestID)
	})

	var imgElements []Image

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetCacheDisabled(true),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if s.networkCondition != nil {
				if err := network.EmulateNetworkConditions(
					s.networkCondition.Offline,
					s.networkCondition.Latency,
					s.networkCondition.DownloadThroughput,
					s.networkCondition.UploadThroughput,
				).Do(ctx); err != nil {
					return fmt.Errorf("failed to set network conditions: %w", err)
				}
				log.Printf("Network conditions set: %+v", s.networkCondition)
			}
			return nil
		}),
		chromedp.Navigate(targetURL),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(
				scrollScript,
			).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.Sleep(8*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("error navigating to URL: %w", err)
	}

	err = chromedp.Run(ctx,
		chromedp.Evaluate(evalScript, &imgElements),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting images from DOM: %w", err)
	}

	var images []Image

	for _, img := range imgElements {
		src := cleanURL(img.Src)
		if src == "" {
			continue
		}

		found := false
		for id, netImg := range imagesByRequestID {
			if cleanURL(netImg.Src) == src {
				netImg.Width = img.Width
				netImg.Height = img.Height
				netImg.Alt = img.Alt
				imagesByRequestID[id] = netImg
				found = true
				break
			}
		}

		if !found {
			images = append(images, Image{
				Src:    src,
				Width:  img.Width,
				Height: img.Height,
				Alt:    img.Alt,
			})
		}
	}

	for _, img := range imagesByRequestID {
		if img.Network.MimeType == "image/gif" || img.Network.MimeType == "text/plain" {
			continue
		}
		images = append(images, img)
	}

	var resourceTimings []ResourceTimingEntry
	err = chromedp.Run(ctx, chromedp.Evaluate(resourceTimingScript, &resourceTimings))
	if err != nil {
		log.Printf("Warning: failed to retrieve resource timing data: %v", err)
	} else {
		timingMap := make(map[string]ResourceTimingEntry)
		for _, rt := range resourceTimings {
			cleaned := cleanURL(rt.Name)
			timingMap[cleaned] = rt
		}

		for i := range images {
			if rt, ok := timingMap[images[i].Src]; ok {
				images[i].Timing = convertTiming(rt)
			}
		}
	}

	log.Printf("Found %d unique images", len(images))
	return images, nil
}

func handleImageEvents(ev interface{}, imagesByRequestID map[network.RequestID]Image) {
	switch ev := ev.(type) {
	case *network.EventRequestWillBeSent:
		if ev.Type == network.ResourceTypeImage {
			handleRequestWillBeSent(ev, imagesByRequestID)
		}

	case *network.EventResponseReceived:
		if ev.Type == network.ResourceTypeImage {
			handleResponseReceived(ev, imagesByRequestID)
		}

	case *network.EventLoadingFinished:
		handleLoadingFinished(ev, imagesByRequestID)
	}
}

func handleRequestWillBeSent(ev *network.EventRequestWillBeSent, imagesByRequestID map[network.RequestID]Image) {
	img := imagesByRequestID[ev.RequestID]
	img.Src = cleanURL(ev.Request.URL)

	img.Network.RequestID = ev.RequestID
	img.Network.DocumentURL = ev.DocumentURL
	img.Network.Method = ev.Request.Method
	img.Network.RequestTime = ev.Timestamp

	if ev.Initiator != nil {
		img.Network.InitiatorType = ev.Initiator.Type
		img.Network.InitiatorURL = ev.Initiator.URL
		img.Network.InitiatorLineNo = ev.Initiator.LineNumber
		img.Network.InitiatorColNo = ev.Initiator.ColumnNumber
	}

	imagesByRequestID[ev.RequestID] = img
}

func handleResponseReceived(ev *network.EventResponseReceived, imagesByRequestID map[network.RequestID]Image) {
	img := imagesByRequestID[ev.RequestID]

	img.Network.Status = ev.Response.Status
	img.Network.MimeType = ev.Response.MimeType
	img.Format = normalizeFormat(ev.Response.MimeType)
	img.Network.Protocol = ev.Response.Protocol
	img.Network.RemoteIPAddress = ev.Response.RemoteIPAddress
	img.Network.RemotePort = ev.Response.RemotePort
	img.Network.ResponseTime = ev.Timestamp

	if img.Network.RequestTime != nil {
		img.Network.LoadTime = ev.Timestamp.Time().Sub(img.Network.RequestTime.Time()).Seconds()
	}

	requestHeaders := make(map[string]string)
	for k, v := range ev.Response.RequestHeaders {
		requestHeaders[k] = fmt.Sprintf("%v", v)
	}
	responseHeaders := make(map[string]string)
	for k, v := range ev.Response.Headers {
		responseHeaders[k] = fmt.Sprintf("%v", v)
	}

	img.Network.RequestHeaders = requestHeaders
	img.Network.ResponseHeaders = responseHeaders

	imagesByRequestID[ev.RequestID] = img
}

func handleLoadingFinished(ev *network.EventLoadingFinished, imagesByRequestID map[network.RequestID]Image) {
	if img, exists := imagesByRequestID[ev.RequestID]; exists {
		img.Network.EncodedDataLength = int(ev.EncodedDataLength)
		img.Size = int(ev.EncodedDataLength)
		imagesByRequestID[ev.RequestID] = img
	}
}

func normalizeFormat(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return "image/jpeg"
	case "image/png":
		return "image/png"
	case "image/webp":
		return "image/webp"
	}
	return mimeType
}

func getNetworkProfiles() map[string]NetworkProfile {
	return map[string]NetworkProfile{
		"No Throttling": {
			Download: -1,
			Upload:   -1,
			Latency:  0,
		},
		"Slow 3G": {
			Download: ((500 * 1000) / 8) * 0.8,
			Upload:   ((500 * 1000) / 8) * 0.8,
			Latency:  400 * 5,
		},
		"Fast 3G": {
			Download: ((1.6 * 1000 * 1000) / 8) * 0.9,
			Upload:   ((750 * 1000) / 8) * 0.9,
			Latency:  150 * 3.75,
		},
	}
}

func cleanURL(imgURL string) string {
	u, err := url.Parse(imgURL)
	if err != nil {
		return imgURL
	}

	q := u.Query()
	if originalURL := q.Get("url"); originalURL != "" {
		decoded, err := url.QueryUnescape(originalURL)
		if err != nil {
			return imgURL
		}
		return decoded
	}

	q.Del("w")
	q.Del("q")
	u.RawQuery = q.Encode()

	return u.String()
}

func convertTiming(rt ResourceTimingEntry) TimingInfo {
	timing := TimingInfo{
		DNSLookup:           rt.DomainLookupEnd - rt.DomainLookupStart,
		ConnectionTime:      rt.ConnectEnd - rt.ConnectStart,
		SSLTime:             sslTime(rt),
		TTFB:                rt.ResponseStart - rt.RequestStart,
		ContentDownloadTime: rt.ResponseEnd - rt.ResponseStart,
		TransferSize:        rt.TransferSize,
		EncodedBodySize:     rt.EncodedBodySize,
		DecodedBodySize:     rt.DecodedBodySize,
	}

	return timing
}

func sslTime(rt ResourceTimingEntry) float64 {
	if rt.SecureConnectionStart > 0 {
		return rt.ConnectEnd - rt.SecureConnectionStart
	}
	return 0
}
