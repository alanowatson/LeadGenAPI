package models

type Playlist struct {
    ID                   int    `json:"id"`
    PlaylisterId         int    `json:"playlisterid" validate:"required,min=1"`
    PlaylistSpotifyId    string `json:"playlistspotifyid" validate:"required,min=10,max=100"`
    NumberOfFollowers    int    `json:"numberoffollowers" validate:"min=0"`
    CurrentPlaylistName  string `json:"current_playlist_name" validate:"required,min=1,max=200"`
    LastFollowerCountDate string `json:"lastfollowercountdate" validate:"omitempty,datetime=2006-01-02"`
    LastExposed          string `json:"last_exposed" validate:"omitempty,datetime=2006-01-02"`
}
