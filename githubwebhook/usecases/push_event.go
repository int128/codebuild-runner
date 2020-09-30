package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"
	"github.com/google/go-github/v32/github"
	"github.com/int128/codebuild-runner/githubwebhook/builder"
)

func PushEvent(ctx context.Context, e github.PushEvent) error {
	if e.HeadCommit == nil {
		return fmt.Errorf("HeadCommit is nil")
	}

	buildInput := calculateBuildInputForPushEvent(e)
	buildOutput, err := builder.Start(ctx, buildInput)
	if err != nil {
		return fmt.Errorf("could not start a build: %w", err)
	}
	log.Printf("build started %+v", buildOutput.Build)

	return nil
}

func calculateBuildInputForPushEvent(e github.PushEvent) *codebuild.StartBuildInput {
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
	log.Printf("commit=%v, changed=%+v", e.HeadCommit.ID, changed)

	// https://docs.aws.amazon.com/codebuild/latest/userguide/sample-source-version.html
	// A tag (for example, refs/tags/mytagv1.0^{full-commit-SHA}).
	// A branch (for example, refs/heads/mydevbranch^{full-commit-SHA}).
	sourceVersion := e.HeadCommit.GetID()
	ref := e.GetRef()
	if ref != "" {
		sourceVersion = fmt.Sprintf("%s^{%s}", ref, e.HeadCommit.GetID())
	}

	//TODO: compute jobs from .codebuild/workflows/*.yaml

	return &codebuild.StartBuildInput{
		ProjectName:   aws.String("codebuild-runner"),
		SourceVersion: aws.String(sourceVersion),
		EnvironmentVariablesOverride: []*types.EnvironmentVariable{
			{
				Name:  aws.String("GITHUB_WEBHOOK_HEADCOMMIT_ID"),
				Value: e.HeadCommit.ID,
				Type:  types.EnvironmentVariableTypePlaintext,
			},
			{
				Name:  aws.String("GITHUB_WEBHOOK_REF"),
				Value: e.Ref,
				Type:  types.EnvironmentVariableTypePlaintext,
			},
		},
	}
}
