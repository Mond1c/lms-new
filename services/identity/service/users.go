package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.User, error)
}

type UsersService struct {
	repo UserRepo
}

var (
	ErrDisplayNameRequired = errors.New("display_name is required")
	ErrInvalidID           = errors.New("invalid id")
)

func NewUsersService(repo UserRepo) *UsersService {
	return &UsersService{repo: repo}
}

type CreateUserInput struct {
	Email       string
	DisplayName string
	Password    string
	TelegramID  string
}

func (s *UsersService) Create(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if input.DisplayName == "" {
		return nil, ErrDisplayNameRequired
	}

	var hash domain.PasswordHash
	if input.Password != "" {
		hash, err = domain.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
	}

	user := &domain.User{
		ID:           ulid.Make().String(),
		Email:        email,
		DisplayName:  input.DisplayName,
		PasswordHash: hash,
		TelegramID:   input.TelegramID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return user, nil
}

func (s *UsersService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return user, nil
}

func (s *UsersService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get by email: %w", err)
	}
	return user, nil
}
