package grpccontroller

import (
	"backend/onlineconversation/internal/controller"
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
	toIds := make([]uuid.UUID, len(in.ToIds))
	for _, toId := range in.ToIds {
		toIds = append(toIds, uuid.UUID(toId))
	}
	gc.controller.RelaySignal(ctx, toIds, fromId, in.Signal)
	return nil, nil
}
