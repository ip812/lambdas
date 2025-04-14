package config

import (
	"context"
	"log"
	"os"
)

type Environment string

const (
	Local Environment = "local"
	Prod  Environment = "prod"
)

func (e Environment) IsValid() bool {
	switch e {
	case Local, Prod:
		return true
	default:
		return false
	}
}

var (
	key = "CONFIG"
)

type Config struct {
	App struct {
		Env Environment
	}

	Database struct {
		Endpoint string
		Username string
		Password string
		SSLMode  string
	}
}

func Inject(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, key, cfg)
}

func Get(ctx context.Context) Config {
	c, ok := ctx.Value(key).(Config)
	if !ok {
		log.Fatal("couldn't get config from context")
	}
	return c
}

func New() *Config {
	cfg := &Config{}

	cfg.App.Env = Environment(os.Getenv("APP_ENV"))
	if !cfg.App.Env.IsValid() {
		cfg.App.Env = Prod
	}

	cfg.Database.Endpoint = os.Getenv("DB_ENDPOINT")
	cfg.Database.Username = os.Getenv("DB_USERNAME")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.SSLMode = os.Getenv("DB_SSL_MODE")

	return cfg
}
