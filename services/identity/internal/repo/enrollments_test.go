package repo

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func newEnrollment(userID, courseID string, role domain.Role) *domain.Enrollment {
	return &domain.Enrollment{
		ID:       ulid.Make().String(),
		UserID:   userID,
		CourseID: courseID,
		Role:     role,
	}
}

func TestEnrollmentRepo_CreateListDelete(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewEnrollmentsRepo(testDB)

	e := newEnrollment("user-1", "course-1", domain.RoleStudent)
	require.NoError(t, r.Create(ctx, e))
	require.False(t, e.EnrolledAt.IsZero())

	t.Run("duplicate enrollment conflicts", func(t *testing.T) {
		dup := newEnrollment("user-1", "course-1", domain.RoleInstructor)
		require.ErrorIs(t, r.Create(ctx, dup), ErrConflict)
	})

	t.Run("list by course preserves role", func(t *testing.T) {
		got, err := r.List(ctx, "course-1", "", 10, 0)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, domain.RoleStudent, got[0].Role)
	})

	t.Run("list by user", func(t *testing.T) {
		require.NoError(t, r.Create(ctx, newEnrollment("user-1", "course-2", domain.RoleInstructor)))
		got, err := r.List(ctx, "", "user-1", 10, 0)
		require.NoError(t, err)
		require.Len(t, got, 2)
	})

	t.Run("delete then delete-missing", func(t *testing.T) {
		require.NoError(t, r.Delete(ctx, "user-1", "course-1"))
		require.ErrorIs(t, r.Delete(ctx, "user-1", "course-1"), ErrNotFound)
	})
}
