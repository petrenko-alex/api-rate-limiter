package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	proto "github.com/petrenko-alex/api-rate-limiter/api"
	"github.com/petrenko-alex/api-rate-limiter/internal/app"
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	proto.UnimplementedRateLimiterServer

	app    app.IApp
	logger slog.Logger
}

func NewService(app app.IApp, logger slog.Logger) *Service {
	return &Service{
		logger: logger,
		app:    app,
	}
}

func (s Service) WhiteListAdd(_ context.Context, req *proto.WhiteListAddRequest) (*proto.WhiteListAddResponse, error) {
	err := s.app.WhiteListAdd(req.IpNet)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed adding to white list: %s", err))

		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &proto.WhiteListAddResponse{}, nil
}

func (s Service) WhiteListDelete(_ context.Context, req *proto.WhiteListDeleteRequest) (*proto.WhiteListDeleteResponse, error) { //nolint:lll
	err := s.app.WhiteListDelete(req.IpNet)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed deleting from white list: %s", err))

		code := codes.Unknown
		if errors.Is(err, ipnet.ErrRuleNotFound) {
			code = codes.NotFound
		}

		return nil, status.Errorf(code, err.Error())
	}

	return &proto.WhiteListDeleteResponse{}, nil
}

func (s Service) BlackListAdd(_ context.Context, req *proto.BlackListAddRequest) (*proto.BlackListAddResponse, error) {
	err := s.app.BlackListAdd(req.IpNet)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed adding to black list: %s", err))

		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &proto.BlackListAddResponse{}, nil
}

func (s Service) BlackListDelete(_ context.Context, req *proto.BlackListDeleteRequest) (*proto.BlackListDeleteResponse, error) { //nolint:lll
	err := s.app.BlackListDelete(req.IpNet)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed deleting from black list: %s", err))

		code := codes.Unknown
		if errors.Is(err, ipnet.ErrRuleNotFound) {
			code = codes.NotFound
		}

		return nil, status.Errorf(code, err.Error())
	}

	return &proto.BlackListDeleteResponse{}, nil
}

func (s Service) BucketReset(_ context.Context, req *proto.BucketResetRequest) (*proto.BucketResetResponse, error) {
	err := s.app.LimitReset(req.Ip, req.Login)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed resetting limits: %s", err))

		code := codes.Unknown

		return nil, status.Errorf(code, err.Error())
	}

	return &proto.BucketResetResponse{}, nil
}

func (s Service) LimitCheck(_ context.Context, req *proto.LimitCheckRequest) (*proto.LimitCheckResponse, error) {
	satisfies, err := s.app.LimitCheck(req.Ip, req.Login, req.Password)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed checking limit: %s", err))

		code := codes.Unknown
		if errors.Is(err, limiter.ErrIncorrectIdentity) {
			code = codes.InvalidArgument
		}

		return nil, status.Errorf(code, err.Error())
	}

	return &proto.LimitCheckResponse{Allowed: satisfies}, nil
}
