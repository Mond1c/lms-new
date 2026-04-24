package handler

import (
	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
)

type Service struct {
	lmsv1connect.UnimplementedVCSServiceHandler
}

var _ lmsv1connect.VCSServiceHandler = (*Service)(nil)

func New() *Service {
	return &Service{}
}
