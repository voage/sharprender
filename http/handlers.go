package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/voage/sharprender-api/internal/imageai"
	"github.com/voage/sharprender-api/internal/imagescraper"
)

type ScraperResponse struct {
	Overview        imagescraper.ImageOverview `json:"overview"`
	Images          []imagescraper.Image       `json:"images"`
	Recommendations *imageai.Recommendation    `json:"recommendations"`
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

	var aiRecommendations *imageai.Recommendation
	if len(results) > 0 {
		aiRecommendations, err = imageai.GetRecommendations(results[0])
		if err != nil {
			log.Printf("Failed to get AI recommendations: %v", err)
		}
	}

	response := ScraperResponse{
		Overview:        overview,
		Images:          results,
		Recommendations: aiRecommendations,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
