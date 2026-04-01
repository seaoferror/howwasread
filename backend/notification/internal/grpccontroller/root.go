package grpccontroller

import (
	"backend/notification/internal/controller"
	pb "backend/proto"
	"context"

	"github.com/google/uuid"
)

type GRPCController struct {
	controller *controller.Controller
	pb.UnimplementedNotificationServiceServer
}

func NewGRPCController(controller *controller.Controller) *GRPCController {
	gc := GRPCController{controller: controller}

	return &gc
}

func (gc *GRPCController) NotifyMessaging(ctx context.Context, in *pb.NotifyMessagingRequest) (*pb.NotifyMessagingResponse, error) {
	roomId := uuid.UUID(in.RoomId)
	var toIds []uuid.UUID
	for _, toId := range in.ToIds {
		toIds = append(toIds, uuid.UUID(toId))
	}
	err := gc.controller.NotifyMessaging(ctx, toIds, roomId, uuid.UUID(in.FromId), in.ContentType, in.Content)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
