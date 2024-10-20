package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/voage/sharprender-api/internal/imagescraper"
)

type ScraperResponse struct {
	Overview imagescraper.ImageOverview `json:"overview"`
	Images   []imagescraper.Image       `json:"images"`
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
