package handlers

type Campaign struct {
	CampaignId      int    `json:"campaign_id"`
	CampaignName    string `json:"campaign_name"`
	ReferenceArtists string `json:"reference_artists"`
	LaunchDate      string `json:"launch_date"`
	PromotedArtist  string `json:"promoted_artist"`
}
