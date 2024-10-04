package models

type Playlister struct {
    ID                int    `json:"id"`
    SpotifyUserID     string `json:"spotify_user_id" validate:"required"`
    CuratorFullName   string `json:"curator_full_name" validate:"required"`
    Email             string `json:"email" validate:"required,email"`
    Instagram         string `json:"instagram"`
    Facebook          string `json:"facebook"`
    Whatsapp          string `json:"whatsapp"`
    LastContacted     string `json:"last_contacted" validate:"omitempty,datetime=2006-01-02"`
    PreferredLanguage string `json:"preferred_language"`
    FollowupStatus    string `json:"followup_status" validate:"oneof=Pending InProgress Completed"`
}
