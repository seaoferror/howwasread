package internal

import (
	"backend/submessaging/internal/kafka/consumer"
	"backend/submessaging/internal/kafka/producer"
	"backend/submessaging/internal/repository"
	"backend/submessaging/internal/service"
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

	ks.GetMessage([]string{"chat.messaging"})
}
