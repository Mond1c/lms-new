package handler

import (
	"context"
	"errors"

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

func getHelper(
	ctx context.Context,
	getter func(context.Context, string) (*domain.User, error),
	key string,
) (*connect.Response[identityv1.User], error) {
	user, err := getter(ctx, key)
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(userToProto(user)), nil
}

func (h *Identity) GetUserById(
	ctx context.Context,
	req *connect.Request[identityv1.GetUserRequestById],
) (*connect.Response[identityv1.User], error) {
	return getHelper(ctx, h.users.GetByID, req.Msg.GetId())
}

func (h *Identity) GetUserByEmail(
	ctx context.Context,
	req *connect.Request[identityv1.GetUserRequestByEmail],
) (*connect.Response[identityv1.User], error) {
	return getHelper(ctx, h.users.GetByEmail, req.Msg.GetEmail())
}

func (h *Identity) ListUsers(
	ctx context.Context,
	req *connect.Request[identityv1.ListUsersRequest],
) (*connect.Response[identityv1.ListUsersResponse], error) {
	users, err := h.users.List(ctx, req.Msg.GetPage().GetPageSize(), req.Msg.GetPage().GetPageToken())
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoUsers := make([]*identityv1.User, len(users))
	for i, user := range users {
		protoUsers[i] = userToProto(user)
	}

	// TODO: I do not like this
	response := &identityv1.ListUsersResponse{
		Users: protoUsers,
		Page: &identityv1.PageResponse{
			NextPageToken: req.Msg.GetPage().GetPageToken() + 1,
		},
	}

	return connect.NewResponse(response), nil
}

func (h *Identity) UpdateUser(
	ctx context.Context,
	req *connect.Request[identityv1.UpdateUserRequest],
) (*connect.Response[identityv1.User], error) {
	user, err := h.users.Update(ctx, &domain.User{
		ID:          req.Msg.GetId(),
		DisplayName: req.Msg.GetDisplayName(),
		TelegramID:  req.Msg.GetTelegramId(),
	})
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
