package internal

import (
	"backend/common"
	"backend/signalrelay/internal/client"
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

	kp := common.NewProducer("producer_signal_relay")

	s := service.NewService(r, kp, client.NewRelayClient())

	ks := consumer.NewKafkaConsumer(s)

	err := ks.GetMessage([]string{"conversation-signal"})
	if err != nil {
		panic(err)
	}
}
