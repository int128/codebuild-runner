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
	"github.com/int128/codebuild-runner/githubwebhook/workflow"
)

func PushEvent(ctx context.Context, e github.PushEvent) error {
	if e.HeadCommit == nil {
		return fmt.Errorf("HeadCommit is nil")
	}

	workflows, err := workflow.Find(ctx, e.GetRepo().GetOwner().GetName(), e.GetRepo().GetName(), e.GetHeadCommit().GetID())
	if err != nil {
		return fmt.Errorf("could not find workflows: %w", err)
	}
	buildInput, err := calculateBuildInputForPushEvent(e, workflows)
	if err != nil {
		return fmt.Errorf("could not calcluate build from workflows: %w", err)
	}
	buildOutput, err := builder.Start(ctx, buildInput)
	if err != nil {
		return fmt.Errorf("could not start a build: %w", err)
	}
	log.Printf("build started %+v", buildOutput.Build)

	return nil
}

func calculateBuildInputForPushEvent(e github.PushEvent, workflows []*workflow.Workflow) (*codebuild.StartBuildInput, error) {
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
	var changedPaths []string
	for name := range changed {
		changedPaths = append(changedPaths, name)
	}

	// https://docs.aws.amazon.com/codebuild/latest/userguide/sample-source-version.html
	// A tag (for example, refs/tags/mytagv1.0^{full-commit-SHA}).
	// A branch (for example, refs/heads/mydevbranch^{full-commit-SHA}).
	sourceVersion := e.HeadCommit.GetID()
	if ref := e.GetRef(); ref != "" {
		sourceVersion = fmt.Sprintf("%s^{%s}", ref, e.HeadCommit.GetID())
	}

	for _, w := range workflows {
		m, err := w.On.Push.MatchPaths(changedPaths)
		if err != nil {
			return nil, fmt.Errorf("invalid workflow: %w", err)
		}
		if m {
			for jobName, j := range w.Jobs {
				// TODO: return multiple builds
				return &codebuild.StartBuildInput{
					ProjectName:       aws.String("codebuild-runner"),
					SourceVersion:     aws.String(sourceVersion),
					BuildspecOverride: aws.String(j.Buildspec),
					EnvironmentVariablesOverride: []*types.EnvironmentVariable{
						{
							Name:  aws.String("RUNNER_WORKFLOW_NAME"),
							Value: aws.String(w.GetName()),
							Type:  types.EnvironmentVariableTypePlaintext,
						},
						{
							Name:  aws.String("RUNNER_JOB_NAME"),
							Value: aws.String(jobName),
							Type:  types.EnvironmentVariableTypePlaintext,
						},
						{
							Name:  aws.String("GITHUB_WEBHOOK_REF"),
							Value: e.Ref,
							Type:  types.EnvironmentVariableTypePlaintext,
						},
						{
							Name:  aws.String("GITHUB_WEBHOOK_COMMIT_ID"),
							Value: e.HeadCommit.ID,
							Type:  types.EnvironmentVariableTypePlaintext,
						},
					},
				}, nil
			}
		}
	}
	return nil, nil
}
