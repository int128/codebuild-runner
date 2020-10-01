package usecases

import (
	"context"
	"log"

	"github.com/int128/codebuild-runner/codebuildevents"
)

func phaseChangeEvent(ctx context.Context, e codebuildevents.CodeBuildEvent) error {
	log.Printf("Detail=%+v", e.Detail)
	return nil
}
