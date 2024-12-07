package simage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
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
		window.scrollTo({
			top: document.documentElement.scrollHeight - document.documentElement.clientHeight,
			behavior: 'smooth'
		});
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
		headless: true,
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
	)

	if err != nil {
		return nil, fmt.Errorf("error navigating to URL: %w", err)
	}

	if err := waitForImagesToLoad(ctx, s.timeout); err != nil {
		log.Printf("Warning: timed out waiting for images to fully load: %v", err)
	}

	err = chromedp.Run(ctx,
		chromedp.Evaluate(evalScript, &imgElements),
	)
	if err != nil {
		return nil, fmt.Errorf("error extracting images from DOM: %w", err)
	}

	cleanedURLToReqID := map[string]network.RequestID{}
	for reqID, netImg := range imagesByRequestID {
		cURL := cleanURL(netImg.Src)
		if cURL != "" {
			cleanedURLToReqID[cURL] = reqID
		}
	}

	var images []Image
	uniqueImages := map[string]bool{}

	for _, img := range imgElements {
		src := cleanURL(img.Src)
		if src == "" {
			continue
		}

		if reqID, ok := cleanedURLToReqID[src]; ok {
			netImg := imagesByRequestID[reqID]
			netImg.Width = img.Width
			netImg.Height = img.Height
			netImg.Alt = img.Alt
			imagesByRequestID[reqID] = netImg
			uniqueImages[src] = true
		} else {
			if !uniqueImages[src] {
				images = append(images, Image{
					Src:    src,
					Width:  img.Width,
					Height: img.Height,
					Alt:    img.Alt,
				})
				uniqueImages[src] = true
			}
		}
	}

	for _, img := range imagesByRequestID {
		if img.Network.MimeType == "image/gif" || img.Network.MimeType == "text/plain" {
			continue
		}
		src := cleanURL(img.Src)
		if src != "" && !uniqueImages[src] {
			images = append(images, img)
			uniqueImages[src] = true
		}
	}

	log.Printf("Found %d unique images", len(images))
	return images, nil
}

func waitForImagesToLoad(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var prevCount, stableCount int
	for time.Now().Before(deadline) {
		var count int
		if err := chromedp.Evaluate(`document.images.length`, &count).Do(ctx); err != nil {
			return err
		}

		if count == prevCount {
			stableCount++
			if stableCount > 2 {
				return nil
			}
		} else {
			stableCount = 0
		}

		prevCount = count
		if err := chromedp.Evaluate(scrollScript, nil).Do(ctx); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout reached while waiting for images to load")
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

func GetImageOverview(images []Image) ImageOverview {
	overview := ImageOverview{
		TotalImages: len(images),
		Formats:     make(map[string]int),
	}

	var totalWidth, totalHeight, totalSize int
	var totalLoadTime float64

	for _, img := range images {
		totalSize += img.Size
		if img.Format != "" {
			overview.Formats[img.Format]++
		}

		totalWidth += img.Width
		totalHeight += img.Height
		totalLoadTime += img.Network.LoadTime
	}

	if overview.TotalImages > 0 {
		overview.AverageSize = totalSize / overview.TotalImages
		overview.AverageWidth = totalWidth / overview.TotalImages
		overview.AverageHeight = totalHeight / overview.TotalImages
		overview.TotalSize = totalSize

		overview.AverageTotalTime = totalLoadTime / float64(overview.TotalImages)
	}

	return overview
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
		log.Printf("Warning: failed to parse URL %q: %v", imgURL, err)
		return imgURL
	}

	q := u.Query()

	if originalURL := q.Get("url"); originalURL != "" {
		decoded, err := url.QueryUnescape(originalURL)
		if err == nil {
			u2, err := url.Parse(decoded)
			if err == nil {
				u = u2
				q = u.Query()
			}
		}
	}

	q.Del("w")
	q.Del("q")
	u.RawQuery = q.Encode()

	return u.String()
}
