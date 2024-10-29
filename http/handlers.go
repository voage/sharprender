package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/voage/sharprender-api/internal/imageai"
	"github.com/voage/sharprender-api/internal/imagescraper"
)

type ScraperResponse struct {
	Overview imagescraper.ImageOverview `json:"overview"`
	Images   []imagescraper.Image       `json:"images"`
}
type AIResponse struct {
	Recommendations imageai.Recommendation `json:"recommendations"`
}

func getScraperResults(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	imageScraper := imagescraper.NewImageScraper()

	results, err := imageScraper.ScrapeImages(r.Context(), urlParam)
	if err != nil {
		http.Error(w, "Failed to scrape", http.StatusInternalServerError)
		return
	}

	overview := imagescraper.GetImageOverview(results)

	response := ScraperResponse{
		Overview: overview,
		Images:   results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getAIReccomendationsResults(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	scraper := imagescraper.NewImageScraper()
	scraper.SetTimeout(2 * time.Second)

	results, err := scraper.ScrapeImages(context.Background(), urlParam)
	if err != nil {
		http.Error(w, "Failed to scrape images", http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		http.Error(w, "No images found at the specified URL", http.StatusNotFound)
		return
	}

	firstImage := results[0]
	aiRecommendations, err := imageai.GetRecommendations(firstImage)
	if err != nil {
		log.Fatalf("Failed to get AI recommendations: %v", err)
	}
	response := AIResponse{
		Recommendations: *aiRecommendations,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
