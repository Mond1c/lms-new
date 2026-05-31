package repo

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestStudentRepoRepo_RegisterGet(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewStudentReposRepo(testDB)

	sr := &domain.StudentRepo{
		ID:            ulid.Make().String(),
		UserID:        "user-1",
		AssignmentID:  "asgn-1",
		Provider:      domain.ProviderRef{Kind: 1, Instance: "gitea.example.com"},
		FullName:      "org/user-1-hw1",
		ExternalID:    7,
		State:         domain.ProvisioningReady,
		CloneURLHTTPS: "https://gitea/org/user-1-hw1.git",
	}
	out, err := r.Register(ctx, sr)
	require.NoError(t, err)
	require.Equal(t, domain.ProvisioningReady, out.State)
	require.False(t, out.CreatedAt.IsZero())

	t.Run("get", func(t *testing.T) {
		got, err := r.Get(ctx, "user-1", "asgn-1")
		require.NoError(t, err)
		require.Equal(t, "org/user-1-hw1", got.FullName)
		require.Equal(t, int64(7), got.ExternalID)
	})

	t.Run("re-register is idempotent and updates", func(t *testing.T) {
		sr.FullName = "org/renamed"
		out, err := r.Register(ctx, sr)
		require.NoError(t, err)
		require.Equal(t, "org/renamed", out.FullName)

		got, err := r.Get(ctx, "user-1", "asgn-1")
		require.NoError(t, err)
		require.Equal(t, "org/renamed", got.FullName)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := r.Get(ctx, "user-x", "asgn-x")
		require.ErrorIs(t, err, ErrNotFound)
	})
}
