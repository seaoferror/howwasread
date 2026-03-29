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
	roomId := uuid.UUID(in.RoomId)
	if len(in.RoomId) == 0 {
		roomId = uuid.Nil
	}
	err := gc.controller.RelayMessaging(ctx, uuid.UUID(in.ToId), roomId, uuid.UUID(in.FromId), in.ContentType, in.Content)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
