package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/voage/sharprender-api/shttp"
)

func main() {
	router := chi.NewRouter()
	shttp.SetupRoutes(router)

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
