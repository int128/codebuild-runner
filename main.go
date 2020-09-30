package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/int128/codebuild-runner/handler"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "GitHubWebhookFunction" {
		lambda.Start(handler.HandleGitHubWebhook)
		return
	}
	if functionName == "CodeBuildEventsFunction" {
		lambda.Start(handler.HandleCodeBuildSNSEvent)
		return
	}
	log.Fatalf("invalid FUNCTION_NAME: %s", functionName)
}
