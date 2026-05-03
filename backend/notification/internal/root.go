package internal

import (
	"backend/common/producer"
	"backend/notification/internal/consumer"
	"backend/notification/internal/controller"
	"backend/notification/internal/repository"
	"backend/notification/internal/service"
	"log/slog"
	"net/http"
	"os"
)

func NewServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	r := repository.NewRepository()

	p := producer.NewProducer("producer.preprocess_notification")

	s := service.NewService(r, p)

	mux := http.NewServeMux()

	controller.SetController(s, mux)

	go func() {
		err := http.ListenAndServe(":8078", mux)
		if err != nil {
			slog.Error("fail to listen and serve http",
				"err", err)
			panic(err)
		}
	}()

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"preprocess_notification"})
	if err != nil {
		panic(err)
	}
}
