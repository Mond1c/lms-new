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

type AssignmentRepo struct {
	q *sqlcgen.Queries
}

func NewAssignmentsRepo(pool *pgxpool.Pool) *AssignmentRepo {
	return &AssignmentRepo{q: sqlcgen.New(pool)}
}

func (r *AssignmentRepo) Create(ctx context.Context, a *domain.Assignment) error {
	row, err := r.q.CreateAssignment(ctx, sqlcgen.CreateAssignmentParams{
		ID:                      a.ID,
		CourseID:                a.CourseID,
		Slug:                    a.Slug,
		Title:                   a.Title,
		DescriptionMarkdown:     a.DescriptionMarkdown,
		Deadline:                pgTimestamp(a.Deadline),
		HardDeadline:            pgTimestamp(a.HardDeadline),
		MaxScore:                a.MaxScore,
		TemplateRepo:            a.TemplateRepo,
		RepoNamingPattern:       a.RepoNamingPattern,
		AutoRequestReviewOnPass: a.AutoRequestReviewOnPass,
		RequiresDefense:         a.RequiresDefense,
		WeightTests:             a.GradingPolicy.WeightTests,
		WeightQuality:           a.GradingPolicy.WeightQuality,
		DefenceMultiplier:       a.GradingPolicy.DefenceMultiplier,
		CustomFormula:           a.GradingPolicy.CustomFormula,
		Runner:                  string(a.Runner),
	})
	if err != nil {
		if isUniqueViolation(err, "") {
			return ErrConflict
		}
		return fmt.Errorf("create assignment: %w", err)
	}
	*a = *assignmentFromRow(row)
	return nil
}

func (r *AssignmentRepo) GetByID(ctx context.Context, id string) (*domain.Assignment, error) {
	row, err := r.q.GetAssignmentById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get assignment by id: %w", err)
	}
	return assignmentFromRow(row), nil
}

func (r *AssignmentRepo) List(ctx context.Context, courseID string, limit, offset int32) ([]*domain.Assignment, error) {
	rows, err := r.q.ListAssignments(ctx, sqlcgen.ListAssignmentsParams{
		CourseID: pgTextFromString(courseID),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}
	results := make([]*domain.Assignment, 0, len(rows))
	for _, row := range rows {
		results = append(results, assignmentFromRow(row))
	}
	return results, nil
}

func assignmentFromRow(row sqlcgen.Assignment) *domain.Assignment {
	return &domain.Assignment{
		ID:                      row.ID,
		CourseID:                row.CourseID,
		Slug:                    row.Slug,
		Title:                   row.Title,
		DescriptionMarkdown:     row.DescriptionMarkdown,
		Deadline:                timeFromPg(row.Deadline),
		HardDeadline:            timeFromPg(row.HardDeadline),
		MaxScore:                row.MaxScore,
		TemplateRepo:            row.TemplateRepo,
		RepoNamingPattern:       row.RepoNamingPattern,
		AutoRequestReviewOnPass: row.AutoRequestReviewOnPass,
		RequiresDefense:         row.RequiresDefense,
		GradingPolicy: domain.GradingPolicy{
			WeightTests:       row.WeightTests,
			WeightQuality:     row.WeightQuality,
			DefenceMultiplier: row.DefenceMultiplier,
			CustomFormula:     row.CustomFormula,
		},
		Runner:    domain.RunnerKind(row.Runner),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
