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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScanHandler struct {
	service *ScanService
	repo    *ScanRepository
}

func NewScanHandler(service *ScanService, repo *ScanRepository) *ScanHandler {
	return &ScanHandler{service: service, repo: repo}
}

func (h *ScanHandler) GetScanResults(w http.ResponseWriter, r *http.Request) {
	// Parse scan ID
	id := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid scan ID", http.StatusBadRequest)
		return
	}

	// Parse query filters
	filters := parseFilterOptions(r)

	// Fetch results from service
	result, err := h.service.fetchScanResult(r.Context(), objectID, filters)
	if err != nil {
		http.Error(w, "Failed to fetch scan results", http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ScanHandler) ScanURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	imageScraper := simage.NewImageScraper()

	imageScraper.SetNetworkProfile("No Throttling")

	results, err := imageScraper.ScrapeImages(r.Context(), req.URL)
	if err != nil {
		log.Printf("Failed to scrape: %v", err)
		http.Error(w, "Failed to scrape", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	resultsWithAI, err := simage.CreateAIRecommendations(ctx, results)
	if err != nil {
		log.Printf("Failed to get AI recommendations: %v", err)
		http.Error(w, "Failed to get AI recommendations", http.StatusInternalServerError)
		return
	}

	scan := Scan{
		URL:       req.URL,
		Images:    resultsWithAI,
		CreatedAt: time.Now(),
	}

	id, err := h.repo.Create(context.Background(), &scan)
	if err != nil {
		http.Error(w, "Failed to create scan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"scan_id": id.Hex()})
}
