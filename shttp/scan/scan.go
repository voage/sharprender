package scan

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewScanRoutes(mongoClient *mongo.Client) *chi.Mux {
	repo := NewScanRepository(mongoClient)
	service := NewScanService(repo)
	handler := NewScanHandler(service, repo)

	router := chi.NewRouter()
	router.Get("/{id}", handler.GetScanResults)
	router.Post("/", handler.ScanURL)

	return router
}
