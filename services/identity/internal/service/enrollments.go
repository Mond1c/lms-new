package service

import (
	"context"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
)

type EnrollmentRepo interface {
	Create(ctx context.Context, e *domain.Enrollment) error
	Delete(ctx context.Context, userID, courseID string) error
	List(ctx context.Context, courseID, userID string, limit, offset int32) ([]*domain.Enrollment, error)
}

type EnrollmentsService struct {
	repo EnrollmentRepo
}

var (
	ErrUserRequired   = fmt.Errorf("user_id is required: %w", ErrValidation)
	ErrCourseRequired = fmt.Errorf("course_id is required: %w", ErrValidation)
	ErrInvalidRole    = fmt.Errorf("invalid role: %w", ErrValidation)
)

func NewEnrollmentsService(repo EnrollmentRepo) *EnrollmentsService {
	return &EnrollmentsService{repo: repo}
}

type EnrollInput struct {
	UserID   string
	CourseID string
	Role     domain.Role
}

func (s *EnrollmentsService) Enroll(ctx context.Context, in EnrollInput) (*domain.Enrollment, error) {
	if in.UserID == "" {
		return nil, ErrUserRequired
	}
	if in.CourseID == "" {
		return nil, ErrCourseRequired
	}
	if !in.Role.Valid() {
		return nil, ErrInvalidRole
	}

	e := &domain.Enrollment{
		ID:       ulid.Make().String(),
		UserID:   in.UserID,
		CourseID: in.CourseID,
		Role:     in.Role,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return e, nil
}

func (s *EnrollmentsService) Unenroll(ctx context.Context, userID, courseID string) error {
	if userID == "" {
		return ErrUserRequired
	}
	if courseID == "" {
		return ErrCourseRequired
	}
	if err := s.repo.Delete(ctx, userID, courseID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s *EnrollmentsService) List(ctx context.Context, courseID, userID string, limit, offset int32) ([]*domain.Enrollment, error) {
	enrollments, err := s.repo.List(ctx, courseID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return enrollments, nil
}
