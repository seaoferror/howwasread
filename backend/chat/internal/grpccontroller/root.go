package grpccontroller

import (
	"backend/chat/internal/controller"
	pb "backend/proto"
	"context"

	"github.com/google/uuid"
)

type GRPCController struct {
	controller *controller.Controller
	pb.UnimplementedMessagingServiceServer
}

func NewGRPCController(controller *controller.Controller) *GRPCController {
	gc := GRPCController{controller: controller}

	return &gc
}

func (gc *GRPCController) RelayMessaging(ctx context.Context, in *pb.RelayMessagingRequest) (*pb.RelayMessagingResponse, error) {
	err := gc.controller.RelayMessaging(ctx, uuid.UUID(in.Id), in.ToIds, uuid.UUID(in.RoomId), uuid.UUID(in.FromId), in.ContentType, in.Content)
	return nil, err

}
