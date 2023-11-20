package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/petrenko-alex/api-rate-limiter/internal/app"
	"github.com/petrenko-alex/api-rate-limiter/internal/config"
	"github.com/petrenko-alex/api-rate-limiter/internal/server"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	configFile, fileErr := os.Open(configFilePath)
	if fileErr != nil {
		log.Println("Error opening config file: ", fileErr)

		return 1
	}

	cfg, configErr := config.New(ctx, configFile)
	if configErr != nil {
		log.Println("Error parsing config file: ", configErr)

		return 1
	}

	ctx = cfg.WithContext(ctx)

	logg := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: cfg.Logger.Level},
	))

	application, appInitErr := app.New(ctx, cfg, logg)
	if appInitErr != nil {
		logg.Error("Failed to init app: " + appInitErr.Error())

		return 1
	}

	srv := server.NewServer(
		server.Options{
			Host:           cfg.Server.Host,
			Port:           cfg.Server.Port,
			ConnectTimeout: cfg.Server.ConnectTimeout,
		},
		application,
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

	return 0
}
