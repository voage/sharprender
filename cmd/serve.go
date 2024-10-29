package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	api "github.com/voage/sharprender-api/http"
	"github.com/voage/sharprender-api/internal/imageai"
	"github.com/voage/sharprender-api/internal/imagescraper"
)

func main() {
	router := chi.NewRouter()
	api.SetupRoutes(router)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		log.Fatalf("OPENAI_KEY not found")
	}

	

	urlParam := "https://www.ycombinator.com/"
	
	scraper := imagescraper.NewImageScraper()
	scraper.SetTimeout(2 * time.Second) 

	
	results, err := scraper.ScrapeImages(context.Background(), urlParam)
	if err != nil {
		log.Fatalf("Failed to scrape images: %v", err)
	}

	
	if len(results) == 0 {
		log.Fatalf("No images found at the URL")
	}

	
	firstImage := results[0]

	
	aiRecommendations, err := imageai.GetRecommendations(firstImage, apiKey)
	if err != nil {
		log.Fatalf("Failed to get AI recommendations: %v", err)
	}

	fmt.Printf("AI Recommendations for %s: %+v\n", firstImage.Src, *aiRecommendations)

	log.Printf("Starting server on :%s...", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
