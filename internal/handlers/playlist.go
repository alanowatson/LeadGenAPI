package handlers

import (
    "net/http"
)

type Playlist struct {
	PlaylistId           int    `json:"playlist_id"`
	PlaylisterId         int    `json:"playlister_id"`
	PlaylistSpotifyId    string `json:"playlist_spotify_id"`
	NumberOfFollowers    int    `json:"number_of_followers"`
	CurrentPlaylistName  string `json:"current_playlist_name"`
	LastFollowerCountDate string `json:"last_follower_count_date"`
	LastExposed          string `json:"last_exposed"`
}

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
    // Implement the handler logic
}
