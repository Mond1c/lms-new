package handler

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo"
	"github.com/Mond1c/lms/services/identity/internal/service"
	"github.com/stretchr/testify/require"
)

func TestToConnectErr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode connect.Code
	}{
		{"nil", nil, connect.Code(0)},
		{"invalid email", domain.ErrInvalidEmail, connect.CodeInvalidArgument},
		{"short password", domain.ErrPasswordTooShort, connect.CodeInvalidArgument},
		{"display name required", service.ErrDisplayNameRequired, connect.CodeInvalidArgument},
		{"not found", repo.ErrNotFound, connect.CodeNotFound},
		{"email taken", repo.ErrEmailTaken, connect.CodeAlreadyExists},
		{"unknown", errors.New("random"), connect.CodeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toConnectErr(tt.err)
			if tt.err == nil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.wantCode, connect.CodeOf(got))
		})
	}
}

func TestUserToProto(t *testing.T) {
	email, _ := domain.NewEmail("foo@bar.com")
	u := &domain.User{
		ID:          "01HX",
		Email:       email,
		DisplayName: "Foo",
		TelegramID:  "@foo",
	}
	p := userToProto(u)
	require.Equal(t, "01HX", p.Id)
	require.Equal(t, "foo@bar.com", p.Email)
	require.Equal(t, "Foo", p.DisplayName)
	require.Equal(t, "@foo", p.TelegramId)
}
