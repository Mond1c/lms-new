package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
	"github.com/Mond1c/lms/services/identity/internal/domain"
	"github.com/Mond1c/lms/services/identity/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Identity struct {
	lmsv1connect.UnimplementedIdentityServiceHandler
	users         *service.UsersService
	courses       *service.CoursesService
	enrollments   *service.EnrollmentsService
	assignments   *service.AssignmentsService
	vcsIdentities *service.VCSIdentitiesService
	studentRepos  *service.StudentReposService
}

var _ lmsv1connect.IdentityServiceHandler = (*Identity)(nil)

func New(
	users *service.UsersService,
	courses *service.CoursesService,
	enrollments *service.EnrollmentsService,
	assignments *service.AssignmentsService,
	vcsIdentities *service.VCSIdentitiesService,
	studentRepos *service.StudentReposService,
) *Identity {
	return &Identity{
		users:         users,
		courses:       courses,
		enrollments:   enrollments,
		assignments:   assignments,
		vcsIdentities: vcsIdentities,
		studentRepos:  studentRepos,
	}
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
	limit, offset := pageParams(req.Msg.GetPage())
	users, err := h.users.List(ctx, limit, offset)
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoUsers := make([]*identityv1.User, len(users))
	for i, user := range users {
		protoUsers[i] = userToProto(user)
	}

	response := &identityv1.ListUsersResponse{
		Users: protoUsers,
		Page: &identityv1.PageResponse{
			NextPageToken: nextPageToken(offset, limit, int32(len(users))),
		},
	}

	return connect.NewResponse(response), nil
}

func (h *Identity) UpdateUser(
	ctx context.Context,
	req *connect.Request[identityv1.UpdateUserRequest],
) (*connect.Response[identityv1.User], error) {
	user, err := h.users.Update(ctx, domain.UserUpdate{
		ID:          req.Msg.GetId(),
		DisplayName: req.Msg.DisplayName,
		TelegramID:  req.Msg.TelegramId,
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
