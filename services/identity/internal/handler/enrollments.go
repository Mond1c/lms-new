package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

func (h *Identity) Enroll(
	ctx context.Context,
	req *connect.Request[identityv1.EnrollRequest],
) (*connect.Response[identityv1.Enrollment], error) {
	m := req.Msg
	enrollment, err := h.enrollments.Enroll(ctx, service.EnrollInput{
		UserID:   m.GetUserId(),
		CourseID: m.GetCourseId(),
		Role:     roleFromProto(m.GetRole()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(enrollmentToProto(enrollment)), nil
}

func (h *Identity) Unenroll(
	ctx context.Context,
	req *connect.Request[identityv1.UnenrollRequest],
) (*connect.Response[identityv1.UnenrollResponse], error) {
	if err := h.enrollments.Unenroll(ctx, req.Msg.GetUserId(), req.Msg.GetCourseId()); err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(&identityv1.UnenrollResponse{}), nil
}

func (h *Identity) ListEnrollments(
	ctx context.Context,
	req *connect.Request[identityv1.ListEnrollmentsRequest],
) (*connect.Response[identityv1.ListEnrollmentsResponse], error) {
	limit, offset := pageParams(req.Msg.GetPage())
	enrollments, err := h.enrollments.List(ctx, req.Msg.GetCourseId(), req.Msg.GetUserId(), limit, offset)
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoEnrollments := make([]*identityv1.Enrollment, len(enrollments))
	for i, enrollment := range enrollments {
		protoEnrollments[i] = enrollmentToProto(enrollment)
	}
	return connect.NewResponse(&identityv1.ListEnrollmentsResponse{
		Enrollments: protoEnrollments,
		Page: &identityv1.PageResponse{
			NextPageToken: nextPageToken(offset, limit, int32(len(enrollments))),
		},
	}), nil
}
