package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ip812/hello/config"
	"github.com/ip812/hello/logger"
)

func Handler(ctx context.Context, event interface{}) (interface{}, error) {
	log := logger.Get(ctx)
	log.Info("Hello, received event is: %v", event)
	return event, nil
}

func main() {
	ctx := context.Background()
	cfg := config.New()
	log := logger.New(cfg)
	ctx = log.Inject(ctx)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
