package models

import (
    "database/sql"
    "encoding/json"
)

type Playlist struct {
    ID                   int            `json:"playlistid" db:"playlistid"`
    PlaylisterId         int            `json:"playlisterid" db:"playlisterid" validate:"required,min=1"`
    PlaylistSpotifyId    string         `json:"playlistspotifyid" db:"playlistspotifyid" validate:"required,min=10,max=100"`
    NumberOfFollowers    int            `json:"numberoffollowers" db:"numberoffollowers" validate:"min=0"`
    CurrentPlaylistName  sql.NullString `json:"current_playlist_name" db:"current_playlist_name"`
    LastFollowerCountDate sql.NullString `json:"lastfollowercountdate" db:"lastfollowercountdate"`
    LastExposed          sql.NullString `json:"last_exposed" db:"last_exposed"`
}

func (p Playlist) MarshalJSON() ([]byte, error) {
    type Alias Playlist
    return json.Marshal(&struct {
        CurrentPlaylistName  string `json:"current_playlist_name"`
        LastFollowerCountDate string `json:"lastfollowercountdate"`
        LastExposed          string `json:"last_exposed"`
        *Alias
    }{
        CurrentPlaylistName:  p.CurrentPlaylistName.String,
        LastFollowerCountDate: p.LastFollowerCountDate.String,
        LastExposed:          p.LastExposed.String,
        Alias:                (*Alias)(&p),
    })
}
