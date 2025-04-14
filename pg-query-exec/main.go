package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ip812/pg-query-exec/config"
	"github.com/ip812/pg-query-exec/logger"
	"github.com/lib/pq"
)

type QueryEvent struct {
	DatabaseName string `json:"database_name"`
	Query        string `json:"query"`
}

func Handler(ctx context.Context, event QueryEvent) (interface{}, error) {
	log := logger.Get(ctx)
	cfg := config.Get(ctx)

	log.Info("event: %+v", event)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Endpoint,
		event.DatabaseName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Info("failed to connect to database: %s", err.Error())
		return nil, err
	}
	defer db.Close()

	row := db.QueryRowContext(ctx, event.Query)
	var result string
	err = row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("no rows found")
			return "", nil
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			log.Info("database already exists, ignoring error")
			return "", nil
		}
		log.Error("failed to execute query: %s", err.Error())
		return "", err
	}

	log.Info("result: %s", result)
	return result, nil
}

func main() {
	ctx := context.Background()
	cfg := config.New()
	log := logger.New(cfg)
	ctx = config.Inject(ctx, *cfg)
	ctx = log.Inject(ctx)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
