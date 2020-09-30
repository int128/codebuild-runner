package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/int128/codebuild-runner/handler"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	lambda.Start(handle)
}

func handle(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	code, err := handler.HandleGitHubWebhook(ctx, r.Headers, r.Body)
	if err != nil {
		log.Printf("error: %+v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: code,
		}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
	}, nil
}
