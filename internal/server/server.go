package server

import (
	"context"
	"log/slog"
	"net"

	proto "api-rate-limiter/api"
	"api-rate-limiter/internal/server/log"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server

	logger *slog.Logger

	options Options
}

func NewServer(options Options, logger *slog.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(options.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			log.NewInterceptor(logger).GetInterceptor(),
		),
	)
	proto.RegisterRateLimiterServer(grpcServer, NewService(*logger))

	return &Server{
		server:  grpcServer,
		logger:  logger,
		options: options,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var lc net.ListenConfig

	listener, err := lc.Listen(
		ctx,
		"tcp",
		net.JoinHostPort(s.options.Host, s.options.Port),
	)
	if err != nil {
		return err
	}

	err = s.server.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}