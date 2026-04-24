package handler

import (
	"github.com/Mond1c/lms/gen/go/lms/v1/lmsv1connect"
)

type Service struct {
	lmsv1connect.UnimplementedGatewayServiceHandler
}

var _ lmsv1connect.GatewayServiceHandler = (*Service)(nil)

func New() *Service {
	return &Service{}
}
