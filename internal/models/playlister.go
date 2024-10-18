package models

type Playlister struct {
	ID                int    `json:"playlisterid"`
	SpotifyUserID     string `json:"spotifyuserid" validate:"required,min=5,max=50"`
	CuratorFullName   string `json:"curatorfullname" validate:"required,min=2,max=100"`
	Email             string `json:"email" validate:"required,email"`
	Instagram         string `json:"instagram" validate:"omitempty,min=3,max=30"`
	Facebook          string `json:"facebook" validate:"omitempty,min=5,max=50"`
	Whatsapp          string `json:"whatsapp" validate:"omitempty,e164"`
	LastContacted     string `json:"lastcontacted" validate:"omitempty,datetime=2006-01-02"`
	PreferredLanguage string `json:"preferredlanguage" validate:"required,iso639_1"`
	FollowupStatus    string `json:"followupstatus" validate:"required,oneof=Pending InProgress Completed"`
}
