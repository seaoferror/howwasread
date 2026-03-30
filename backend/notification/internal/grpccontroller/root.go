package grpccontroller

import (
	"backend/notification/internal/service"
)

type GRPCController struct {
	service *service.Service
}

func NewGRPCController(service *service.Service) *GRPCController {
	gc := GRPCController{service: service}

	return &gc
}
