package scraper

import (
	"context"
	"fmt"
	"log"
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
	mimeTypeImagePrefix = "image/"
)

type Scraper struct {
	timeout time.Duration
}

func NewScraper() *Scraper {
	return &Scraper{
		timeout: defaultTimeout,
	}
}

func (s *Scraper) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *Scraper) ScrapeImages(ctx context.Context, url string) ([]Image, error) {
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
					Alt:    SafeString(imgData["alt"]),
					Width:  SafeInt(imgData["width"]),
					Height: SafeInt(imgData["height"]),
				}

				if netImg, ok := imageURLs[src]; ok {
					img.Format = netImg.Format
					img.Size = netImg.Size
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
	case *network.EventResponseReceived:
		if ev.Type == network.ResourceTypeImage {
			img := Image{
				Src: ev.Response.URL,
			}
			img.Format = ev.Response.MimeType[len(mimeTypeImagePrefix):]
			imageURLs[ev.Response.URL] = img
			imageRequestIDToURL[ev.RequestID] = ev.Response.URL
		}
	case *network.EventLoadingFinished:
		if url, ok := imageRequestIDToURL[ev.RequestID]; ok {
			if img, exists := imageURLs[url]; exists {
				img.Size = int(ev.EncodedDataLength)
				imageURLs[url] = img
			}
		}
	}
}
