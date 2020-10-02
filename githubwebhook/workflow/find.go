package workflow

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

func Find(ctx context.Context, owner, repo, commitSHA string) ([]*Workflow, error) {
	r, err := find(ctx, owner, repo, commitSHA)
	if err != nil {
		return nil, fmt.Errorf("could not find: %w", err)
	}
	return parseFindResponse(r)
}

func parseFindResponse(r *findResponse) ([]*Workflow, error) {
	var workflows []*Workflow
	for _, entry := range r.Repository.Object.Commit.File.Object.Tree.Entries {
		if !strings.HasSuffix(entry.Name, ".yaml") {
			continue
		}
		workflow, err := parseWorkflow(entry.Object.Blob.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid workflow %s", entry.Name)
		}
		workflow.basename = strings.TrimSuffix(entry.Name, ".yaml")
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func parseWorkflow(s string) (*Workflow, error) {
	var w Workflow
	if err := yaml.Unmarshal([]byte(s), &w); err != nil {
		return nil, fmt.Errorf("could not unmarshal yaml: %w", err)
	}
	return &w, nil
}

type findResponse struct {
	Repository struct {
		Object struct {
			Commit struct {
				File struct {
					Object struct {
						Tree struct {
							Entries []struct {
								Name   string
								Object struct {
									Blob struct {
										Text string
									} `graphql:"... on Blob"`
								}
							}
						} `graphql:"... on Tree"`
					}
				} `graphql:"file(path: $path)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(oid: $commitSHA)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func find(ctx context.Context, owner, repo, commitSHA string) (*findResponse, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)
	var r findResponse
	v := map[string]interface{}{
		"owner":     githubv4.String(owner),
		"repo":      githubv4.String(repo),
		"commitSHA": githubv4.GitObjectID(commitSHA),
		"path":      githubv4.String("/.codebuild/workflows"),
	}
	if err := client.Query(ctx, &r, v); err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	return &r, nil
}
