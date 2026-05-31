package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo/sqlcgen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrEmailTaken = errors.New("email already taken")
	ErrConflict   = errors.New("conflict")
)

type UserRepo struct {
	q *sqlcgen.Queries
}

func NewUsersRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{q: sqlcgen.New(pool)}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:           u.ID,
		Email:        u.Email.String(),
		DisplayName:  u.DisplayName,
		PasswordHash: pgTextFromPasswordHash(u.PasswordHash),
		TelegramID:   pgTextFromString(u.TelegramID),
	})
	if err != nil {
		if isUniqueViolation(err, "users_email_key") {
			return ErrEmailTaken
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get by email: %w", err)
	}
	return userFromRow(row), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row, err := r.q.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return userFromRow(row), nil
}

func (r *UserRepo) Update(ctx context.Context, u domain.UserUpdate) (*domain.User, error) {
	params := sqlcgen.UpdateUserParams{
		ID:          u.ID,
		DisplayName: u.DisplayName,
	}
	// A non-nil TelegramID means "set it" — including to the empty string, which
	// clears it (stored as NULL). A nil TelegramID leaves the value untouched.
	if u.TelegramID != nil {
		params.SetTelegram = ptrBool(true)
		params.TelegramID = pgTextFromString(*u.TelegramID)
	}

	row, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return userFromRow(row), nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	users, err := r.q.ListUsers(ctx, sqlcgen.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	results := make([]*domain.User, 0, len(users))
	for _, user := range users {
		results = append(results, userFromRow(user))
	}

	return results, nil
}
