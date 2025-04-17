package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ip812/ecr-push-notifier/config"
	"github.com/ip812/ecr-push-notifier/logger"
)

type ECRDetail struct {
	ActionType        string `json:"action-type"`
	RepositoryName    string `json:"repository-name"`
	ImageTag          string `json:"image-tag"`
	ImageDigest       string `json:"image-digest"`
	ArtifactMediaType string `json:"artifact-media-type"`
	ManifestMediaType string `json:"manifest-media-type"`
	Result            string `json:"result"`
}

var (
	ErrInvalidEvent = fmt.Errorf("invalid event")
)

func Handler(ctx context.Context, event events.EventBridgeEvent) (interface{}, error) {
	log := logger.Get(ctx)
	log.Info("Hello, received event is: %v", event)
	var detail ECRDetail
	if err := json.Unmarshal(event.Detail, &detail); err != nil {
		return nil, ErrInvalidEvent
	}

	log.Info("📦 New ECR image pushed!")
	log.Info("Repository: %s", detail.RepositoryName)
	log.Info("Tag: %s", detail.ImageTag)
	log.Info("Digest: %s", detail.ImageDigest)

	return event, nil
}

func main() {
	ctx := context.Background()
	cfg := config.New()
	log := logger.New(cfg)
	ctx = log.Inject(ctx)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
