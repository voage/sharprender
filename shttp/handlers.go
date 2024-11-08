package shttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/voage/sharprender-api/db"
	"github.com/voage/sharprender-api/internal/simage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	mongoClient *db.MongoClient
}

func NewHandler(client *db.MongoClient) *Handler {
	return &Handler{mongoClient: client}
}

type ScraperResponse struct {
	Overview        simage.ImageOverview   `json:"overview"`
	Images          []simage.Image         `json:"images"`
	Recommendations *simage.Recommendation `json:"recommendations"`
}

type AIResponse struct {
	Recommendations simage.Recommendation `json:"recommendations"`
}

type ScanResult struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	URL       string               `bson:"url"`
	Overview  simage.ImageOverview `bson:"overview"`
	Images    []simage.Image       `bson:"images"`
	CreatedAt time.Time            `bson:"created_at"`
}

func (h *Handler) getScraperResults(w http.ResponseWriter, r *http.Request) {
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

	response := ScraperResponse{
		Overview: overview,
		Images:   results,
	}

	err = insertScanResult(h.mongoClient, ScanResult{
		URL:      urlParam,
		Overview: overview,
		Images:   results,
	})
	if err != nil {
		log.Printf("Error inserting scan result: %v", err)
		http.Error(w, "Failed to insert scan result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func insertScanResult(client *db.MongoClient, result ScanResult) error {
	collection := client.Database("sharprenderdb").Collection("scans")
	result.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, result)
	if err != nil {
		log.Printf("Error inserting scan result: %v", err)
		return err
	}

	log.Println("Scan result inserted successfully")
	return nil
}
