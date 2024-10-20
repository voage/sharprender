package main

import (
	"context"
	"log"

	"github.com/voage/sharprender-api/internal/scraper"
)

func main() {
	scraper := scraper.NewScraper()

	_, err := scraper.ScrapeImages(context.Background(), "https://www.ycombinator.com")
	if err != nil {
		log.Fatalf("Error scraping: %v", err)
	}

}
