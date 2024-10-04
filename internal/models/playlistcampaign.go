package models

type PlaylistCampaign struct {
    PlaylistID       int    `json:"playlist_id" validate:"required"`
    CampaignID       int    `json:"campaign_id" validate:"required"`
    PlaylisterId     int    `json:"playlister_id" validate:"required"`
    ReferenceArtists string `json:"reference_artists"`
    PlacementStatus  string `json:"placement_status"`
    NumberOfMessages int    `json:"number_of_messages"`
    Purchased        bool   `json:"purchased"`
}
