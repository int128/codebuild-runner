package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/int128/codebuild-runner/handler"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	lambda.Start(handle)
}

func handle(ctx context.Context, e events.SNSEvent) {
	for _, record := range e.Records {
		handleRecord(ctx, record)
	}
}

func handleRecord(ctx context.Context, r events.SNSEventRecord) {
	var e events.CodeBuildEvent
	if err := json.Unmarshal([]byte(r.SNS.Message), &e); err != nil {
		log.Printf("invalid SNS message: %+v", err)
		return
	}
	if err := handler.HandleCodeBuildEvent(ctx, e); err != nil {
		log.Printf("error: %+v", err)
		return
	}
}
