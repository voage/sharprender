package scraper

import (
	"context"
	"fmt"
	"time"
	"path"
	"strings"
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

func getImageFormat(src string) string {
	// Check if the image is a base64-encoded image
	if strings.HasPrefix(src, "data:image/") {
		// Extract the format from the data URL (e.g., "data:image/png;base64,...")
		parts := strings.Split(src, ";")
		if len(parts) > 0 && strings.HasPrefix(parts[0], "data:image/") {
			return strings.TrimPrefix(parts[0], "data:image/")
		}
		return "none"
	}
	ext := strings.ToLower(path.Ext(src))
	if len(ext) > 1 {
		return ext[1:] // remove dot
	}

	return "unknown"	// if no extension
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
	for i := range images {
		images[i].Format = getImageFormat(images[i].Src)	}
	for _, img := range images {
		fmt.Printf("Src: %s, Alt: %s, Width: %d, Height: %d\n, Format: %s\n", img.Src, img.Alt, img.Width, img.Height, img.Format)
	}
	fmt.Println(images)

}
