package repo

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func newUser(t *testing.T, email string) *domain.User {
	t.Helper()
	e, err := domain.NewEmail(email)
	require.NoError(t, err)
	return &domain.User{
		ID:          ulid.Make().String(),
		Email:       e,
		DisplayName: "Test " + email,
	}
}

func TestUserRepo_CreateAndGet(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	u := newUser(t, "alice@example.com")
	require.NoError(t, r.Create(ctx, u))

	t.Run("GetByEmail", func(t *testing.T) {
		got, err := r.GetByEmail(ctx, u.Email.String())
		require.NoError(t, err)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "alice@example.com", got.Email.String())
		require.Equal(t, u.DisplayName, got.DisplayName)
		require.False(t, got.CreatedAt.IsZero(), "created_at should be set")
	})

	t.Run("GetByID", func(t *testing.T) {
		got, err := r.GetByID(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, u.Email.String(), got.Email.String())
	})
}

func TestUserRepo_NotFound(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	_, err := r.GetByID(ctx, "does-not-exist")
	require.ErrorIs(t, err, ErrNotFound)

	missing, _ := domain.NewEmail("nobody@example.com")
	_, err = r.GetByEmail(ctx, missing.String())
	require.ErrorIs(t, err, ErrNotFound)
}

func TestUserRepo_DuplicateEmail(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	u1 := newUser(t, "dup@example.com")
	require.NoError(t, r.Create(ctx, u1))

	u2 := newUser(t, "dup@example.com")
	err := r.Create(ctx, u2)
	require.ErrorIs(t, err, ErrEmailTaken)
}

func TestUserRepo_NullableFields(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	t.Run("no password, no telegram", func(t *testing.T) {
		u := newUser(t, "nopass@example.com")
		require.NoError(t, r.Create(ctx, u))

		got, err := r.GetByID(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, domain.PasswordHash(""), got.PasswordHash)
		require.Equal(t, "", got.TelegramID)
	})

	t.Run("with password and telegram", func(t *testing.T) {
		u := newUser(t, "withpass@example.com")
		h, err := domain.HashPassword("supersecret")
		require.NoError(t, err)
		u.PasswordHash = h
		u.TelegramID = "@alice"
		require.NoError(t, r.Create(ctx, u))

		got, err := r.GetByID(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, h, got.PasswordHash)
		require.Equal(t, "@alice", got.TelegramID)
	})
}

func TestUserRepo_UpdateTelegram(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	u := newUser(t, "tg@example.com")
	u.TelegramID = "@before"
	require.NoError(t, r.Create(ctx, u))

	t.Run("set new value", func(t *testing.T) {
		tg := "@after"
		got, err := r.Update(ctx, domain.UserUpdate{ID: u.ID, TelegramID: &tg})
		require.NoError(t, err)
		require.Equal(t, "@after", got.TelegramID)
	})

	t.Run("nil leaves telegram unchanged", func(t *testing.T) {
		name := "New Name"
		got, err := r.Update(ctx, domain.UserUpdate{ID: u.ID, DisplayName: &name})
		require.NoError(t, err)
		require.Equal(t, "New Name", got.DisplayName)
		require.Equal(t, "@after", got.TelegramID, "telegram must survive a name-only update")
	})

	t.Run("empty string clears telegram", func(t *testing.T) {
		empty := ""
		got, err := r.Update(ctx, domain.UserUpdate{ID: u.ID, TelegramID: &empty})
		require.NoError(t, err)
		require.Equal(t, "", got.TelegramID)

		reloaded, err := r.GetByID(ctx, u.ID)
		require.NoError(t, err)
		require.Equal(t, "", reloaded.TelegramID, "cleared telegram must persist as unset")
	})

	t.Run("missing id", func(t *testing.T) {
		name := "x"
		_, err := r.Update(ctx, domain.UserUpdate{ID: "does-not-exist", DisplayName: &name})
		require.ErrorIs(t, err, ErrNotFound)
	})
}

func TestUserRepo_List(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewUsersRepo(testDB)

	for i, email := range []string{"a@x.com", "b@x.com", "c@x.com"} {
		u := newUser(t, email)
		u.DisplayName = string(rune('A' + i))
		require.NoError(t, r.Create(ctx, u))
	}

	t.Run("list all", func(t *testing.T) {
		got, err := r.List(ctx, 10, 0)
		require.NoError(t, err)
		require.Len(t, got, 3)
	})

	t.Run("paginate", func(t *testing.T) {
		got, err := r.List(ctx, 2, 0)
		require.NoError(t, err)
		require.Len(t, got, 2)

		next, err := r.List(ctx, 2, 2)
		require.NoError(t, err)
		require.Len(t, next, 1)
	})
}
