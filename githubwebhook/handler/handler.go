package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"github.com/int128/codebuild-runner/githubwebhook/usecases"
)

func Handle(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	code, err := handle(ctx, r.Headers, r.Body)
	if err != nil {
		log.Printf("error: %+v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: code}, nil
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: code}, nil
}

func handle(ctx context.Context, header map[string]string, body string) (int, error) {
	event := header["x-github-event"]
	if event == "push" {
		var e github.PushEvent
		if err := json.Unmarshal([]byte(body), &e); err != nil {
			return 500, fmt.Errorf("could not decode json: %w", err)
		}
		if err := usecases.PushEvent(ctx, e); err != nil {
			return 500, fmt.Errorf("push: %w", err)
		}
		return 200, nil
	}
	if event == "pull_request" {
		var e github.PullRequestEvent
		if err := json.Unmarshal([]byte(body), &e); err != nil {
			return 500, fmt.Errorf("could not decode json: %w", err)
		}
		if err := usecases.PullRequestEvent(ctx, e); err != nil {
			return 500, fmt.Errorf("pull_request: %w", err)
		}
		return 200, nil
	}
	return 404, fmt.Errorf("unknown event `%s` in header %+v", event, header)
}
