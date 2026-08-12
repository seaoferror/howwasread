package internal

import (
	"backend/turn/internal/manager"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServer() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	t, err := manager.RunTurnManager(20 * time.Second)
	if err != nil {
		panic(err)
	}
	defer t.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
