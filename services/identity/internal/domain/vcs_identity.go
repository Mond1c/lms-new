package domain

import "time"

type VCSIdentity struct {
	UserID         string
	Provider       ProviderRef
	ExternalUserID int64
	ExternalLogin  string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      *time.Time
	TokenValid     bool
	LinkedAt       time.Time
}
