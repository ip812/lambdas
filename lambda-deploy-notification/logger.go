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
	log *zerolog.Logger
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
		log: &log,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log.Debug().Msgf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log.Info().Msgf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log.Warn().Msgf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log.Error().Msgf(format, v...)
}

func InjectLogger(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func GetLogger(ctx context.Context) *Logger {
	c, ok := ctx.Value(key).(*Logger)
	if !ok {
		log.Fatal("couldn't get logger from context")
	}
	return c
}
