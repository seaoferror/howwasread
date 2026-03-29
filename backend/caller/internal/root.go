package internal

import (
	"backend/caller/internal/kafka/consumer"
	"backend/caller/internal/repository"
	"backend/caller/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	s := service.NewService(r)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal", "chat.messaging"})
}
