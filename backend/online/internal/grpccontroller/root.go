package grpccontroller

import (
	"backend/online/internal/controller"
	pb "backend/proto"
	"context"

	"github.com/google/uuid"
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
	fromId := uuid.UUID(in.FromId)
	toId := uuid.UUID(in.ToId)
	err := gc.controller.RelaySignal(ctx, fromId, toId, in.Signal)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
