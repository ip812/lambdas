package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

type QueryEvent struct {
	DatabaseName string `json:"database_name"`
	Query        string `json:"query"`
	SSLMode      string `json:"ssl_mode"`
}

func Handler(ctx context.Context, event QueryEvent) (interface{}, error) {
	log := logger(ctx)

	log.Info().Msgf("event: %+v", event)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		event.DatabaseName,
		event.SSLMode,
	)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Info().Msgf("failed to connect to database: %s", err.Error())
		return nil, err
	}
	defer db.Close()

	row := db.QueryRowContext(ctx, event.Query)
	var result string
	err = row.Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	log.Info().Msgf("result: %s", result)
	return result, nil
}

func main() {
	ctx := context.Background()
	ctx = inject(ctx, NewLogger().Log)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
