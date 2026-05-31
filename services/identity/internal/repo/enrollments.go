package repo

import (
	"context"
	"fmt"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo/sqlcgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepo struct {
	q *sqlcgen.Queries
}

func NewEnrollmentsRepo(pool *pgxpool.Pool) *EnrollmentRepo {
	return &EnrollmentRepo{q: sqlcgen.New(pool)}
}

func (r *EnrollmentRepo) Create(ctx context.Context, e *domain.Enrollment) error {
	row, err := r.q.CreateEnrollment(ctx, sqlcgen.CreateEnrollmentParams{
		ID:       e.ID,
		UserID:   e.UserID,
		CourseID: e.CourseID,
		Role:     string(e.Role),
	})
	if err != nil {
		if isUniqueViolation(err, "") {
			return ErrConflict
		}
		return fmt.Errorf("create enrollment: %w", err)
	}
	*e = *enrollmentFromRow(row)
	return nil
}

func (r *EnrollmentRepo) Delete(ctx context.Context, userID, courseID string) error {
	rows, err := r.q.DeleteEnrollment(ctx, sqlcgen.DeleteEnrollmentParams{
		UserID:   userID,
		CourseID: courseID,
	})
	if err != nil {
		return fmt.Errorf("delete enrollment: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *EnrollmentRepo) List(ctx context.Context, courseID, userID string, limit, offset int32) ([]*domain.Enrollment, error) {
	rows, err := r.q.ListEnrollments(ctx, sqlcgen.ListEnrollmentsParams{
		CourseID: pgTextFromString(courseID),
		UserID:   pgTextFromString(userID),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list enrollments: %w", err)
	}
	results := make([]*domain.Enrollment, 0, len(rows))
	for _, row := range rows {
		results = append(results, enrollmentFromRow(row))
	}
	return results, nil
}

func enrollmentFromRow(row sqlcgen.Enrollment) *domain.Enrollment {
	return &domain.Enrollment{
		ID:         row.ID,
		UserID:     row.UserID,
		CourseID:   row.CourseID,
		Role:       domain.Role(row.Role),
		EnrolledAt: row.EnrolledAt.Time,
	}
}
