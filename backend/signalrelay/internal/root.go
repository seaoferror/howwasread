package internal

import (
	"backend/common/producer"
	"backend/signalrelay/internal/consumer"
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

	kp := producer.NewProducer("producer_signal_relay")

	s := service.NewService(r, kp)

	ks := consumer.NewKafkaConsumer(s)

	ks.GetMessage([]string{"conversation.signal"})
}
