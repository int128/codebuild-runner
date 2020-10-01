package usecases

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func stateChangeEvent(ctx context.Context, e events.CodeBuildEvent) error {
	commitID := findEnvironmentVariable(e, "GITHUB_WEBHOOK_HEADCOMMIT_ID")
	if commitID == "" {
		return fmt.Errorf("could not find the commit id")
	}

	ghc := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})))
	_, _, err := ghc.Repositories.CreateStatus(ctx,
		"int128", "codebuild-runner", commitID,
		&github.RepoStatus{
			Context:     github.String("CodeBuild/example"),
			State:       github.String(determineCommitStatus(e)),
			Description: github.String(string(e.Detail.BuildStatus)),
		})
	if err != nil {
		return fmt.Errorf("could not create a commit status: %w", err)
	}
	return nil
}

func determineCommitStatus(e events.CodeBuildEvent) string {
	switch e.Detail.BuildStatus {
	case events.CodeBuildPhaseStatusQueued:
		return "pending"
	case events.CodeBuildPhaseStatusInProgress:
		return "pending"
	case events.CodeBuildPhaseStatusSucceeded:
		return "success"
	case events.CodeBuildPhaseStatusStopped:
		return "error"
	case events.CodeBuildPhaseStatusFailed:
		return "failure"
	case events.CodeBuildPhaseStatusFault:
		return "error"
	case events.CodeBuildPhaseStatusTimedOut:
		return "error"
	}
	return ""
}
