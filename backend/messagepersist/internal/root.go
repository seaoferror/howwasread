package internal

import (
	"backend/messagepersist/internal/kafka/consumer"
	"backend/messagepersist/internal/repository"
	"backend/messagepersist/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	s := service.NewService(r)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"manage_message.prepared"})
	if err != nil {
		panic(err)
	}
}
