package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
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
	e, err := parseCodeBuildEvent(r.SNS.Message)
	if err != nil {
		return fmt.Errorf("could not parse the SNS message: %w", err)
	}
	if err := usecases.CodeBuildEvent(ctx, e); err != nil {
		return fmt.Errorf("could not process the event: %w", err)
	}
	return nil
}

func parseCodeBuildEvent(m string) (events.CodeBuildEvent, error) {
	var e struct {
		events.CodeBuildEvent

		// CodeBuildEvent type has key `detail-type` but actual SNS message has key `detailType`.
		// Here we use `detailType` if `detail-type` is empty.
		DetailType string `json:"detailType"`
	}
	if err := json.Unmarshal([]byte(m), &e); err != nil {
		return events.CodeBuildEvent{}, fmt.Errorf("invalid json: %w", err)
	}
	if e.CodeBuildEvent.DetailType == "" {
		e.CodeBuildEvent.DetailType = e.DetailType
	}
	return e.CodeBuildEvent, nil
}
