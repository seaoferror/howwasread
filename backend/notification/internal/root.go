package internal

import (
	"backend/notification/internal/controller"
	"backend/notification/internal/grpccontroller"
	"backend/notification/internal/repository"
	"backend/notification/internal/service"
	pb "backend/proto"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
)

func NewServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	s := service.NewService(r)

	mux := http.NewServeMux()

	c := controller.NewController(s, mux)

	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			slog.Error("fail to listen and serve http",
				"err", err)
			return
		}
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panic("fail to create tcp listener at port 50051")
	}
	gc := grpccontroller.NewGRPCController(c)
	g := grpc.NewServer()
	pb.RegisterNotificationServiceServer(g, gc)
	err = g.Serve(lis)
	if err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}
