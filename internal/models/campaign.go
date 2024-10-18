package models

type Campaign struct {
    ID               int    `json:"campaignid"`
    CampaignName     string `json:"campaignname" validate:"required,min=1,max=100"`
    ReferenceArtists string `json:"referenceartists" validate:"required"`
    TrelloLink       string `json:"trello_link"`
    SpotifyLink      string `json:"spotify_link"`
    LaunchDate       string `json:"launch_date" validate:"required,datetime=2006-01-02"`
    PromotedArtist   string `json:"promoted_artist" validate:"required,min=1,max=100"`
}
