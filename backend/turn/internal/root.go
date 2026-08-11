package internal

import (
	"backend/turn/internal/manager"
	"log"
	"log/slog"
	"net"
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

	pc, err := net.ListenPacket("udp4", "0.0.0.0:3478")
	if err != nil {
		log.Panicf("fail to create udp listener: %v", err)
	}
	defer pc.Close()
	t, err := manager.RunTurnManager(pc, 20*time.Second)
	if err != nil {
		panic(err)
	}
	defer t.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
