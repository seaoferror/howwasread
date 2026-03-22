package grpccontroller

import (
	"backend/online/server/controller"
	pb "backend/proto"
	"context"
)

type GRPCController struct {
	pb.UnimplementedSignalServiceServer
	controller *controller.Controller
}

func NewGRPCController(controller *controller.Controller) *GRPCController {
	gc := GRPCController{controller: controller}

	return &gc
}

func (gc *GRPCController) RelaySignal(ctx context.Context, in *pb.RelaySignalRequest) (*pb.RelaySignalResponse, error) {
	err := gc.controller.RelaySignal(ctx, in.FromId, in.ToId, in.Signal)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
