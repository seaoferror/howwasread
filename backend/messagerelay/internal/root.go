package internal

import (
	"backend/messagerelay/internal/kafka/consumer"
	"backend/messagerelay/internal/kafka/producer"
	"backend/messagerelay/internal/repository"
	"backend/messagerelay/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	p := producer.NewProducer()

	s := service.NewService(r, p)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"manage_message.prepared"})
	if err != nil {
		panic(err)
	}
}
