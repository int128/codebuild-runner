package usecases

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func CodeBuildEvent(ctx context.Context, e events.CodeBuildEvent) error {
	if e.DetailType == events.CodeBuildPhaseChangeDetailType {
		return phaseChangeEvent(ctx, e)
	}
	if e.DetailType == events.CodeBuildStateChangeDetailType {
		return statusChangeEvent(ctx, e)
	}
	return fmt.Errorf("unknown event detail type `%s`", e.DetailType)
}
