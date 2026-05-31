package service

import (
	"context"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
)

type CourseRepo interface {
	Create(ctx context.Context, c *domain.Course) error
	GetByID(ctx context.Context, id string) (*domain.Course, error)
	List(ctx context.Context, instructorID string, limit, offset int32) ([]*domain.Course, error)
}

type CoursesService struct {
	repo CourseRepo
}

var (
	ErrCourseCodeRequired  = fmt.Errorf("course code is required: %w", ErrValidation)
	ErrCourseTitleRequired = fmt.Errorf("course title is required: %w", ErrValidation)
	ErrInstructorRequired  = fmt.Errorf("instructor_id is required: %w", ErrValidation)
)

func NewCoursesService(repo CourseRepo) *CoursesService {
	return &CoursesService{repo: repo}
}

type CreateCourseInput struct {
	Code         string
	Title        string
	Description  string
	InstructorID string
}

func (s *CoursesService) Create(ctx context.Context, in CreateCourseInput) (*domain.Course, error) {
	if in.Code == "" {
		return nil, ErrCourseCodeRequired
	}
	if in.Title == "" {
		return nil, ErrCourseTitleRequired
	}
	if in.InstructorID == "" {
		return nil, ErrInstructorRequired
	}

	c := &domain.Course{
		ID:           ulid.Make().String(),
		Code:         in.Code,
		Title:        in.Title,
		Description:  in.Description,
		InstructorID: in.InstructorID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return c, nil
}

func (s *CoursesService) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	return c, nil
}

func (s *CoursesService) List(ctx context.Context, instructorID string, limit, offset int32) ([]*domain.Course, error) {
	courses, err := s.repo.List(ctx, instructorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return courses, nil
}
