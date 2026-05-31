package repo

import (
	"context"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo/sqlcgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VCSIdentityRepo struct {
	q *sqlcgen.Queries
}

func NewVCSIdentitiesRepo(pool *pgxpool.Pool) *VCSIdentityRepo {
	return &VCSIdentityRepo{q: sqlcgen.New(pool)}
}

func (r *VCSIdentityRepo) Upsert(ctx context.Context, vi *domain.VCSIdentity) (*domain.VCSIdentity, error) {
	row, err := r.q.UpsertVCSIdentity(ctx, sqlcgen.UpsertVCSIdentityParams{
		UserID:           vi.UserID,
		ProviderKind:     vi.Provider.Kind,
		ProviderInstance: vi.Provider.Instance,
		ExternalUserID:   vi.ExternalUserID,
		ExternalLogin:    vi.ExternalLogin,
		AccessToken:      pgTextFromString(vi.AccessToken),
		RefreshToken:     pgTextFromString(vi.RefreshToken),
		ExpiresAt:        pgTimestamp(vi.ExpiresAt),
	})
	if err != nil {
		return nil, fmt.Errorf("upsert vcs identity: %w", err)
	}
	return vcsIdentityFromRow(row), nil
}

func (r *VCSIdentityRepo) Delete(ctx context.Context, userID string, provider domain.ProviderRef) error {
	rows, err := r.q.DeleteVCSIdentity(ctx, sqlcgen.DeleteVCSIdentityParams{
		UserID:           userID,
		ProviderKind:     provider.Kind,
		ProviderInstance: provider.Instance,
	})
	if err != nil {
		return fmt.Errorf("delete vcs identity: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *VCSIdentityRepo) List(ctx context.Context, userID string) ([]*domain.VCSIdentity, error) {
	rows, err := r.q.ListVCSIdentities(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list vcs identities: %w", err)
	}
	results := make([]*domain.VCSIdentity, 0, len(rows))
	for _, row := range rows {
		results = append(results, vcsIdentityFromRow(row))
	}
	return results, nil
}

func vcsIdentityFromRow(row sqlcgen.VcsIdentity) *domain.VCSIdentity {
	return &domain.VCSIdentity{
		UserID: row.UserID,
		Provider: domain.ProviderRef{
			Kind:     row.ProviderKind,
			Instance: row.ProviderInstance,
		},
		ExternalUserID: row.ExternalUserID,
		ExternalLogin:  row.ExternalLogin,
		AccessToken:    strFromPtr(row.AccessToken),
		RefreshToken:   strFromPtr(row.RefreshToken),
		ExpiresAt:      timeFromPg(row.ExpiresAt),
		TokenValid:     row.TokenValid,
		LinkedAt:       row.LinkedAt.Time,
	}
}
