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

type StudentRepoRepo struct {
	q *sqlcgen.Queries
}

func NewStudentReposRepo(pool *pgxpool.Pool) *StudentRepoRepo {
	return &StudentRepoRepo{q: sqlcgen.New(pool)}
}

func (r *StudentRepoRepo) Register(ctx context.Context, sr *domain.StudentRepo) (*domain.StudentRepo, error) {
	row, err := r.q.RegisterStudentRepo(ctx, sqlcgen.RegisterStudentRepoParams{
		ID:               sr.ID,
		UserID:           sr.UserID,
		AssignmentID:     sr.AssignmentID,
		ProviderKind:     sr.Provider.Kind,
		ProviderInstance: sr.Provider.Instance,
		FullName:         sr.FullName,
		ExternalID:       sr.ExternalID,
		State:            string(sr.State),
		CloneUrlHttps:    sr.CloneURLHTTPS,
		CloneUrlSsh:      sr.CloneURLSSH,
	})
	if err != nil {
		return nil, fmt.Errorf("register student repo: %w", err)
	}
	return studentRepoFromRow(row), nil
}

func (r *StudentRepoRepo) Get(ctx context.Context, userID, assignmentID string) (*domain.StudentRepo, error) {
	row, err := r.q.GetStudentRepo(ctx, sqlcgen.GetStudentRepoParams{
		UserID:       userID,
		AssignmentID: assignmentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get student repo: %w", err)
	}
	return studentRepoFromRow(row), nil
}

func studentRepoFromRow(row sqlcgen.StudentRepo) *domain.StudentRepo {
	return &domain.StudentRepo{
		ID:           row.ID,
		UserID:       row.UserID,
		AssignmentID: row.AssignmentID,
		Provider: domain.ProviderRef{
			Kind:     row.ProviderKind,
			Instance: row.ProviderInstance,
		},
		FullName:      row.FullName,
		ExternalID:    row.ExternalID,
		State:         domain.ProvisioningState(row.State),
		CloneURLHTTPS: row.CloneUrlHttps,
		CloneURLSSH:   row.CloneUrlSsh,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
