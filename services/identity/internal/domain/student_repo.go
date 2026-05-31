package domain

import "time"

// ProvisioningState mirrors the proto ProvisioningState enum.
type ProvisioningState string

const (
	ProvisioningPending ProvisioningState = "pending"
	ProvisioningReady   ProvisioningState = "ready"
	ProvisioningFailed  ProvisioningState = "failed"
)

func (s ProvisioningState) Valid() bool {
	switch s {
	case ProvisioningPending, ProvisioningReady, ProvisioningFailed:
		return true
	}
	return false
}

type StudentRepo struct {
	ID            string
	UserID        string
	AssignmentID  string
	Provider      ProviderRef
	FullName      string
	ExternalID    int64
	State         ProvisioningState
	CloneURLHTTPS string
	CloneURLSSH   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
