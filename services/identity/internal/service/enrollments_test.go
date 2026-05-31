package service

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo"
	"github.com/stretchr/testify/require"
)

type mockEnrollmentRepo struct {
	createFn func(ctx context.Context, e *domain.Enrollment) error
	deleteFn func(ctx context.Context, userID, courseID string) error
}

func (m *mockEnrollmentRepo) Create(ctx context.Context, e *domain.Enrollment) error {
	return m.createFn(ctx, e)
}

func (m *mockEnrollmentRepo) Delete(ctx context.Context, userID, courseID string) error {
	return m.deleteFn(ctx, userID, courseID)
}

func (m *mockEnrollmentRepo) List(ctx context.Context, courseID, userID string, limit, offset int32) ([]*domain.Enrollment, error) {
	return nil, nil
}

func TestEnrollmentsService_Enroll(t *testing.T) {
	ctx := context.Background()

	t.Run("valid", func(t *testing.T) {
		var captured *domain.Enrollment
		svc := NewEnrollmentsService(&mockEnrollmentRepo{
			createFn: func(_ context.Context, e *domain.Enrollment) error { captured = e; return nil },
		})
		got, err := svc.Enroll(ctx, EnrollInput{UserID: "u1", CourseID: "c1", Role: domain.RoleStudent})
		require.NoError(t, err)
		require.NotEmpty(t, got.ID)
		require.Equal(t, domain.RoleStudent, captured.Role)
	})

	t.Run("invalid role rejected", func(t *testing.T) {
		svc := NewEnrollmentsService(&mockEnrollmentRepo{
			createFn: func(_ context.Context, _ *domain.Enrollment) error {
				t.Fatal("repo.Create must not be called for an invalid role")
				return nil
			},
		})
		_, err := svc.Enroll(ctx, EnrollInput{UserID: "u1", CourseID: "c1", Role: ""})
		require.ErrorIs(t, err, ErrInvalidRole)
	})

	t.Run("missing ids", func(t *testing.T) {
		svc := NewEnrollmentsService(&mockEnrollmentRepo{})
		_, err := svc.Enroll(ctx, EnrollInput{CourseID: "c1", Role: domain.RoleStudent})
		require.ErrorIs(t, err, ErrUserRequired)
		_, err = svc.Enroll(ctx, EnrollInput{UserID: "u1", Role: domain.RoleStudent})
		require.ErrorIs(t, err, ErrCourseRequired)
	})
}

func TestEnrollmentsService_Unenroll(t *testing.T) {
	ctx := context.Background()

	t.Run("propagates not found", func(t *testing.T) {
		svc := NewEnrollmentsService(&mockEnrollmentRepo{
			deleteFn: func(_ context.Context, _, _ string) error { return repo.ErrNotFound },
		})
		err := svc.Unenroll(ctx, "u1", "c1")
		require.ErrorIs(t, err, repo.ErrNotFound)
	})
}
