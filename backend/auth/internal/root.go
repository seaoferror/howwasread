package internal

import (
	"backend/auth/internal/network"
	"backend/auth/internal/repository"
	"backend/auth/internal/service"
	"log"
	"log/slog"
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

	n := network.NewNetwork(s)

	n.Start()
}
