package internal

import (
	"backend/messagepreprocess/internal/kafka/consumer"
	"backend/messagepreprocess/internal/kafka/producer"
	"backend/messagepreprocess/internal/repository"
	"backend/messagepreprocess/internal/service"
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
