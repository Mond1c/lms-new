package domain

import "time"

type Course struct {
	ID           string
	Code         string
	Title        string
	Description  string
	InstructorID string
	VCS          *VCSBinding
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VCSBinding struct {
	Provider       ProviderRef
	TargetOrg      string
	StudentTeam    string
	ReviewerTeam   string
	ReviewerLogins []string
}
