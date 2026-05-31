package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mond1c/lms/services/identity/internal/domain"
)

type VCSIdentityRepo interface {
	Upsert(ctx context.Context, vi *domain.VCSIdentity) (*domain.VCSIdentity, error)
	Delete(ctx context.Context, userID string, provider domain.ProviderRef) error
	List(ctx context.Context, userID string) ([]*domain.VCSIdentity, error)
}

type VCSIdentitiesService struct {
	repo VCSIdentityRepo
}

var ErrExternalLoginRequired = fmt.Errorf("external_login is required: %w", ErrValidation)

func NewVCSIdentitiesService(repo VCSIdentityRepo) *VCSIdentitiesService {
	return &VCSIdentitiesService{repo: repo}
}

type LinkVCSIdentityInput struct {
	UserID         string
	Provider       domain.ProviderRef
	ExternalUserID int64
	ExternalLogin  string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      *time.Time
}

func (s *VCSIdentitiesService) Link(ctx context.Context, in LinkVCSIdentityInput) (*domain.VCSIdentity, error) {
	if in.UserID == "" {
		return nil, ErrUserRequired
	}
	if in.ExternalLogin == "" {
		return nil, ErrExternalLoginRequired
	}

	vi := &domain.VCSIdentity{
		UserID:         in.UserID,
		Provider:       in.Provider,
		ExternalUserID: in.ExternalUserID,
		ExternalLogin:  in.ExternalLogin,
		AccessToken:    in.AccessToken,
		RefreshToken:   in.RefreshToken,
		ExpiresAt:      in.ExpiresAt,
		TokenValid:     true,
	}
	out, err := s.repo.Upsert(ctx, vi)
	if err != nil {
		return nil, fmt.Errorf("upsert: %w", err)
	}
	return out, nil
}

func (s *VCSIdentitiesService) Unlink(ctx context.Context, userID string, provider domain.ProviderRef) error {
	if userID == "" {
		return ErrUserRequired
	}
	if err := s.repo.Delete(ctx, userID, provider); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s *VCSIdentitiesService) List(ctx context.Context, userID string) ([]*domain.VCSIdentity, error) {
	if userID == "" {
		return nil, ErrUserRequired
	}
	identities, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return identities, nil
}
