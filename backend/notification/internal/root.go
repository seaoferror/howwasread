package internal

import (
	"backend/notification/internal/grpccontroller"
	"backend/notification/internal/repository"
	"backend/notification/internal/service"
	pb "backend/proto"
	"log"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
)

func NewServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	s := service.NewService(r)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panic("fail to create tcp listener at port 50051")
	}
	gc := grpccontroller.NewGRPCController(s)
	g := grpc.NewServer()
	pb.RegisterNotificationServiceServer(g, gc)
	err = g.Serve(lis)
	if err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}
