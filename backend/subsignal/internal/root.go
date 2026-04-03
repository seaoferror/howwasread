package internal

import (
	"backend/subsignal/internal/kafka/consumer"
	"backend/subsignal/internal/kafka/producer"
	"backend/subsignal/internal/repository"
	"backend/subsignal/internal/service"
	"log/slog"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	kp := producer.NewKafkaProducer()

	s := service.NewService(r, kp)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})
}
