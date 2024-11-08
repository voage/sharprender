package shttp

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/voage/sharprender-api/internal/simage"
)

type ScraperResponse struct {
	Overview        simage.ImageOverview   `json:"overview"`
	Images          []simage.Image         `json:"images"`
	Recommendations *simage.Recommendation `json:"recommendations"`
}

type AIResponse struct {
	Recommendations simage.Recommendation `json:"recommendations"`
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

	imageScraper := simage.NewImageScraper()

	imageScraper.SetNetworkProfile("Fast 3G")

	results, err := imageScraper.ScrapeImages(r.Context(), urlParam)
	if err != nil {
		log.Printf("Failed to scrape: %v", err)
		http.Error(w, "Failed to scrape", http.StatusInternalServerError)
		return
	}

	overview := simage.GetImageOverview(results)

	var aiRecommendations *simage.Recommendation
	if len(results) > 0 {
		aiRecommendations, err = simage.GetRecommendations(results[0])
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
