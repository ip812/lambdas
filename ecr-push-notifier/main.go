package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ip812/ecr-push-notifier/config"
	"github.com/ip812/ecr-push-notifier/git"
	"github.com/ip812/ecr-push-notifier/logger"
	"github.com/ip812/ecr-push-notifier/notifier"
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
	ErrInvalidEvent  = fmt.Errorf("invalid event")
	ErrInvalidDetail = fmt.Errorf("invalid detail")
)

func pickGitTarget(detail ECRDetail) (*git.Target, error) {
	switch detail.RepositoryName {
	case "ip812/hello", "ip812/pg-query-exec", "ip812/ecr-push-notifier":
		return &git.Target{
			Type:          git.Lambda,
			RepositroyURL: "https://github.com/ip812/infra.git",
			FilePath:      "prod/lambdas.tf",
			Branch:        "main",
			ImageName:     detail.RepositoryName,
			ImageTag:      detail.ImageTag,
		}, nil
	case "ip812/go-template":
		return &git.Target{
			Type:          git.Service,
			RepositroyURL: "https://github.com/ip812/apps.git",
			FilePath:      "manifests/prod/go-template/deployment.yaml",
			Branch:        "main",
			ImageName:     detail.RepositoryName,
			ImageTag:      detail.ImageTag,
		}, nil
	default:
		return nil, fmt.Errorf("unknown repository: %s", detail.RepositoryName)
	}
}

func Handler(ctx context.Context, event events.EventBridgeEvent) (interface{}, error) {
	log := logger.Get(ctx)
	cfg := config.Get(ctx)

	var detail ECRDetail
	if err := json.Unmarshal(event.Detail, &detail); err != nil {
		log.Error("Failed to unmarshal event detail: %v", err)
		return nil, ErrInvalidEvent
	}

	if detail.ActionType != "PUSH" || detail.Result != "SUCCESS" {
		log.Error("Not a push event or not successful: %v", detail)
		return nil, ErrInvalidDetail
	}

	trg, err := pickGitTarget(detail)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Info("Target: %v", trg)

	notifier := notifier.NewSlack(cfg.Slack.BotToken, log)
	var slackTargetChannel string
	if trg.Type == git.Lambda {
		slackTargetChannel = cfg.Slack.AWSChannelID
	} else {
		slackTargetChannel = cfg.Slack.K8sChannelID
	}

	git, err := git.New(
		log,
		cfg.Git.Username,
		cfg.Git.AccessToken,
		*trg,
	)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer git.Close()
	log.Info("Git client initialized")

	if err := git.ReplaceImageTag(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Info("Image tag replaced")

	if err := git.Push(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Info("Push successful")

	err = notifier.SendSuccessNotification(
		slackTargetChannel,
		detail.RepositoryName,
		detail.ImageTag,
	)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return event, nil
}

func main() {
	ctx := context.Background()
	cfg := config.New()
	log := logger.New(cfg)
	ctx = log.Inject(ctx)
	ctx = config.Inject(ctx, *cfg)

	lambda.StartWithOptions(Handler, lambda.WithContext(ctx))
}
