package internal

import (
	"backend/messageprocess/internal/kafka/consumer"
	"backend/messageprocess/internal/kafka/producer"
	"backend/messageprocess/internal/repository"
	"backend/messageprocess/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	kp := producer.NewProducer()

	s := service.NewService(r, kp)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"chat.message"})
	if err != nil {
		panic(err)
	}
}
