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
}

func Handler(ctx context.Context, event QueryEvent) (interface{}, error) {
	log := logger(ctx)

	log.Info().Msgf("event: %+v", event)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		event.DatabaseName,
	)
	log.Info().Msgf("dbURL: %s", dbURL)

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
			return "", fmt.Errorf("no rows found")
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
