package usecases

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func phaseChangeEvent(ctx context.Context, e events.CodeBuildEvent) error {
	commitID := findEnvironmentVariable(e, "GITHUB_WEBHOOK_HEADCOMMIT_ID")
	if commitID == "" {
		return fmt.Errorf("could not find the commit id")
	}

	ghc := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})))
	_, _, err := ghc.Repositories.CreateStatus(ctx,
		"int128", "codebuild-runner", commitID,
		&github.RepoStatus{
			Context:     github.String("CodeBuild/example"),
			State:       github.String("pending"),
			TargetURL:   github.String(e.Detail.BuildID),
			Description: github.String(string(e.Detail.CompletedPhase)),
		})
	if err != nil {
		return fmt.Errorf("could not create a commit status: %w", err)
	}
	return nil
}

func findEnvironmentVariable(e events.CodeBuildEvent, key string) string {
	for _, v := range e.Detail.AdditionalInformation.Environment.EnvironmentVariables {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}
