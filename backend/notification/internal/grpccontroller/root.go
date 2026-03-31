package grpccontroller

import (
	"backend/notification/internal/service"
	pb "backend/proto"
	"context"

	"github.com/google/uuid"
)

type GRPCController struct {
	service *service.Service
	pb.UnimplementedNotificationServiceServer
}

func NewGRPCController(service *service.Service) *GRPCController {
	gc := GRPCController{service: service}

	return &gc
}

func (gc *GRPCController) NotifyMessaging(ctx context.Context, in *pb.NotifyMessagingRequest) (*pb.NotifyMessagingResponse, error) {
	roomId := uuid.UUID(in.RoomId)
	if len(in.RoomId) == 0 {
		roomId = uuid.Nil
	}
	var toIds []uuid.UUID
	for _, toId := range in.ToIds {
		toIds = append(toIds, uuid.UUID(toId))
	}
	err := gc.service.NotifyMessaging(ctx, toIds, roomId, uuid.UUID(in.FromId), in.ContentType, in.Content)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
