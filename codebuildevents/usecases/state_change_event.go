package usecases

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-github/v32/github"
	"github.com/int128/codebuild-runner/codebuildevents"
	"golang.org/x/oauth2"
)

func stateChangeEvent(ctx context.Context, e codebuildevents.CodeBuildEvent) error {
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

func calculateStatus(e codebuildevents.CodeBuildEvent) (*commitStatus, error) {
	commitID := findEnvironmentVariable(e, "GITHUB_WEBHOOK_HEADCOMMIT_ID")
	if commitID == "" {
		return nil, fmt.Errorf("could not find the commit id")
	}
	owner, repo, err := parseGitHubURL(e.Detail.AdditionalInformation.Source.Location)
	if err != nil {
		return nil, fmt.Errorf("could not determine GitHub URL: %w", err)
	}
	codeBuildURL, err := computeCodeBuildURL(e.Detail.BuildID)
	if err != nil {
		return nil, fmt.Errorf("could not determine CodeBuild URL: %w", err)
	}
	return &commitStatus{
		owner:  owner,
		repo:   repo,
		commit: commitID,
		repoStatus: github.RepoStatus{
			Context:     github.String("CodeBuild/example"),
			State:       github.String(determineCommitStatus(e.Detail.BuildStatus)),
			Description: github.String(string(e.Detail.BuildStatus)),
			TargetURL:   github.String(codeBuildURL),
		},
	}, nil
}

// arn:aws:codebuild:REGION:ACCOUNT:build/PROJECT:BUILD
var regexpBuildID = regexp.MustCompile(`^arn:aws:codebuild:(.+?):(.+?):build/(.+?):(.+?)$`)

// https://REGION.console.aws.amazon.com/codesuite/codebuild/ACCOUNT/projects/PROJECT/build/PROJECT:BUILD/config?region=REGION
func computeCodeBuildURL(buildID string) (string, error) {
	m := regexpBuildID.FindStringSubmatch(buildID)
	if m == nil {
		return "", fmt.Errorf("invalid build-id `%s`", buildID)
	}
	return fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codebuild/%s/projects/%s/build/%s:%s/config?region=%s",
		m[1], // region
		m[2], // account
		m[3], // project
		m[3], // project
		m[4], // build
		m[1], // region
	), nil
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

func determineCommitStatus(codeBuildPhaseStatus events.CodeBuildPhaseStatus) string {
	switch codeBuildPhaseStatus {
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
