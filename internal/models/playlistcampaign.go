package models

import (
    "database/sql"
    "encoding/json"
)

type PlaylistCampaign struct {
    PlaylistID       int            `json:"playlistid" validate:"required,min=1"`
    CampaignID       int            `json:"campaignid" validate:"required,min=1"`
    PlaylisterId     int            `json:"playlisterid" validate:"required,min=1"`
    ReferenceArtists sql.NullString `json:"referenceartists" validate:"required"`
    PlacementStatus  sql.NullString `json:"placementstatus" validate:"required,oneof=Pending Placed Rejected"`
    NumberOfMessages int            `json:"numberofmessages" validate:"min=0"`
    Purchased        bool           `json:"purchased"`
}

func (pc PlaylistCampaign) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        PlaylistID       int    `json:"playlistid"`
        CampaignID       int    `json:"campaignid"`
        PlaylisterId     int    `json:"playlisterid"`
        ReferenceArtists string `json:"referenceartists"`
        PlacementStatus  string `json:"placementstatus"`
        NumberOfMessages int    `json:"numberofmessages"`
        Purchased        bool   `json:"purchased"`
    }{
        PlaylistID:       pc.PlaylistID,
        CampaignID:       pc.CampaignID,
        PlaylisterId:     pc.PlaylisterId,
        ReferenceArtists: pc.ReferenceArtists.String,
        PlacementStatus:  pc.PlacementStatus.String,
        NumberOfMessages: pc.NumberOfMessages,
        Purchased:        pc.Purchased,
    })
}
