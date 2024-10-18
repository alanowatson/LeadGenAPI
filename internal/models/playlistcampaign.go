package models

type PlaylistCampaign struct {
	PlaylistID       int    `json:"playlistid" validate:"required,min=1"`
	CampaignID       int    `json:"campaignid" validate:"required,min=1"`
	PlaylisterId     int    `json:"playlisterid" validate:"required,min=1"`
	ReferenceArtists string `json:"referenceartists" validate:"required"`
	PlacementStatus  string `json:"placementstatus" validate:"required,oneof=Pending Placed Rejected"`
	NumberOfMessages int    `json:"numberofmessages" validate:"min=0"`
	Purchased        bool   `json:"purchased"`
}
