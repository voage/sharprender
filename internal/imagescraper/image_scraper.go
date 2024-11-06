package imagescraper

import (
	"context"
	"fmt"
	"log"
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

	imageURLs := make(map[string]Image)
	imageRequestIDToURL := make(map[network.RequestID]string)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		handleImageEvents(ev, imageURLs, imageRequestIDToURL)
	})

	var images []Image

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(s.timeout),
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result []map[string]interface{}
			if err := chromedp.Evaluate(evalScript, &result).Do(ctx); err != nil {
				return err
			}

			for _, imgData := range result {
				src, ok := imgData["src"].(string)
				if !ok {
					continue
				}

				img := &Image{
					Src:    src,
					Alt:    safeString(imgData["alt"]),
					Width:  safeInt(imgData["width"]),
					Height: safeInt(imgData["height"]),
				}

				if netImg, ok := imageURLs[src]; ok {
					img.Format = netImg.Format
					img.Size = netImg.Size
					img.Network = netImg.Network
					img.Timing = netImg.Timing
					img.Cache = netImg.Cache
				}

				images = append(images, *img)
			}
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("error scraping images: %w", err)
	}

	return images, nil
}

func handleImageEvents(ev interface{}, imageURLs map[string]Image, imageRequestIDToURL map[network.RequestID]string) {
	switch ev := ev.(type) {
	case *network.EventRequestWillBeSent:
		if ev.Type == network.ResourceTypeImage {
			handleRequestWillBeSent(ev, imageURLs, imageRequestIDToURL)
		}

	case *network.EventResponseReceived:
		if ev.Type == network.ResourceTypeImage {
			handleResponseReceived(ev, imageURLs)
		}

	case *network.EventLoadingFinished:
		handleLoadingFinished(ev, imageURLs, imageRequestIDToURL)
	}
}

func handleRequestWillBeSent(ev *network.EventRequestWillBeSent, imageURLs map[string]Image, imageRequestIDToURL map[network.RequestID]string) {
	url := ev.Request.URL
	img := imageURLs[url]

	img.Network.RequestID = string(ev.RequestID)
	img.Network.DocumentURL = ev.DocumentURL
	img.Network.Method = ev.Request.Method
	img.Network.RequestTime = ev.WallTime

	if ev.Initiator != nil {
		img.Network.InitiatorType = string(ev.Initiator.Type)
		img.Network.InitiatorURL = ev.Initiator.URL
		img.Network.InitiatorLineNo = int(ev.Initiator.LineNumber)
		img.Network.InitiatorColNo = int(ev.Initiator.ColumnNumber)
	}

	imageURLs[url] = img
	imageRequestIDToURL[ev.RequestID] = url
}

func handleResponseReceived(ev *network.EventResponseReceived, imageURLs map[string]Image) {
	url := ev.Response.URL
	img := imageURLs[url]

	img.Network.Status = int(ev.Response.Status)
	img.Network.MimeType = ev.Response.MimeType
	img.Format = ev.Response.MimeType
	img.Network.Protocol = ev.Response.Protocol
	img.Network.RemoteIPAddress = ev.Response.RemoteIPAddress
	img.Network.RemotePort = int(ev.Response.RemotePort)
	img.Network.ResponseTime = ev.Timestamp

	img.Cache.FromCache = ev.Response.FromDiskCache || ev.Response.FromPrefetchCache
	img.Cache.CacheHit = ev.Response.FromDiskCache || ev.Response.FromPrefetchCache

	handleCacheHeaders(ev.Response.Headers, &img.Cache)

	imageURLs[url] = img
}

func handleLoadingFinished(ev *network.EventLoadingFinished, imageURLs map[string]Image, imageRequestIDToURL map[network.RequestID]string) {
	if url, ok := imageRequestIDToURL[ev.RequestID]; ok {
		if img, exists := imageURLs[url]; exists {
			img.Network.EncodedDataLength = int(ev.EncodedDataLength)
			img.Size = int(ev.EncodedDataLength)
			imageURLs[url] = img
		}
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
		overview.Formats[img.Format]++

		if img.Cache.CacheHit {
			totalCacheHits++
		}

		totalWidth += img.Width
		totalHeight += img.Height

		if img.Network.RequestTime != nil && img.Network.ResponseTime != nil {
			requestTime := float64(img.Network.RequestTime.Time().Unix())
			responseTime := float64(img.Network.ResponseTime.Time().Unix())
			totalRequestTime += responseTime - requestTime
			totalResponseTime += responseTime
		}

		totalTiming += img.Timing.TotalTime
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

func safeString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func safeInt(val interface{}) int {
	if num, ok := val.(float64); ok {
		return int(num)
	}
	return 0
}
