package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

func (h *Identity) RegisterStudentRepo(
	ctx context.Context,
	req *connect.Request[identityv1.RegisterStudentRepoRequest],
) (*connect.Response[identityv1.StudentRepo], error) {
	m := req.Msg
	repo, err := h.studentRepos.Register(ctx, service.RegisterStudentRepoInput{
		UserID:        m.GetUserId(),
		AssignmentID:  m.GetAssignmentId(),
		Provider:      providerFromProto(m.GetProvider()),
		FullName:      m.GetFullName(),
		ExternalID:    m.GetExternalId(),
		CloneURLHTTPS: m.GetCloneUrlHttps(),
		CloneURLSSH:   m.GetCloneUrlSsh(),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(studentRepoToProto(repo)), nil
}

func (h *Identity) GetStudentRepo(
	ctx context.Context,
	req *connect.Request[identityv1.GetStudentRepoRequest],
) (*connect.Response[identityv1.StudentRepo], error) {
	repo, err := h.studentRepos.Get(ctx, req.Msg.GetUserId(), req.Msg.GetAssignmentId())
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(studentRepoToProto(repo)), nil
}
