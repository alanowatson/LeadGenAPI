package models

type Campaign struct {
    ID               int    `json:"id"`
    CampaignName     string `json:"campaign_name" validate:"required,min=1,max=100"`
    ReferenceArtists string `json:"reference_artists" validate:"required"`
    LaunchDate       string `json:"launch_date" validate:"required,datetime=2006-01-02"`
    PromotedArtist   string `json:"promoted_artist" validate:"required,min=1,max=100"`
}
