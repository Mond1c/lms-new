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

type CourseRepo struct {
	q *sqlcgen.Queries
}

func NewCoursesRepo(pool *pgxpool.Pool) *CourseRepo {
	return &CourseRepo{q: sqlcgen.New(pool)}
}

func (r *CourseRepo) Create(ctx context.Context, c *domain.Course) error {
	row, err := r.q.CreateCourse(ctx, sqlcgen.CreateCourseParams{
		ID:           c.ID,
		Code:         c.Code,
		Title:        c.Title,
		Description:  c.Description,
		InstructorID: c.InstructorID,
	})
	if err != nil {
		if isUniqueViolation(err, "") {
			return ErrConflict
		}
		return fmt.Errorf("create course: %w", err)
	}
	*c = *courseFromRow(row)
	return nil
}

func (r *CourseRepo) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	row, err := r.q.GetCourseById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get course by id: %w", err)
	}
	return courseFromRow(row), nil
}

func (r *CourseRepo) List(ctx context.Context, instructorID string, limit, offset int32) ([]*domain.Course, error) {
	rows, err := r.q.ListCourses(ctx, sqlcgen.ListCoursesParams{
		InstructorID: pgTextFromString(instructorID),
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	results := make([]*domain.Course, 0, len(rows))
	for _, row := range rows {
		results = append(results, courseFromRow(row))
	}
	return results, nil
}

func courseFromRow(row sqlcgen.Course) *domain.Course {
	c := &domain.Course{
		ID:           row.ID,
		Code:         row.Code,
		Title:        row.Title,
		Description:  row.Description,
		InstructorID: row.InstructorID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	if row.VcsProviderKind != nil {
		c.VCS = &domain.VCSBinding{
			Provider: domain.ProviderRef{
				Kind:     *row.VcsProviderKind,
				Instance: strFromPtr(row.VcsProviderInstance),
			},
			TargetOrg:      strFromPtr(row.VcsTargetOrg),
			StudentTeam:    strFromPtr(row.VcsStudentTeam),
			ReviewerTeam:   strFromPtr(row.VcsReviewerTeam),
			ReviewerLogins: row.VcsReviewerLogins,
		}
	}
	return c
}
