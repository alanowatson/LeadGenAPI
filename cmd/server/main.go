package main

import (
    "log"
    "net/http"

    "github.com/alanowatson/LeadGenAPI/internal/handlers"
    "github.com/alanowatson/LeadGenAPI/internal/middleware"
    "github.com/alanowatson/LeadGenAPI/internal/db"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }


    if err := db.InitDB(); err != nil {
        log.Fatalf("Error initializing database: %v", err)
    }

    r := mux.NewRouter()

    // Public routes
    r.HandleFunc("/login", handlers.Login).Methods("POST")

    // Protected routes - Playlisters
    r.HandleFunc("/playlisters", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylisters))).Methods("GET")
    r.HandleFunc("/playlisters", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreatePlaylister))).Methods("POST")
    r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylister))).Methods("GET")
    r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.UpdatePlaylister))).Methods("PUT")
    r.HandleFunc("/playlisters/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.DeletePlaylister))).Methods("DELETE")

    // Protected routes - Playlists
    r.HandleFunc("/playlists", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylists))).Methods("GET")
    r.HandleFunc("/playlists", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreatePlaylist))).Methods("POST")
    r.HandleFunc("/playlists/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylist))).Methods("GET")
    r.HandleFunc("/playlists/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.UpdatePlaylist))).Methods("PUT")
    r.HandleFunc("/playlists/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.DeletePlaylist))).Methods("DELETE")

    // Protected routes - Campaigns
    r.HandleFunc("/campaigns", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetCampaigns))).Methods("GET")
    r.HandleFunc("/campaigns", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreateCampaign))).Methods("POST")
    r.HandleFunc("/campaigns/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetCampaign))).Methods("GET")
    r.HandleFunc("/campaigns/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.UpdateCampaign))).Methods("PUT")
    r.HandleFunc("/campaigns/{id}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.DeleteCampaign))).Methods("DELETE")

    // Protected routes - PlaylistCampaigns
    r.HandleFunc("/playlistcampaigns", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylistCampaigns))).Methods("GET")
    r.HandleFunc("/playlistcampaigns", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreatePlaylistCampaign))).Methods("POST")
    r.HandleFunc("/playlistcampaigns/{playlistId}/{campaignId}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetPlaylistCampaign))).Methods("GET")
    r.HandleFunc("/playlistcampaigns/{playlistId}/{campaignId}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.UpdatePlaylistCampaign))).Methods("PUT")
    r.HandleFunc("/playlistcampaigns/{playlistId}/{campaignId}", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.DeletePlaylistCampaign))).Methods("DELETE")

    // Start the cleanup goroutine
    go middleware.CleanupVisitors()

    log.Println("Starting server on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
