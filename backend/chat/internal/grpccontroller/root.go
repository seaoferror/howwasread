package grpccontroller

import (
	"backend/chat/internal/controller"
	pb "backend/proto"
)

type GRPCController struct {
	controller *controller.Controller
	pb.UnimplementedSignalServiceServer
}

func NewGRPCController(controller *controller.Controller) *GRPCController {
	gc := GRPCController{controller: controller}

	return &gc
}
