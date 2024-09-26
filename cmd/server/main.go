package main

import (
    "log"
    "net/http"
    "github.com/joho/godotenv"

    "github.com/gorilla/mux"
    "github.com/alanowatson/LeadGenAPI/internal/handlers"
    "github.com/alanowatson/LeadGenAPI/internal/middleware"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	r.HandleFunc("/playlisters", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylisters))).Methods("GET")
	r.HandleFunc("/playlisters", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreatePlaylister))).Methods("POST")
	r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylister))).Methods("GET")
	r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.UpdatePlaylister))).Methods("PUT")
	r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.DeletePlaylister))).Methods("DELETE")

	// Start the cleanup goroutine
	go middleware.CleanupVisitors()

    log.Println("Starting server on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
