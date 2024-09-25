package handlers

import (
    "net/http"
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

func GetPlaylistersList(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

func GetPlaylister(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// ... other handler functions ...
