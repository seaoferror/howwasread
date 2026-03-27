package internal

import (
	"backend/online/internal/controller"
	"backend/online/internal/grpccontroller"
	"backend/online/internal/kafka/consumer"
	"backend/online/internal/kafka/producer"
	"backend/online/internal/repository"
	"backend/online/internal/service"
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

	kp := producer.NewKafkaProducer()

	r := repository.NewRepository()

	s := service.NewService(r, kp)

	ks := consumer.NewKafkaConsumer(s)

	go func() {
		err := ks.GetMessage([]string{"auth.new_member_id"})
		if err != nil {
			slog.Error("fail to get payload from kafka",
				"err", err)
			return
		}
	}()

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
	pb.RegisterSignalServiceServer(g, gc)
	err = g.Serve(lis)
	if err != nil {
		log.Fatalf("fail to serve: %v", err)
	}

}
