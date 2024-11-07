package models

import (
    "database/sql"
    "encoding/json"
)

type Playlist struct {
    ID                   int            `json:"playlistid"`
    PlaylisterId         int            `json:"playlisterid" validate:"required,min=1"`
    PlaylistSpotifyId    sql.NullString `json:"playlistspotifyid" validate:"required,min=10,max=100"`
    NumberOfFollowers    int            `json:"numberoffollowers" validate:"min=0"`
    CurrentPlaylistName  sql.NullString `json:"current_playlist_name" validate:"required,min=1,max=200"`
    LastFollowerCountDate sql.NullString `json:"lastfollowercountdate" validate:"omitempty,datetime=2006-01-02"`
    LastExposed          sql.NullString `json:"last_exposed" validate:"omitempty,datetime=2006-01-02"`
}

// MarshalJSON implements a custom JSON marshaler for Playlist
func (p Playlist) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID                   int    `json:"playlistid"`
        PlaylisterId         int    `json:"playlisterid"`
        PlaylistSpotifyId    string `json:"playlistspotifyid"`
        NumberOfFollowers    int    `json:"numberoffollowers"`
        CurrentPlaylistName  string `json:"current_playlist_name"`
        LastFollowerCountDate string `json:"lastfollowercountdate"`
        LastExposed          string `json:"last_exposed"`
    }{
        ID:                   p.ID,
        PlaylisterId:         p.PlaylisterId,
        PlaylistSpotifyId:    stringOrEmpty(p.PlaylistSpotifyId),
        NumberOfFollowers:    p.NumberOfFollowers,
        CurrentPlaylistName:  stringOrEmpty(p.CurrentPlaylistName),
        LastFollowerCountDate: stringOrEmpty(p.LastFollowerCountDate),
        LastExposed:          stringOrEmpty(p.LastExposed),
    })
}


