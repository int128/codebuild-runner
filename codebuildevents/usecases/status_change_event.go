package usecases

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func statusChangeEvent(ctx context.Context, e events.CodeBuildEvent) error {
	return nil
}
