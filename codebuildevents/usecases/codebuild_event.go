package usecases

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/int128/codebuild-runner/codebuildevents"
)

func CodeBuildEvent(ctx context.Context, e codebuildevents.CodeBuildEvent) error {
	if e.DetailType == events.CodeBuildPhaseChangeDetailType {
		return phaseChangeEvent(ctx, e)
	}
	if e.DetailType == events.CodeBuildStateChangeDetailType {
		return stateChangeEvent(ctx, e)
	}
	return fmt.Errorf("unknown event detail type `%s`", e.DetailType)
}
