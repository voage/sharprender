package shttp

import (
	"github.com/go-chi/chi/v5"
	"github.com/voage/sharprender-api/db"
)

func SetupRoutes(r chi.Router, mongoClient *db.MongoClient) {
	handler := NewHandler(mongoClient)

	r.Get("/scan", handler.getScraperResults)
}
