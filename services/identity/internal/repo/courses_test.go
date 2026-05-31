package repo

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func newCourse(code, instructor string) *domain.Course {
	return &domain.Course{
		ID:           ulid.Make().String(),
		Code:         code,
		Title:        "Course " + code,
		Description:  "desc",
		InstructorID: instructor,
	}
}

func TestCourseRepo_CreateGetList(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewCoursesRepo(testDB)

	c := newCourse("CS101", "inst-1")
	require.NoError(t, r.Create(ctx, c))
	require.False(t, c.CreatedAt.IsZero(), "Create should backfill audit fields")

	t.Run("get", func(t *testing.T) {
		got, err := r.GetByID(ctx, c.ID)
		require.NoError(t, err)
		require.Equal(t, "CS101", got.Code)
		require.Equal(t, "inst-1", got.InstructorID)
		require.Nil(t, got.VCS, "no VCS binding set on create")
	})

	t.Run("duplicate code conflicts", func(t *testing.T) {
		dup := newCourse("CS101", "inst-2")
		require.ErrorIs(t, r.Create(ctx, dup), ErrConflict)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := r.GetByID(ctx, "nope")
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("list filters by instructor", func(t *testing.T) {
		require.NoError(t, r.Create(ctx, newCourse("CS201", "inst-2")))

		mine, err := r.List(ctx, "inst-1", 10, 0)
		require.NoError(t, err)
		require.Len(t, mine, 1)

		all, err := r.List(ctx, "", 10, 0)
		require.NoError(t, err)
		require.Len(t, all, 2)
	})
}
