package workflow

import (
	"fmt"

	"github.com/gobwas/glob"
)

// Workflow represents a workflow.
// https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-syntax-for-github-actions
type Workflow struct {
	// basename of workflow file
	basename string

	Name string         `yaml:"name"`
	On   On             `yaml:"on"`
	Jobs map[string]Job `yaml:"jobs"`
}

func (w Workflow) GetName() string {
	if w.Name != "" {
		return w.Name
	}
	return w.basename
}

type On struct {
	Push        Event `yaml:"push"`
	PullRequest Event `yaml:"pull_request"`
}

type Event struct {
	Paths    []string `yaml:"paths"`
	Branches []string `yaml:"branches"`
	Tags     []string `yaml:"tags"`
}

func (e Event) MatchPaths(paths []string) (bool, error) {
	for _, pattern := range e.Paths {
		g, err := glob.Compile(pattern)
		if err != nil {
			return false, fmt.Errorf("invalid path pattern `%s`: %w", pattern, err)
		}
		for _, changed := range paths {
			if g.Match(changed) {
				return true, nil
			}
		}
	}
	return false, nil
}

type Job struct {
	Name string `yaml:"name"`

	// path to buildspec.yml
	Buildspec string `yaml:"buildspec"`
}
