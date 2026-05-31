package repo

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestVCSIdentityRepo_UpsertListDelete(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewVCSIdentitiesRepo(testDB)

	provider := domain.ProviderRef{Kind: 1, Instance: "gitea.example.com"}
	vi := &domain.VCSIdentity{
		UserID:         "user-1",
		Provider:       provider,
		ExternalUserID: 42,
		ExternalLogin:  "alice",
		AccessToken:    "tok",
		TokenValid:     true,
	}
	out, err := r.Upsert(ctx, vi)
	require.NoError(t, err)
	require.Equal(t, "alice", out.ExternalLogin)
	require.True(t, out.TokenValid)
	require.False(t, out.LinkedAt.IsZero())

	t.Run("re-link updates login", func(t *testing.T) {
		vi.ExternalLogin = "alice2"
		vi.AccessToken = "tok2"
		out, err := r.Upsert(ctx, vi)
		require.NoError(t, err)
		require.Equal(t, "alice2", out.ExternalLogin)

		list, err := r.List(ctx, "user-1")
		require.NoError(t, err)
		require.Len(t, list, 1, "re-link must not create a duplicate")
		require.Equal(t, "alice2", list[0].ExternalLogin)
		require.Equal(t, "tok2", list[0].AccessToken)
	})

	t.Run("delete then delete-missing", func(t *testing.T) {
		require.NoError(t, r.Delete(ctx, "user-1", provider))
		require.ErrorIs(t, r.Delete(ctx, "user-1", provider), ErrNotFound)
	})
}
