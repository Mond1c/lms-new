package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
)

type Identity struct {
	lmsv1connect.UnimplementedIdentityServiceHandler
}

var _ lmsv1connect.IdentityServiceHandler = (*Identity)(nil)

func New() *Identity {
	return &Identity{}
}

func (h *Identity) CreateUser(
	ctx context.Context,
	req *connect.Request[identityv1.CreateUserRequest]) (*connect.Response[identityv1.User], error) {
	return nil, connect.NewError(
		connect.CodeUnimplemented,
		errors.New("CreateUser: not implemented"))
}
