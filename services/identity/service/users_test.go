package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo"
	"github.com/stretchr/testify/require"
)

type mockUsersRepo struct {
	createFn     func(ctx context.Context, u *domain.User) error
	getByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	listFn       func(ctx context.Context, limit, offset int32) ([]*domain.User, error)
	updateFn     func(ctx context.Context, u *domain.User) (*domain.User, error)
}

func (m *mockUsersRepo) Create(ctx context.Context, u *domain.User) error {
	return m.createFn(ctx, u)
}

func (m *mockUsersRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockUsersRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockUsersRepo) List(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	return m.listFn(ctx, limit, offset)
}

func (m *mockUsersRepo) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	return m.updateFn(ctx, u)
}

func TestUsersService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("valid input, no password", func(t *testing.T) {
		var captured *domain.User

		mock := &mockUsersRepo{
			createFn: func(ctx context.Context, u *domain.User) error {
				captured = u
				return nil
			},
		}
		svc := NewUsersService(mock)

		got, err := svc.Create(ctx, CreateUserInput{
			Email:       "Foo@Bar.com",
			DisplayName: "Foo",
		})
		require.NoError(t, err)
		require.NotEmpty(t, got.ID, "ULID should be generated")
		require.Equal(t, "foo@bar.com", got.Email.String(), "email should be normalized")
		require.Equal(t, domain.PasswordHash(""), got.PasswordHash, "no password => empty hash")
		require.Same(t, captured, got, "service should return the same user it sent to repo")
	})

	t.Run("valid input, with password", func(t *testing.T) {
		mock := &mockUsersRepo{
			createFn: func(_ context.Context, u *domain.User) error { return nil },
		}
		svc := NewUsersService(mock)

		got, err := svc.Create(ctx, CreateUserInput{
			Email:       "foo@bar.com",
			DisplayName: "Foo",
			Password:    "supersecret",
		})
		require.NoError(t, err)
		require.NotEmpty(t, got.PasswordHash)
		require.NoError(t, got.PasswordHash.Verify("supersecret"))
	})

	t.Run("invalid email", func(t *testing.T) {
		mock := &mockUsersRepo{
			createFn: func(_ context.Context, u *domain.User) error {
				t.Fatal("repo.Create should not be called on invalid email")
				return nil
			},
		}
		svc := NewUsersService(mock)

		_, err := svc.Create(ctx, CreateUserInput{
			Email:       "not-an-email",
			DisplayName: "Foo",
		})
		require.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("empty display name", func(t *testing.T) {
		svc := NewUsersService(&mockUsersRepo{})
		_, err := svc.Create(ctx, CreateUserInput{
			Email:       "foo@bar.com",
			DisplayName: "",
		})
		require.ErrorIs(t, err, ErrDisplayNameRequired)
	})

	t.Run("too short password", func(t *testing.T) {
		svc := NewUsersService(&mockUsersRepo{})
		_, err := svc.Create(ctx, CreateUserInput{
			Email:       "foo@bar.com",
			DisplayName: "Foo",
			Password:    "short",
		})
		require.ErrorIs(t, err, domain.ErrPasswordTooShort)
	})

	t.Run("repo error is wrapped", func(t *testing.T) {
		mock := &mockUsersRepo{
			createFn: func(_ context.Context, u *domain.User) error {
				return repo.ErrEmailTaken
			},
		}
		svc := NewUsersService(mock)

		_, err := svc.Create(ctx, CreateUserInput{
			Email:       "taken@bar.com",
			DisplayName: "Foo",
		})
		require.ErrorIs(t, err, repo.ErrEmailTaken, "repo errors must be unwrappable")
	})
}

func TestUsersService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		want := &domain.User{ID: "01HX"}
		mock := &mockUsersRepo{
			getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				require.Equal(t, "01HX", id)
				return want, nil
			},
		}
		svc := NewUsersService(mock)

		got, err := svc.GetByID(ctx, "01HX")
		require.NoError(t, err)
		require.Same(t, want, got)
	})

	t.Run("empty id", func(t *testing.T) {
		svc := NewUsersService(&mockUsersRepo{})
		_, err := svc.GetByID(ctx, "")
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("not found propagates", func(t *testing.T) {
		mock := &mockUsersRepo{
			getByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, repo.ErrNotFound
			},
		}
		svc := NewUsersService(mock)

		_, err := svc.GetByID(ctx, "nope")
		require.ErrorIs(t, err, repo.ErrNotFound)
		require.True(t, errors.Is(err, repo.ErrNotFound))
	})
}
