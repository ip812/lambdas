package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	key = "LOGGER"
)

type Logger struct {
	Log *zerolog.Logger
}

func NewLogger() *Logger {
	var log zerolog.Logger
	if Profile == "local" {
		log = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	} else {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = time.Now().UTC
	return &Logger{
		Log: &log,
	}
}

func inject(ctx context.Context, log *zerolog.Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func logger(ctx context.Context) *zerolog.Logger {
	c, ok := ctx.Value(key).(*zerolog.Logger)
	if !ok {
		log.Fatal("couldn't get logger from context")
	}
	return c
}
