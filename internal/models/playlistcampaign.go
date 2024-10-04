package models

type PlaylistCampaign struct {
	PlaylistID       int    `json:"playlist_id" validate:"required,min=1"`
	CampaignID       int    `json:"campaign_id" validate:"required,min=1"`
	PlaylisterId     int    `json:"playlister_id" validate:"required,min=1"`
	ReferenceArtists string `json:"reference_artists" validate:"required"`
	PlacementStatus  string `json:"placement_status" validate:"required,oneof=Pending Placed Rejected"`
	NumberOfMessages int    `json:"number_of_messages" validate:"min=0"`
	Purchased        bool   `json:"purchased"`
}
