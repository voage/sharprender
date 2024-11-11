package scan

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/voage/sharprender-api/internal/simage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScanHandler struct {
	repo *ScanRepository
}

func NewScanHandler(repo *ScanRepository) *ScanHandler {
	return &ScanHandler{repo: repo}
}

func (h *ScanHandler) GetScanResults(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	scan, err := h.repo.FindOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(scan)
}

func (h *ScanHandler) ScanURL(w http.ResponseWriter, r *http.Request) {

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

	scan := Scan{
		URL:       urlParam,
		Overview:  overview,
		Images:    results,
		CreatedAt: time.Now(),
	}

	err = h.repo.Create(context.Background(), &scan)
	if err != nil {
		http.Error(w, "Failed to create scan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scan)
}
