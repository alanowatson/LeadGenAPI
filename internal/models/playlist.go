package models

type Playlist struct {
    ID                   int    `json:"id"`
    PlaylisterId         int    `json:"playlister_id" validate:"required,min=1"`
    PlaylistSpotifyId    string `json:"playlist_spotify_id" validate:"required,min=10,max=100"`
    NumberOfFollowers    int    `json:"number_of_followers" validate:"min=0"`
    CurrentPlaylistName  string `json:"current_playlist_name" validate:"required,min=1,max=200"`
    LastFollowerCountDate string `json:"last_follower_count_date" validate:"omitempty,datetime=2006-01-02"`
    LastExposed          string `json:"last_exposed" validate:"omitempty,datetime=2006-01-02"`
}
