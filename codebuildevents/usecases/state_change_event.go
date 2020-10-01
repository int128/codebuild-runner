package usecases

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func stateChangeEvent(ctx context.Context, e events.CodeBuildEvent) error {
	s, err := calculateStatus(e)
	if err != nil {
		return fmt.Errorf("could not calculate the commit status: %w", err)
	}

	ghc := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})))
	if _, _, err := ghc.Repositories.CreateStatus(ctx, s.owner, s.repo, s.commit, &s.repoStatus); err != nil {
		return fmt.Errorf("could not create a commit status: %w", err)
	}
	return nil
}

type commitStatus struct {
	owner      string
	repo       string
	commit     string
	repoStatus github.RepoStatus
}

func calculateStatus(e events.CodeBuildEvent) (*commitStatus, error) {
	commitID := findEnvironmentVariable(e, "GITHUB_WEBHOOK_HEADCOMMIT_ID")
	if commitID == "" {
		return nil, fmt.Errorf("could not find the commit id")
	}
	owner, repo, err := parseGitHubURL(e.Detail.AdditionalInformation.Source.Location)
	if err != nil {
		return nil, fmt.Errorf("could not find GitHub URL: %w", err)
	}
	return &commitStatus{
		owner:  owner,
		repo:   repo,
		commit: commitID,
		repoStatus: github.RepoStatus{
			Context:     github.String("CodeBuild/example"),
			State:       github.String(determineCommitStatus(e)),
			Description: github.String(string(e.Detail.BuildStatus)),
		},
	}, nil
}

func parseGitHubURL(l string) (string, string, error) {
	if !strings.HasPrefix(l, "https://github.com/") {
		return "", "", fmt.Errorf("not GitHub URL `%s`", l)
	}
	pair := strings.TrimSuffix(strings.TrimPrefix(l, "https://github.com/"), ".git")
	e := strings.SplitN(pair, "/", 2)
	if len(e) < 2 {
		return "", "", fmt.Errorf("could not find owner/repo from `%s`", l)
	}
	return e[0], e[1], nil
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
