package simage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	defaultTimeout = 2 * time.Second
	evalScript     = `
		Array.from(document.images).map(img => ({
			src: img.src,
			alt: img.alt,
			width: img.width,
			height: img.height
		}))
	`
	scrollScript = `
		window.scrollTo(0,document.body.scrollHeight);
	`
)

type ImageScraper struct {
	timeout time.Duration
}

func NewImageScraper() *ImageScraper {
	return &ImageScraper{
		timeout: defaultTimeout,
	}
}

func (s *ImageScraper) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *ImageScraper) ScrapeImages(ctx context.Context, url string) ([]Image, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	imagesByRequestID := make(map[network.RequestID]Image)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		handleImageEvents(ev, imagesByRequestID)
	})

	var images []Image
	var imgElements []Image

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Evaluate(scrollScript, nil).Do(ctx)
		}),
		chromedp.Sleep(s.timeout),

		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Evaluate(evalScript, &imgElements).Do(ctx)
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("error scraping images: %w", err)
	}

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
		images = append(images, img)
	}

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
	img.Format = ev.Response.MimeType
	img.Network.Protocol = ev.Response.Protocol
	img.Network.RemoteIPAddress = ev.Response.RemoteIPAddress
	img.Network.RemotePort = ev.Response.RemotePort
	img.Network.ResponseTime = ev.Timestamp

	if img.Network.RequestTime != nil {
		img.Network.LoadTime = float64(ev.Timestamp.Time().Sub(img.Network.RequestTime.Time()).Seconds())
	}

	img.Cache.FromCache = ev.Response.FromDiskCache || ev.Response.FromPrefetchCache
	img.Cache.CacheHit = ev.Response.FromDiskCache || ev.Response.FromPrefetchCache

	handleCacheHeaders(ev.Response.Headers, &img.Cache)

	imagesByRequestID[ev.RequestID] = img
}

func handleLoadingFinished(ev *network.EventLoadingFinished, imagesByRequestID map[network.RequestID]Image) {
	if img, exists := imagesByRequestID[ev.RequestID]; exists {
		img.Network.EncodedDataLength = int(ev.EncodedDataLength)
		img.Size = int(ev.EncodedDataLength)
		imagesByRequestID[ev.RequestID] = img
	}
}

func handleCacheHeaders(headers map[string]interface{}, cache *CacheInfo) {
	if cacheState, ok := headers["x-cache"]; ok {
		cache.CacheState = cacheState.(string)
	}
	if age, ok := headers["age"]; ok {
		cache.Age, _ = strconv.Atoi(age.(string))
	}
	if expires, ok := headers["expires"]; ok {
		cache.ExpirationTime = expires.(string)
	}
	if lastModified, ok := headers["last-modified"]; ok {
		cache.LastModified = lastModified.(string)
	}
	if etag, ok := headers["etag"]; ok {
		cache.ETag = etag.(string)
	}

	if cache.ETag != "" {
		cache.CacheValidation = "ETag"
	} else if cache.LastModified != "" {
		cache.CacheValidation = "Last-Modified"
	}
}

func GetImageOverview(images []Image) ImageOverview {
	overview := ImageOverview{
		TotalImages: len(images),
		Formats:     make(map[string]int),
	}

	var totalWidth, totalHeight, totalSize, totalCacheHits int
	var totalRequestTime, totalResponseTime, totalTiming float64

	for _, img := range images {
		totalSize += img.Size
		if img.Format != "" {
			overview.Formats[img.Format]++
		}

		if img.Cache.CacheHit {
			totalCacheHits++
		}

		totalWidth += img.Width
		totalHeight += img.Height

		if img.Network.RequestTime != nil && img.Network.ResponseTime != nil {
			requestTime := img.Network.RequestTime.Time()
			responseTime := img.Network.ResponseTime.Time()
			totalRequestTime += float64(responseTime.UnixNano()-requestTime.UnixNano()) / 1e6
			totalResponseTime += float64(responseTime.UnixNano()) / 1e6
		}

	}

	if overview.TotalImages > 0 {
		overview.AverageSize = totalSize / overview.TotalImages
		overview.AverageWidth = totalWidth / overview.TotalImages
		overview.AverageHeight = totalHeight / overview.TotalImages
		overview.TotalSize = totalSize
		overview.CacheHits = totalCacheHits

		overview.AverageRequestTime = totalRequestTime / float64(overview.TotalImages)
		overview.AverageResponseTime = totalResponseTime / float64(overview.TotalImages)
		overview.AverageTotalTime = totalTiming / float64(overview.TotalImages)
	}

	return overview
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
