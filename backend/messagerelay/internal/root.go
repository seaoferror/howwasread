package internal

import (
	"backend/common/producer"
	"backend/messagerelay/internal/consumer"
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

	p := producer.NewProducer("producer_message_notification")

	s := service.NewService(r, p)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"manage_message.prepared"})
	if err != nil {
		panic(err)
	}
}
