package main

import (
	"fmt"
	"log"

	"github.com/voage/sharprender-api/internal/scraper"
)

func main() {
	scraper := scraper.NewScraper()

	images, err := scraper.Scrape("https://www.ycombinator.com")
	if err != nil {
		log.Fatalf("Error scraping: %v", err)
	}

	for _, img := range images {
		fmt.Println(img)
	}
}
