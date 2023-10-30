package server

import (
	"context"
	"log/slog"

	proto "api-rate-limiter/api"
)

type Service struct {
	proto.UnimplementedRateLimiterServer

	logger slog.Logger
}

func NewService(logger slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s Service) WhiteListAdd(_ context.Context, _ *proto.WhiteListAddRequest) (*proto.WhiteListAddResponse, error) {
	s.logger.Info("WhiteListAdd executing...")

	return &proto.WhiteListAddResponse{}, nil
}

func (s Service) WhiteListDelete(_ context.Context, _ *proto.WhiteListDeleteRequest) (*proto.WhiteListDeleteResponse, error) { //nolint:lll
	s.logger.Info("WhiteListDelete executing...")

	return &proto.WhiteListDeleteResponse{}, nil
}

func (s Service) BlackListAdd(_ context.Context, _ *proto.BlackListAddRequest) (*proto.BlackListAddResponse, error) {
	s.logger.Info("BlackListAdd executing...")

	return &proto.BlackListAddResponse{}, nil
}

func (s Service) BlackListDelete(_ context.Context, _ *proto.BlackListDeleteRequest) (*proto.BlackListDeleteResponse, error) { //nolint:lll
	s.logger.Info("BlackListDelete executing...")

	return &proto.BlackListDeleteResponse{}, nil
}

func (s Service) BucketReset(_ context.Context, _ *proto.BucketResetRequest) (*proto.BucketResetResponse, error) {
	s.logger.Info("BucketReset executing...")

	return &proto.BucketResetResponse{}, nil
}

func (s Service) LimitCheck(_ context.Context, _ *proto.LimitCheckRequest) (*proto.LimitCheckResponse, error) {
	s.logger.Info("LimitCheck executing...")

	return &proto.LimitCheckResponse{}, nil
}
