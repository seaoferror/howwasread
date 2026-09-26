package internal

import (
	"backend/common"
	"backend/messagerelay/internal/client"
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

	p := common.NewProducer("producer_message_notification")

	s := service.NewService(r, p, client.NewRelayClient())

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"prepared-message"})
	if err != nil {
		panic(err)
	}
}
