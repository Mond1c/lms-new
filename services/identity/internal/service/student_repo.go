package service

import (
	"context"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
)

type StudentRepoRepo interface {
	Register(ctx context.Context, sr *domain.StudentRepo) (*domain.StudentRepo, error)
	Get(ctx context.Context, userID, assignmentID string) (*domain.StudentRepo, error)
}

type StudentReposService struct {
	repo StudentRepoRepo
}

var (
	ErrStudentRepoAssignmentRequired = fmt.Errorf("assignment_id is required: %w", ErrValidation)
	ErrStudentRepoFullNameRequired   = fmt.Errorf("full_name is required: %w", ErrValidation)
)

func NewStudentReposService(repo StudentRepoRepo) *StudentReposService {
	return &StudentReposService{repo: repo}
}

type RegisterStudentRepoInput struct {
	UserID        string
	AssignmentID  string
	Provider      domain.ProviderRef
	FullName      string
	ExternalID    int64
	CloneURLHTTPS string
	CloneURLSSH   string
}

func (s *StudentReposService) Register(ctx context.Context, in RegisterStudentRepoInput) (*domain.StudentRepo, error) {
	if in.UserID == "" {
		return nil, ErrUserRequired
	}
	if in.AssignmentID == "" {
		return nil, ErrStudentRepoAssignmentRequired
	}
	if in.FullName == "" {
		return nil, ErrStudentRepoFullNameRequired
	}

	sr := &domain.StudentRepo{
		ID:            ulid.Make().String(),
		UserID:        in.UserID,
		AssignmentID:  in.AssignmentID,
		Provider:      in.Provider,
		FullName:      in.FullName,
		ExternalID:    in.ExternalID,
		State:         domain.ProvisioningReady,
		CloneURLHTTPS: in.CloneURLHTTPS,
		CloneURLSSH:   in.CloneURLSSH,
	}
	out, err := s.repo.Register(ctx, sr)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	return out, nil
}

func (s *StudentReposService) Get(ctx context.Context, userID, assignmentID string) (*domain.StudentRepo, error) {
	if userID == "" {
		return nil, ErrUserRequired
	}
	if assignmentID == "" {
		return nil, ErrStudentRepoAssignmentRequired
	}
	sr, err := s.repo.Get(ctx, userID, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	return sr, nil
}
