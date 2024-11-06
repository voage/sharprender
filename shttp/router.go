package shttp

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	r.Get("/scrape", getScraperResults)

}
