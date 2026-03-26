package logger

import (
	"backend/auth/internal/kafka/producer"
	"log"
	"log/slog"
	"os"
)

func SetLogger(kafkaProducer *producer.KafkaProducer) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	log.Print("success to set logger")
}
