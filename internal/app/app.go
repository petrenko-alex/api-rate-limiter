package app

import (
	"context"
	"log/slog"

	"github.com/petrenko-alex/api-rate-limiter/internal/config"
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
)

type App struct {
	ruleService    ipnet.IRuleService
	limiterService limiter.ILimitService

	ctx    context.Context
	logger *slog.Logger
	config *config.Config
}

func New(ctx context.Context, config *config.Config, logger *slog.Logger) (*App, error) {
	ruleStorage := ipnet.NewRuleStorage(config.DB.DSN)
	err := ruleStorage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	limitStorage := limiter.NewLimitStorage(config.DB.DSN)
	err = limitStorage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ruleService := ipnet.NewRuleService(ruleStorage)

	limiterService := limiter.NewLoginFormLimiter(
		ruleService,
		limitStorage,
		limiter.NewRefillRate(config.App.RefillRate.Count, config.App.RefillRate.Time),
	)

	return &App{
		ruleService:    ruleService,
		limiterService: limiterService,

		logger: logger,
		ctx:    ctx,
		config: config,
	}, nil
}

func (a *App) LimitCheck(ip, login, password string) (bool, error) {
	return a.limiterService.SatisfyLimit(limiter.UserIdentityDto{
		limiter.IPLimit.String():       ip,
		limiter.LoginLimit.String():    login,
		limiter.PasswordLimit.String(): password,
	})
}

func (a *App) WhiteListAdd(ip string) error {
	return a.ruleService.WhiteListAdd(ip)
}

func (a *App) WhiteListDelete(ip string) error {
	return a.ruleService.WhiteListDelete(ip)
}

func (a *App) BlackListAdd(ip string) error {
	return a.ruleService.BlackListAdd(ip)
}

func (a *App) BlackListDelete(ip string) error {
	return a.ruleService.BlackListDelete(ip)
}
