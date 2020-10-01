package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/int128/codebuild-runner/codebuildevents"
	"github.com/int128/codebuild-runner/codebuildevents/usecases"
)

func Handle(ctx context.Context, e events.SNSEvent) {
	for _, record := range e.Records {
		err := handleCodeBuildSNSEventRecord(ctx, record)
		if err != nil {
			log.Printf("error: %s", err)
		}
	}
}

func handleCodeBuildSNSEventRecord(ctx context.Context, r events.SNSEventRecord) error {
	e, err := codebuildevents.ParseCodeBuildEvent(r.SNS.Message)
	if err != nil {
		return fmt.Errorf("could not parse the SNS message: %w", err)
	}
	if err := usecases.CodeBuildEvent(ctx, e); err != nil {
		return fmt.Errorf("could not process the event: %w", err)
	}
	return nil
}
