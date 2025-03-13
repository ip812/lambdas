package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event interface{}) (interface{}, error) {
	log := logger(ctx)
	log.Info().Msgf("Hello, received event is: %v", event)
	return event, nil
}

func main() {
	ctx := context.Background()
	ctx = inject(ctx, NewLogger().Log)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
