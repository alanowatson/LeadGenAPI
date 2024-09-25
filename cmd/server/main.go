package main

import (
	"log"
	"net/http"

	"github.com/alanowatson/LeadGenAPI/internal/handlers"
	"github.com/alanowatson/LeadGenAPI/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Apply middleware
    r.Use(middleware.Authentication)
    r.Use(middleware.RateLimit)

    // Set up routes
    handlers.SetupRoutes(r)

    log.Println("Starting server on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
