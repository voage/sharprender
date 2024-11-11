package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/voage/sharprender-api/db"
	"github.com/voage/sharprender-api/shttp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient, err := db.InitMongoDB(ctx)
	if err != nil {
		log.Fatalf("Error initializing MongoDB: %s", err)
	}
	defer mongoClient.Disconnect(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := shttp.NewRouter(mongoClient)

	log.Printf("Starting server on :%s", port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
