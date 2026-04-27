package repo

import (
	"errors"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo/sqlcgen"
	"github.com/jackc/pgx/v5/pgconn"
)

const PGDuplicateValueErrCode = "23505"

func pgTextFromPasswordHash(passwordHash domain.PasswordHash) *string {
	if passwordHash == "" {
		return nil
	}
	return new(string(passwordHash))
}

func pgTextFromString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func isUniqueViolation(err error, constraint string) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if pgErr.Code != PGDuplicateValueErrCode {
			return false
		}
		return constraint == "" || pgErr.ConstraintName == constraint
	}
	return false
}

func userFromRow(row sqlcgen.User) *domain.User {
	user := &domain.User{
		ID:          row.ID,
		Email:       domain.EmailFromTrusted(row.Email),
		DisplayName: row.DisplayName,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	if row.PasswordHash != nil {
		user.PasswordHash = domain.PasswordHash(*row.PasswordHash)
	}
	if row.TelegramID != nil {
		user.TelegramID = *row.TelegramID
	}
	return user
}
