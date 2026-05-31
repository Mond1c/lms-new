package repo

import (
	"errors"
	"time"

	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo/sqlcgen"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func strFromPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func ptrBool(b bool) *bool { return &b }

func pgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

func timeFromPg(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
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
