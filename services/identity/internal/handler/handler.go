package handler

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/repo"
	"github.com/Mond1c/lms/services/identity/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Identity struct {
	lmsv1connect.UnimplementedIdentityServiceHandler
	users *service.UsersService
}

var _ lmsv1connect.IdentityServiceHandler = (*Identity)(nil)

func New(users *service.UsersService) *Identity {
	return &Identity{users: users}
}

func (h *Identity) CreateUser(
	ctx context.Context,
	req *connect.Request[identityv1.CreateUserRequest],
) (*connect.Response[identityv1.User], error) {
	m := req.Msg
	user, err := h.users.Create(ctx, service.CreateUserInput{
		Email:       m.GetEmail(),
		DisplayName: m.GetDisplayName(),
		TelegramID:  m.GetTelegramId(),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(userToProto(user)), nil
}

// GetUser TODO: maybe make separate handlers for this
func (h *Identity) GetUser(
	ctx context.Context,
	req *connect.Request[identityv1.GetUserRequest],
) (*connect.Response[identityv1.User], error) {
	var user *domain.User
	var err error
	switch req.Msg.GetIdentity().(type) {
	case *identityv1.GetUserRequest_Id:
		user, err = h.users.GetByID(ctx, req.Msg.GetId())
	case *identityv1.GetUserRequest_Email:
		user, err = h.users.GetByEmail(ctx, req.Msg.GetEmail())
	case *identityv1.GetUserRequest_VcsLogin:
		// TODO: implement me
	default:
		err = fmt.Errorf("unknown identity type") // TODO: make with errors.New
	}
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(userToProto(user)), nil
}

func userToProto(user *domain.User) *identityv1.User {
	return &identityv1.User{
		Id:          user.ID,
		Email:       user.Email.String(),
		DisplayName: user.DisplayName,
		TelegramId:  user.TelegramID,
		Audit: &identityv1.AuditFields{
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		// TODO: add vcs info here
	}
}

// TODO: this is bad function, maybe store something like mapping error to code?
func toConnectErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrPasswordTooShort),
		errors.Is(err, domain.ErrPasswordTooLong),
		errors.Is(err, service.ErrDisplayNameRequired),
		errors.Is(err, service.ErrInvalidID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, repo.ErrNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, repo.ErrEmailTaken):
		return connect.NewError(connect.CodeAlreadyExists, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}
