package shttp

import (
	"github.com/go-chi/chi/v5"
	"github.com/voage/sharprender-api/db"
	"github.com/voage/sharprender-api/shttp/scan"
)

func NewRouter(mongoClient *db.MongoClient) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/scan", scan.NewScanRoutes(mongoClient.Client))

	return router
}
