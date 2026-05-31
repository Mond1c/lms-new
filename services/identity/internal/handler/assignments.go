package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

func (h *Identity) CreateAssignment(
	ctx context.Context,
	req *connect.Request[identityv1.CreateAssignmentRequest],
) (*connect.Response[identityv1.Assignment], error) {
	m := req.Msg
	assignment, err := h.assignments.Create(ctx, service.CreateAssignmentInput{
		CourseID:                m.GetCourseId(),
		Slug:                    m.GetSlug(),
		Title:                   m.GetTitle(),
		DescriptionMarkdown:     m.GetDescriptionMarkdown(),
		TemplateRepo:            m.GetTemplateRepo(),
		RepoNamingPattern:       m.GetRepoNamingPattern(),
		Deadline:                tsFromProto(m.GetDeadline()),
		HardDeadline:            tsFromProto(m.GetHardDeadline()),
		MaxScore:                m.GetMaxScore(),
		AutoRequestReviewOnPass: m.GetAutoRequestReviewOnPass(),
		RequiresDefense:         m.GetRequiresDefense(),
		GradingPolicy:           gradingPolicyFromProto(m.GetGradingPolicy()),
		Runner:                  runnerFromProto(m.GetRunner()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(assignmentToProto(assignment)), nil
}

func (h *Identity) GetAssignment(
	ctx context.Context,
	req *connect.Request[identityv1.GetAssignmentRequest],
) (*connect.Response[identityv1.Assignment], error) {
	assignment, err := h.assignments.GetByID(ctx, req.Msg.GetId())
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(assignmentToProto(assignment)), nil
}

func (h *Identity) ListAssignments(
	ctx context.Context,
	req *connect.Request[identityv1.ListAssignmentsRequest],
) (*connect.Response[identityv1.ListAssignmentsResponse], error) {
	limit, offset := pageParams(req.Msg.GetPage())
	assignments, err := h.assignments.List(ctx, req.Msg.GetCourseId(), limit, offset)
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoAssignments := make([]*identityv1.Assignment, len(assignments))
	for i, assignment := range assignments {
		protoAssignments[i] = assignmentToProto(assignment)
	}
	return connect.NewResponse(&identityv1.ListAssignmentsResponse{
		Assignments: protoAssignments,
		Page: &identityv1.PageResponse{
			NextPageToken: nextPageToken(offset, limit, int32(len(assignments))),
		},
	}), nil
}
