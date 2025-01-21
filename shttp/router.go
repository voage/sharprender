package shttp

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/voage/sharprender-api/db"
	"github.com/voage/sharprender-api/shttp/scan"
)

func NewRouter(mongoClient *db.MongoClient) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Cache preflight response for 5 minutes
	}))
	router.Mount("/scan", scan.NewScanRoutes(mongoClient.Client))

	return router
}
