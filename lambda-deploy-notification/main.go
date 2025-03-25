package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event interface{}) (interface{}, error) {
	log := GetLogger(ctx)
	log.Info("Hello, received event is: %v", event)
	return event, nil
}

func main() {
	ctx := context.Background()
	ctx = InjectLogger(ctx, NewLogger())

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
