package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-rate-limiter/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	logg := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	srv := server.NewServer(
		server.Options{
			Host:           "localhost",
			Port:           "4242",
			ConnectTimeout: time.Second * 1,
		},
		logg,
	)

	logg.Info("Starting GRPC server...")

	go func() {
		<-ctx.Done()

		logg.Info("Stopping GRPC server...")
		if err := srv.Stop(ctx); err != nil {
			logg.Error("Failed to stop GRPC server: " + err.Error())
		}

		logg.Info("Server stopped.")
	}()

	if err := srv.Start(ctx); err != nil {
		logg.Error("Failed to start GRPC server: " + err.Error())
		cancel()
	}
}
