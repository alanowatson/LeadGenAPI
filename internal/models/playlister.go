package models

import (
    "database/sql"
    "encoding/json"
)

type Playlister struct {
    ID                int            `json:"playlisterid"`
    SpotifyUserID     sql.NullString `json:"spotifyuserid" validate:"required,min=5,max=50"`
    CuratorFullName   sql.NullString `json:"curatorfullname" validate:"required,min=2,max=100"`
    Email             sql.NullString `json:"email" validate:"required,email"`
    Instagram         sql.NullString `json:"instagram" validate:"omitempty,min=3,max=30"`
    Facebook          sql.NullString `json:"facebook" validate:"omitempty,min=5,max=50"`
    Whatsapp          sql.NullString `json:"whatsapp" validate:"omitempty,e164"`
    LastContacted     sql.NullString `json:"lastcontacted" validate:"omitempty,datetime=2006-01-02"`
    PreferredLanguage sql.NullString `json:"preferredlanguage" validate:"required,iso639_1"`
    FollowupStatus    sql.NullString `json:"followupstatus" validate:"required,oneof=Pending InProgress Completed"`
}

// MarshalJSON implements a custom JSON marshaler for Playlister
func (p Playlister) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID                int    `json:"playlisterid"`
        SpotifyUserID     string `json:"spotifyuserid"`
        CuratorFullName   string `json:"curatorfullname"`
        Email             string `json:"email"`
        Instagram         string `json:"instagram"`
        Facebook          string `json:"facebook"`
        Whatsapp          string `json:"whatsapp"`
        LastContacted     string `json:"lastcontacted"`
        PreferredLanguage string `json:"preferredlanguage"`
        FollowupStatus    string `json:"followupstatus"`
    }{
        ID:                p.ID,
        SpotifyUserID:     stringOrEmpty(p.SpotifyUserID),
        CuratorFullName:   stringOrEmpty(p.CuratorFullName),
        Email:             stringOrEmpty(p.Email),
        Instagram:         stringOrEmpty(p.Instagram),
        Facebook:          stringOrEmpty(p.Facebook),
        Whatsapp:          stringOrEmpty(p.Whatsapp),
        LastContacted:     stringOrEmpty(p.LastContacted),
        PreferredLanguage: stringOrEmpty(p.PreferredLanguage),
        FollowupStatus:    stringOrEmpty(p.FollowupStatus),
    })
}

// Helper function to handle NULL strings
func stringOrEmpty(s sql.NullString) string {
    if !s.Valid || s.String == "NULL" {
        return ""
    }
    return s.String
}
