package server

import (
	"context"
	"log/slog"
	"net"

	proto "github.com/petrenko-alex/api-rate-limiter/api"
	"github.com/petrenko-alex/api-rate-limiter/internal/app"
	"github.com/petrenko-alex/api-rate-limiter/internal/server/log"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server

	logger *slog.Logger

	options Options
}

func NewServer(options Options, app app.IApp, logger *slog.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(options.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			log.NewInterceptor(logger).GetInterceptor(),
		),
	)
	proto.RegisterRateLimiterServer(grpcServer, NewService(app, *logger))

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
