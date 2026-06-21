package internal

import (
	"backend/common/producer"
	"backend/messagepreprocess/internal/consumer"
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

	p := producer.NewProducer("message_preprocess")

	s := service.NewService(r, p)

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"chat-message"})
	if err != nil {
		panic(err)
	}
}
