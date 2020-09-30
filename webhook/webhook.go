package webhook

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v32/github"
)

func PushEvent(ctx context.Context, e github.PushEvent) error {
	if e.HeadCommit == nil {
		return fmt.Errorf("HeadCommit is nil")
	}

	changed := make(map[string]int)
	for _, name := range e.HeadCommit.Added {
		changed[name]++
	}
	for _, name := range e.HeadCommit.Removed {
		changed[name]++
	}
	for _, name := range e.HeadCommit.Modified {
		changed[name]++
	}
	log.Printf("changed=%+v", changed)
	return nil
}

func PullRequestEvent(ctx context.Context, e github.PullRequestEvent) error {
	log.Printf("payload=%+v", e)
	return nil
}
