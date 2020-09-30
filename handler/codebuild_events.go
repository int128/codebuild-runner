package handler

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func HandleCodeBuildEvent(_ context.Context, e events.CodeBuildEvent) error {
	je := json.NewEncoder(os.Stderr)
	je.SetIndent("", "  ")
	if err := je.Encode(&e); err != nil {
		return err
	}
	return nil
}
