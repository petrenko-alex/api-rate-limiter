package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/petrenko-alex/api-rate-limiter/internal/config"
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
)

type App struct {
	ruleService    ipnet.IRuleService
	limiterService limiter.ILimitService
	limiterGB      limiter.ITokenBucketGB

	ctx    context.Context
	logger *slog.Logger
	config *config.Config
}

func New(ctx context.Context, config *config.Config, logger *slog.Logger) (*App, error) {
	// Init Rule Service
	ruleStorage := ipnet.NewRuleStorage(config.DB.DSN)
	err := ruleStorage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ruleService := ipnet.NewRuleService(ruleStorage)

	// Init Limiter Service
	limitStorage := limiter.NewLimitStorage(config.DB.DSN)
	err = limitStorage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	bucketLimiter := limiter.NewCompositeBucketLimiter(
		limitStorage,
		limiter.NewRefillRate(config.App.RefillRate.Count, config.App.RefillRate.Time),
	)
	limiterService := limiter.NewLoginFormLimiter(ruleService, bucketLimiter)

	// Init Limiter Garbage Collector
	limiterGB := limiter.NewTokenBucketGB(bucketLimiter, config.App.GarbageCollector.TTL)

	return &App{
		ruleService:    ruleService,
		limiterService: limiterService,
		limiterGB:      limiterGB,

		logger: logger,
		ctx:    ctx,
		config: config,
	}, nil
}

func (a *App) RunBackground() {
	a.RunGB()
}

func (a *App) LimitCheck(ip, login, password string) (bool, error) {
	return a.limiterService.SatisfyLimit(limiter.UserIdentityDto{
		limiter.IPLimit.String():       ip,
		limiter.LoginLimit.String():    login,
		limiter.PasswordLimit.String(): password,
	})
}

func (a *App) LimitReset(ip, login string) error {
	return a.limiterService.ResetLimit(limiter.UserIdentityDto{
		limiter.IPLimit.String():    ip,
		limiter.LoginLimit.String(): login,
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

func (a *App) RunGB() {
	if !a.config.App.GarbageCollector.Enabled {
		return
	}

	a.logger.Info("Starting Garbage Collector...")
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				a.logger.Info("GB finished.")

				return
			case <-time.After(a.config.App.GarbageCollector.Interval):
				a.logger.Info("GB sweeping..")

				err := a.limiterGB.Sweep()
				if err != nil {
					a.logger.Error("GB error", "error", err)
				}
			}
		}
	}()
}
