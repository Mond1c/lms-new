package service

import (
	"context"
	"testing"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/stretchr/testify/require"
)

type mockAssignmentRepo struct {
	createFn func(ctx context.Context, a *domain.Assignment) error
}

func (m *mockAssignmentRepo) Create(ctx context.Context, a *domain.Assignment) error {
	return m.createFn(ctx, a)
}

func (m *mockAssignmentRepo) GetByID(ctx context.Context, id string) (*domain.Assignment, error) {
	return nil, nil
}

func (m *mockAssignmentRepo) List(ctx context.Context, courseID string, limit, offset int32) ([]*domain.Assignment, error) {
	return nil, nil
}

func validAssignmentInput() CreateAssignmentInput {
	return CreateAssignmentInput{CourseID: "c1", Slug: "hw1", Title: "Homework 1"}
}

func TestAssignmentsService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("defaults runner and grading weights", func(t *testing.T) {
		var captured *domain.Assignment
		svc := NewAssignmentsService(&mockAssignmentRepo{
			createFn: func(_ context.Context, a *domain.Assignment) error { captured = a; return nil },
		})

		got, err := svc.Create(ctx, validAssignmentInput())
		require.NoError(t, err)
		require.NotEmpty(t, got.ID)
		require.Equal(t, domain.RunnerExternalCI, got.Runner, "empty runner should default to external CI")
		require.Equal(t, defaultWeightTests, got.GradingPolicy.WeightTests)
		require.Equal(t, defaultWeightQuality, got.GradingPolicy.WeightQuality)
		require.Same(t, captured, got)
	})

	t.Run("keeps an explicit grading policy", func(t *testing.T) {
		svc := NewAssignmentsService(&mockAssignmentRepo{
			createFn: func(_ context.Context, a *domain.Assignment) error { return nil },
		})
		in := validAssignmentInput()
		in.GradingPolicy = domain.GradingPolicy{WeightTests: 0.9, WeightQuality: 0.1}
		in.Runner = domain.RunnerSelfHosted

		got, err := svc.Create(ctx, in)
		require.NoError(t, err)
		require.Equal(t, 0.9, got.GradingPolicy.WeightTests)
		require.Equal(t, domain.RunnerSelfHosted, got.Runner)
	})

	t.Run("invalid runner", func(t *testing.T) {
		svc := NewAssignmentsService(&mockAssignmentRepo{})
		in := validAssignmentInput()
		in.Runner = "bogus"
		_, err := svc.Create(ctx, in)
		require.ErrorIs(t, err, ErrInvalidRunner)
		require.ErrorIs(t, err, ErrValidation)
	})

	t.Run("required fields", func(t *testing.T) {
		svc := NewAssignmentsService(&mockAssignmentRepo{})
		for _, tc := range []struct {
			name string
			in   CreateAssignmentInput
			want error
		}{
			{"no course", CreateAssignmentInput{Slug: "s", Title: "t"}, ErrAssignmentCourseRequired},
			{"no slug", CreateAssignmentInput{CourseID: "c", Title: "t"}, ErrAssignmentSlugRequired},
			{"no title", CreateAssignmentInput{CourseID: "c", Slug: "s"}, ErrAssignmentTitleRequired},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := svc.Create(ctx, tc.in)
				require.ErrorIs(t, err, tc.want)
			})
		}
	})
}
