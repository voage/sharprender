package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/voage/sharprender-api/internal/scraper"
)

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

	scraper := scraper.NewScraper()

	results, err := scraper.ScrapeImages(context.Background(), urlParam)
	if err != nil {
		http.Error(w, "Failed to scrape", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
