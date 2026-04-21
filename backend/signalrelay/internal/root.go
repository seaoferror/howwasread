package internal

import (
	"backend/signalrelay/internal/kafka/consumer"
	"backend/signalrelay/internal/kafka/producer"
	"backend/signalrelay/internal/repository"
	"backend/signalrelay/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	kp := producer.NewKafkaProducer()

	s := service.NewService(r, kp)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})
}
