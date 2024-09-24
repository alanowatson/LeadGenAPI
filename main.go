package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

type Playlister struct {
	PlaylisterId      int    `json:"playlister_id"`
	SpotifyUserId     string `json:"spotify_user_id"`
	CuratorFullName   string `json:"curator_full_name"`
	Email             string `json:"email"`
	Instagram         string `json:"instagram"`
	Facebook          string `json:"facebook"`
	Whatsapp          string `json:"whatsapp"`
	LastContacted     string `json:"last_contacted"`
	PreferredLanguage string `json:"preferred_language"`
	FollowupStatus    string `json:"followup_status"`
}

type Playlist struct {
	PlaylistId           int    `json:"playlist_id"`
	PlaylisterId         int    `json:"playlister_id"`
	PlaylistSpotifyId    string `json:"playlist_spotify_id"`
	NumberOfFollowers    int    `json:"number_of_followers"`
	CurrentPlaylistName  string `json:"current_playlist_name"`
	LastFollowerCountDate string `json:"last_follower_count_date"`
	LastExposed          string `json:"last_exposed"`
}

type Campaign struct {
	CampaignId      int    `json:"campaign_id"`
	CampaignName    string `json:"campaign_name"`
	ReferenceArtists string `json:"reference_artists"`
	LaunchDate      string `json:"launch_date"`
	PromotedArtist  string `json:"promoted_artist"`
}

type PlaylistCampaign struct {
	PlaylistId        int    `json:"playlist_id"`
	CampaignId        int    `json:"campaign_id"`
	PlaylisterId      int    `json:"playlister_id"`
	ReferenceArtists  string `json:"reference_artists"`
	PlacementStatus   string `json:"placement_status"`
	NumberOfMessages  int    `json:"number_of_messages"`
	Purchased         bool   `json:"purchased"`
}

var (
	tokenCache *cache.Cache
	limiter    *rate.Limiter
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func getPagination(r *http.Request) (int, int) {
	page := 1
	perPage := 25
	// In a real application, you'd parse these from query parameters
	return page, perPage
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		// For testing purposes, accept any non-empty token
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getPlaylistersHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement database query
	playlisters := []Playlister{} // This should be populated from the database
	page, perPage := getPagination(r)

	// TODO: Implement filtering and sorting

	// Simulate pagination
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(playlisters) {
		end = len(playlisters)
	}
	paginatedPlaylisters := playlisters[start:end]

	respondWithJSON(w, http.StatusOK, paginatedPlaylisters)
}

func getPlaylisterHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]

	// TODO: Implement database query to get playlister by ID
	playlister := Playlister{} // This should be populated from the database

	respondWithJSON(w, http.StatusOK, playlister)
}

func createPlaylisterHandler(w http.ResponseWriter, r *http.Request) {
	var playlister Playlister
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playlister); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: Implement database insertion

	respondWithJSON(w, http.StatusCreated, playlister)
}

func updatePlaylisterHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]

	var playlister Playlister
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playlister); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// TODO: Implement database update

	respondWithJSON(w, http.StatusOK, playlister)
}

func deletePlaylisterHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]

	// TODO: Implement database deletion

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func main() {
	fmt.Println("Starting the Beast API server...")

	tokenCache = cache.New(24*time.Hour, 1*time.Hour)
	limiter = rate.NewLimiter(rate.Limit(2), 1) // 2 requests per second

	router := mux.NewRouter()

	router.Use(authenticationMiddleware)
	router.Use(rateLimitMiddleware)

	router.HandleFunc("/api/playlisters", getPlaylistersHandler).Methods("GET")
	router.HandleFunc("/api/playlisters/{id}", getPlaylisterHandler).Methods("GET")
	router.HandleFunc("/api/playlisters", createPlaylisterHandler).Methods("POST")
	router.HandleFunc("/api/playlisters/{id}", updatePlaylisterHandler).Methods("PUT")
	router.HandleFunc("/api/playlisters/{id}", deletePlaylisterHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
