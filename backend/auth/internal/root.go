package internal

import (
	"backend/auth/internal/controller"
	"backend/auth/internal/repository"
	"backend/auth/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func NewServer() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	slog.SetDefault(logger)

	log.Print("success to set logger")

	r := repository.NewRepository()

	s := service.NewService(r)

	mux := http.NewServeMux()

	controller.NewController(s, mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
