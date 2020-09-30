package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"github.com/int128/codebuild-runner/webhook"
)

func HandleGitHubWebhook(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	code, err := handleGitHubWebhook(ctx, r.Headers, r.Body)
	if err != nil {
		log.Printf("error: %+v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: code,
		}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
	}, nil
}

func handleGitHubWebhook(ctx context.Context, header map[string]string, body string) (int, error) {
	eventKind := header["x-github-event"]
	if eventKind == "push" {
		var e github.PushEvent
		if err := json.Unmarshal([]byte(body), &e); err != nil {
			return 500, fmt.Errorf("could not decode json: %w", err)
		}
		if err := webhook.PushEvent(ctx, e); err != nil {
			return 500, fmt.Errorf("push: %w", err)
		}
		return 200, nil
	}
	if eventKind == "pull_request" {
		var e github.PullRequestEvent
		if err := json.Unmarshal([]byte(body), &e); err != nil {
			return 500, fmt.Errorf("could not decode json: %w", err)
		}
		if err := webhook.PullRequestEvent(ctx, e); err != nil {
			return 500, fmt.Errorf("pull_request: %w", err)
		}
		return 200, nil
	}
	return 404, fmt.Errorf("unknown event `%s` in header %+v", eventKind, header)
}
