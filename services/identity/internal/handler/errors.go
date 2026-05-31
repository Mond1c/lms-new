package handler

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

// errCodeTable maps sentinel errors to Connect codes. It is consulted in order;
// the first errors.Is match wins, so unwrappable category errors (e.g.
// service.ErrValidation, which every service validation error wraps) catch all
// their specific variants with a single entry. Anything unmatched is Internal.
var errCodeTable = []struct {
	err  error
	code connect.Code
}{
	{domain.ErrInvalidEmail, connect.CodeInvalidArgument},
	{domain.ErrPasswordTooShort, connect.CodeInvalidArgument},
	{domain.ErrPasswordTooLong, connect.CodeInvalidArgument},
	{service.ErrValidation, connect.CodeInvalidArgument},
	{repo.ErrNotFound, connect.CodeNotFound},
	{repo.ErrEmailTaken, connect.CodeAlreadyExists},
	{repo.ErrConflict, connect.CodeAlreadyExists},
}

func toConnectErr(err error) error {
	if err == nil {
		return nil
	}
	for _, e := range errCodeTable {
		if errors.Is(err, e.err) {
			return connect.NewError(e.code, err)
		}
	}
	return connect.NewError(connect.CodeInternal, err)
}
