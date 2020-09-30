package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/int128/codebuild-runner/codebuildevents/usecases"
)

func Handle(ctx context.Context, e events.SNSEvent) {
	for _, record := range e.Records {
		handleCodeBuildSNSEventRecord(ctx, record)
	}
}

func handleCodeBuildSNSEventRecord(ctx context.Context, r events.SNSEventRecord) {
	var e events.CodeBuildEvent
	if err := json.Unmarshal([]byte(r.SNS.Message), &e); err != nil {
		log.Printf("invalid SNS message: %+v", err)
		return
	}
	if err := usecases.CodeBuildEvent(ctx, e); err != nil {
		log.Printf("error: %+v", err)
		return
	}
}
