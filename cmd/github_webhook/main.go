package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/int128/codebuild-runner/handler"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	lambda.Start(handle)
}

func handle(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if r.HTTPMethod == "POST" && r.Path == "/github" {
		code, err := handler.HandleGitHubWebhook(ctx, r.MultiValueQueryStringParameters, r.Body)
		if err != nil {
			log.Printf("error: %+v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: code,
				Headers:    map[string]string{"content-type": "text/plain"},
				Body:       fmt.Sprintf("Error %d", code),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: code,
			Headers:    map[string]string{"content-type": "text/plain"},
			Body:       "OK",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       "Not Found",
	}, nil
}
