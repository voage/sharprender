package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router) {
	r.Get("/ping", pingHandler)

	r.Get("/scrape", getScraperResults)

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
