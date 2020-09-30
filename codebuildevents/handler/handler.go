package handler

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
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
	if err := handleCodeBuildEvent(ctx, e); err != nil {
		log.Printf("error: %+v", err)
		return
	}
}

func handleCodeBuildEvent(_ context.Context, e events.CodeBuildEvent) error {
	je := json.NewEncoder(os.Stderr)
	je.SetIndent("", "  ")
	if err := je.Encode(&e); err != nil {
		return err
	}
	return nil
}
