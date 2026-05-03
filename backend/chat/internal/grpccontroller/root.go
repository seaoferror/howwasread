package grpccontroller

import (
	"backend/chat/internal/controller"
	"backend/common/proto"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type GRPCController struct {
	controller *controller.Controller
	proto.UnimplementedMessagingServiceServer
}

func NewGRPCController(controller *controller.Controller) *GRPCController {
	gc := GRPCController{controller: controller}

	return &gc
}

func (gc *GRPCController) RelayMessaging(ctx context.Context, in *proto.RelayMessagingRequest) (*proto.RelayMessagingResponse, error) {
	slog.Info("relay message GRPC request incoming...")
	pushToIds, _ := gc.controller.RelayMessaging(ctx, uuid.UUID(in.Id), in.ToIds, uuid.UUID(in.RoomId), uuid.UUID(in.FromId), in.ContentType, in.Contents)
	return &proto.RelayMessagingResponse{
		PushToIds: pushToIds,
	}, nil
}
