package config

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	_ "github.com/lib/pq" // driver import
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"gopkg.in/yaml.v3"
)

type configCtxKey struct{}

var ErrLimitsNotFound = errors.New("no rate limits found")

type Config struct {
	Server struct {
		Host, Port     string
		ConnectTimeout time.Duration `yaml:"connectTimeout"`
	}

	DB struct {
		DSN           string
		MigrationsDir string `yaml:"migrationsDir"`
	}

	Limits struct {
		Login    *limiter.Limit
		Password *limiter.Limit
		IP       *limiter.Limit
	}

	Logger struct {
		Level slog.Level
	}

	App struct {
		RefillRate struct {
			Count int
			Time  time.Duration
		} `yaml:"refillRate"`
	}
}

func (c Config) WithContext(parentCtx context.Context) context.Context {
	return context.WithValue(parentCtx, configCtxKey{}, c)
}

func New(ctx context.Context, configFile io.Reader) (*Config, error) {
	config, err := ForMigrator(ctx, configFile)
	if err != nil {
		return nil, err
	}

	limitStorage := limiter.NewLimitStorage(config.DB.DSN)
	if err = limitStorage.Connect(ctx); err != nil {
		return nil, err
	}

	limits, err := limitStorage.GetLimits()
	if err != nil || len(*limits) == 0 {
		return nil, ErrLimitsNotFound
	}

	for _, limit := range *limits {
		limit := limit
		switch limit.LimitType {
		case limiter.LoginLimit:
			config.Limits.Login = &limit
		case limiter.PasswordLimit:
			config.Limits.Password = &limit
		case limiter.IPLimit:
			config.Limits.IP = &limit
		}
	}

	if config.Limits.Login == nil || config.Limits.Password == nil || config.Limits.IP == nil {
		return nil, ErrLimitsNotFound
	}

	return config, nil
}

func ForMigrator(_ context.Context, configFile io.Reader) (*Config, error) {
	config := &Config{}

	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
