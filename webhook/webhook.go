package webhook

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/google/go-github/v32/github"
	"github.com/int128/codebuild-runner/builder"
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

	//TODO: compute jobs from .codebuild/workflows/*.yaml

	buildOutput, err := builder.Start(ctx, &codebuild.StartBuildInput{
		ProjectName: aws.String("codebuild-runner"),
	})
	if err != nil {
		return fmt.Errorf("could not start a build: %w", err)
	}
	log.Printf("build started %+v", buildOutput.Build)

	return nil
}

func PullRequestEvent(_ context.Context, e github.PullRequestEvent) error {
	log.Printf("payload=%+v", e)
	return nil
}
