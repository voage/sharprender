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
)

// Housekeeping stuff
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

// Bread and butter
func (s *Scraper) Scrape(url string) ([]*Image, error) {
	if s.timeout == 0 {
		s.timeout = defaultTimeout
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	imageURLs := make(map[string]*Image)
	imageRequestIDToURL := make(map[network.RequestID]string)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeImage {
				img, err := NewImage(ev.Response.URL, "", 0, 0, "", 0)
				if err != nil {
					return
				}

				imageURLs[ev.Response.URL] = img

				img.Format = ev.Response.MimeType[6:]
				imageRequestIDToURL[ev.RequestID] = ev.Response.URL

			}
		case *network.EventLoadingFinished:
			if img, ok := imageURLs[imageRequestIDToURL[ev.RequestID]]; ok {
				img.Size = int(ev.EncodedDataLength)
			}

		}
	})

	var images []*Image

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(s.timeout),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result []map[string]interface{}
			err := chromedp.Evaluate(evalScript, &result).Do(ctx)

			if err != nil {
				return err
			}

			for _, imgData := range result {
				src := imgData["src"].(string)
				img, err := NewImage(src, imgData["alt"].(string), int(imgData["width"].(float64)), int(imgData["height"].(float64)), "", 0)
				if err != nil {
					return err
				}

				if networkImg, ok := imageURLs[src]; ok {
					img.Format = networkImg.Format
					img.Size = networkImg.Size
				}

				images = append(images, img)
			}

			return nil
		}),
	)

	if err != nil {
		fmt.Println("Error scraping images:", err)
		return nil, err
	}

	return images, nil
}
