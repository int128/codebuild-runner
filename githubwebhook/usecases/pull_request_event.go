package usecases

import (
	"context"
	"log"

	"github.com/google/go-github/v32/github"
)

func PullRequestEvent(_ context.Context, e github.PullRequestEvent) error {
	log.Printf("payload=%+v", e)
	return nil
}
