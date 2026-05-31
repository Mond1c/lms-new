package repo

import (
	"context"
	"testing"
	"time"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestAssignmentRepo_CreateGetList(t *testing.T) {
	truncate(t)
	ctx := context.Background()
	r := NewAssignmentsRepo(testDB)

	deadline := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Microsecond)
	a := &domain.Assignment{
		ID:                  ulid.Make().String(),
		CourseID:            "course-1",
		Slug:                "hw1",
		Title:               "Homework 1",
		DescriptionMarkdown: "# hw",
		Deadline:            &deadline,
		MaxScore:            100,
		RequiresDefense:     true,
		GradingPolicy: domain.GradingPolicy{
			WeightTests:       0.6,
			WeightQuality:     0.4,
			DefenceMultiplier: true,
			CustomFormula:     "min(tests, quality) * defence",
		},
		Runner: domain.RunnerSelfHosted,
	}
	require.NoError(t, r.Create(ctx, a))
	require.False(t, a.CreatedAt.IsZero())

	t.Run("get round-trips policy, runner, deadline", func(t *testing.T) {
		got, err := r.GetByID(ctx, a.ID)
		require.NoError(t, err)
		require.Equal(t, "hw1", got.Slug)
		require.True(t, got.RequiresDefense)
		require.Equal(t, domain.RunnerSelfHosted, got.Runner)
		require.Equal(t, 0.6, got.GradingPolicy.WeightTests)
		require.Equal(t, "min(tests, quality) * defence", got.GradingPolicy.CustomFormula)
		require.NotNil(t, got.Deadline)
		require.WithinDuration(t, deadline, *got.Deadline, time.Second)
		require.Nil(t, got.HardDeadline)
	})

	t.Run("duplicate (course, slug) conflicts", func(t *testing.T) {
		dup := &domain.Assignment{ID: ulid.Make().String(), CourseID: "course-1", Slug: "hw1", Title: "x", Runner: domain.RunnerExternalCI}
		require.ErrorIs(t, r.Create(ctx, dup), ErrConflict)
	})

	t.Run("list by course", func(t *testing.T) {
		got, err := r.List(ctx, "course-1", 10, 0)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})
}
