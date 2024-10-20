package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Image struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Size   int    `json:"size"`
}

func Scrape() {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	url := "https://www.ycombinator.com"
	var images []Image
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('img')).map(img => ({
				src: img.src,
				alt: img.alt,
				width: img.naturalWidth || img.width,
				height: img.naturalHeight || img.height
			}))
		`, &images),
	)
	if err != nil {
		fmt.Println("Error scraping images:", err)
		return
	}
	for _, img := range images {
		fmt.Printf("Src: %s, Alt: %s, Width: %d, Height: %d\n", img.Src, img.Alt, img.Width, img.Height)
	}
	fmt.Println(images)

}
