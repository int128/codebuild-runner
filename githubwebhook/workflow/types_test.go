package workflow

import "testing"

func TestEvent_MatchPaths(t *testing.T) {
	// TODO: check whether GitHub webhook changes have head-slash or not
	e := Event{
		Paths: []string{"**.go"},
	}
	m, err := e.MatchPaths([]string{"foo/main.go"})
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if !m {
		t.Errorf("MatchPaths wants true but false")
	}
}
