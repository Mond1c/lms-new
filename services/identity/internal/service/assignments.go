package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
)

type AssignmentRepo interface {
	Create(ctx context.Context, a *domain.Assignment) error
	GetByID(ctx context.Context, id string) (*domain.Assignment, error)
	List(ctx context.Context, courseID string, limit, offset int32) ([]*domain.Assignment, error)
}

type AssignmentsService struct {
	repo AssignmentRepo
}

var (
	ErrAssignmentCourseRequired = fmt.Errorf("course_id is required: %w", ErrValidation)
	ErrAssignmentSlugRequired   = fmt.Errorf("slug is required: %w", ErrValidation)
	ErrAssignmentTitleRequired  = fmt.Errorf("title is required: %w", ErrValidation)
	ErrInvalidRunner            = fmt.Errorf("invalid runner: %w", ErrValidation)
)

// Default grading weights when an assignment is created without an explicit
// policy, matching the reference course (architecture §4).
const (
	defaultWeightTests   = 0.7
	defaultWeightQuality = 0.3
)

func NewAssignmentsService(repo AssignmentRepo) *AssignmentsService {
	return &AssignmentsService{repo: repo}
}

type CreateAssignmentInput struct {
	CourseID                string
	Slug                    string
	Title                   string
	DescriptionMarkdown     string
	TemplateRepo            string
	RepoNamingPattern       string
	Deadline                *time.Time
	HardDeadline            *time.Time
	MaxScore                int32
	AutoRequestReviewOnPass bool
	RequiresDefense         bool
	GradingPolicy           domain.GradingPolicy
	Runner                  domain.RunnerKind
}

func (s *AssignmentsService) Create(ctx context.Context, in CreateAssignmentInput) (*domain.Assignment, error) {
	if in.CourseID == "" {
		return nil, ErrAssignmentCourseRequired
	}
	if in.Slug == "" {
		return nil, ErrAssignmentSlugRequired
	}
	if in.Title == "" {
		return nil, ErrAssignmentTitleRequired
	}

	runner := in.Runner
	if runner == "" {
		runner = domain.RunnerExternalCI
	}
	if !runner.Valid() {
		return nil, ErrInvalidRunner
	}

	policy := in.GradingPolicy
	if policy.WeightTests == 0 && policy.WeightQuality == 0 && policy.CustomFormula == "" {
		policy.WeightTests = defaultWeightTests
		policy.WeightQuality = defaultWeightQuality
	}

	a := &domain.Assignment{
		ID:                      ulid.Make().String(),
		CourseID:                in.CourseID,
		Slug:                    in.Slug,
		Title:                   in.Title,
		DescriptionMarkdown:     in.DescriptionMarkdown,
		Deadline:                in.Deadline,
		HardDeadline:            in.HardDeadline,
		MaxScore:                in.MaxScore,
		TemplateRepo:            in.TemplateRepo,
		RepoNamingPattern:       in.RepoNamingPattern,
		AutoRequestReviewOnPass: in.AutoRequestReviewOnPass,
		RequiresDefense:         in.RequiresDefense,
		GradingPolicy:           policy,
		Runner:                  runner,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return a, nil
}

func (s *AssignmentsService) GetByID(ctx context.Context, id string) (*domain.Assignment, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return a, nil
}

func (s *AssignmentsService) List(ctx context.Context, courseID string, limit, offset int32) ([]*domain.Assignment, error) {
	assignments, err := s.repo.List(ctx, courseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return assignments, nil
}
