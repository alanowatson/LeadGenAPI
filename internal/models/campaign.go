package models

type Campaign struct {
    ID               int    `json:"id" validate:"required"`
    CampaignName     string `json:"campaign_name"`
    ReferenceArtists string `json:"reference_artists"`
    LaunchDate       string `json:"launch_date"`
    PromotedArtist   string `json:"promoted_artist"`
}
