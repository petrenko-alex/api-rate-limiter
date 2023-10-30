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

func (s Service) WhiteListAdd(ctx context.Context, request *proto.WhiteListAddRequest) (*proto.WhiteListAddResponse, error) {
	s.logger.Info("WhiteListAdd executing...")

	return &proto.WhiteListAddResponse{}, nil
}

func (s Service) WhiteListDelete(ctx context.Context, request *proto.WhiteListDeleteRequest) (*proto.WhiteListDeleteResponse, error) {
	s.logger.Info("WhiteListDelete executing...")

	return &proto.WhiteListDeleteResponse{}, nil
}

func (s Service) BlackListAdd(ctx context.Context, request *proto.BlackListAddRequest) (*proto.BlackListAddResponse, error) {
	s.logger.Info("BlackListAdd executing...")

	return &proto.BlackListAddResponse{}, nil
}

func (s Service) BlackListDelete(ctx context.Context, request *proto.BlackListDeleteRequest) (*proto.BlackListDeleteResponse, error) {
	s.logger.Info("BlackListDelete executing...")

	return &proto.BlackListDeleteResponse{}, nil
}

func (s Service) BucketReset(ctx context.Context, request *proto.BucketResetRequest) (*proto.BucketResetResponse, error) {
	s.logger.Info("BucketReset executing...")

	return &proto.BucketResetResponse{}, nil
}

func (s Service) LimitCheck(ctx context.Context, request *proto.LimitCheckRequest) (*proto.LimitCheckResponse, error) {
	s.logger.Info("LimitCheck executing...")

	return &proto.LimitCheckResponse{}, nil
}
