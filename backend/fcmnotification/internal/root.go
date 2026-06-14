package internal

import (
	"backend/common/producer"
	"backend/fcmnotification/internal/consumer"
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

	p := producer.NewProducer("producer_fcm_notification")

	r := repository.NewRepository()

	s := service.NewService(r, p)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"fcm_notification"})
	if err != nil {
		panic(err)
	}
}
