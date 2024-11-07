package models

import (
	"database/sql"
	"encoding/json"
)

type Campaign struct {
    ID               int            `json:"campaignid"`
    CampaignName     sql.NullString `json:"campaignname" validate:"required,min=1,max=100"`
    ReferenceArtists sql.NullString `json:"referenceartists" validate:"required"`
    TrelloLink       sql.NullString `json:"trello_link"`
    SpotifyLink      sql.NullString `json:"spotify_link"`
    LaunchDate       sql.NullString `json:"launch_date" validate:"required,datetime=2006-01-02"`
    PromotedArtist   sql.NullString `json:"promoted_artist" validate:"required,min=1,max=100"`
}

// MarshalJSON implements a custom JSON marshaler for Campaign
func (c Campaign) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        ID               int    `json:"campaignid"`
        CampaignName     string `json:"campaignname"`
        ReferenceArtists string `json:"referenceartists"`
        TrelloLink       string `json:"trello_link"`
        SpotifyLink      string `json:"spotify_link"`
        LaunchDate       string `json:"launch_date"`
        PromotedArtist   string `json:"promoted_artist"`
    }{
        ID:               c.ID,
        CampaignName:     stringOrEmpty(c.CampaignName),
        ReferenceArtists: stringOrEmpty(c.ReferenceArtists),
        TrelloLink:       stringOrEmpty(c.TrelloLink),
        SpotifyLink:      stringOrEmpty(c.SpotifyLink),
        LaunchDate:       stringOrEmpty(c.LaunchDate),
        PromotedArtist:   stringOrEmpty(c.PromotedArtist),
    })
}
