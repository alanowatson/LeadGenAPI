package models

type Playlist struct {
    ID                   int    `json:"id"  validate:"required"`
    PlaylisterId         int    `json:"playlister_id"  validate:"required"`
    PlaylistSpotifyId    string `json:"playlist_spotify_id"`
    NumberOfFollowers    int    `json:"number_of_followers"`
    CurrentPlaylistName  string `json:"current_playlist_name"`
    LastFollowerCountDate string `json:"last_follower_count_date"`
    LastExposed          string `json:"last_exposed"`
}
