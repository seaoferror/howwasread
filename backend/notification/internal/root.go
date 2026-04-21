package internal

import (
	"backend/notification/internal/controller"
	"backend/notification/internal/kafka/consumer"
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

	s := service.NewService(r)

	mux := http.NewServeMux()

	controller.SetController(s, mux)

	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			slog.Error("fail to listen and serve http",
				"err", err)
			panic(err)
		}
	}()

	c := consumer.NewConsumer(s)

	err := c.GetMessage([]string{"notify_message"})
	if err != nil {
		panic(err)
	}
}
