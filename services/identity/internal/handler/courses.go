package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

func (h *Identity) CreateCourse(
	ctx context.Context,
	req *connect.Request[identityv1.CreateCourseRequest],
) (*connect.Response[identityv1.Course], error) {
	m := req.Msg
	course, err := h.courses.Create(ctx, service.CreateCourseInput{
		Code:         m.GetCode(),
		Title:        m.GetTitle(),
		Description:  m.GetDescription(),
		InstructorID: m.GetInstructorId(),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(courseToProto(course)), nil
}

func (h *Identity) GetCourse(
	ctx context.Context,
	req *connect.Request[identityv1.GetCourseRequest],
) (*connect.Response[identityv1.Course], error) {
	course, err := h.courses.GetByID(ctx, req.Msg.GetId())
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(courseToProto(course)), nil
}

func (h *Identity) ListCourses(
	ctx context.Context,
	req *connect.Request[identityv1.ListCoursesRequest],
) (*connect.Response[identityv1.ListCoursesResponse], error) {
	limit, offset := pageParams(req.Msg.GetPage())
	courses, err := h.courses.List(ctx, req.Msg.GetInstructorId(), limit, offset)
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoCourses := make([]*identityv1.Course, len(courses))
	for i, course := range courses {
		protoCourses[i] = courseToProto(course)
	}
	return connect.NewResponse(&identityv1.ListCoursesResponse{
		Courses: protoCourses,
		Page: &identityv1.PageResponse{
			NextPageToken: nextPageToken(offset, limit, int32(len(courses))),
		},
	}), nil
}
