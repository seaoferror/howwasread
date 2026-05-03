package internal

import (
	"backend/fcmnotification/internal/kafka/consumer"
	"backend/fcmnotification/internal/repository"
	"backend/fcmnotification/internal/service"
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

	err := c.GetMessage([]string{"fcm_notification"})
	if err != nil {
		panic(err)
	}
}
