package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	codebuildEventsHandler "github.com/int128/codebuild-runner/codebuildevents/handler"
	gitHubWebhookHandler "github.com/int128/codebuild-runner/githubwebhook/handler"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	// switch the role by ROLE environment variable to
	// use single binary for multiple functions
	role := os.Getenv("ROLE")
	if role == "GitHubWebhookFunction" {
		lambda.Start(gitHubWebhookHandler.Handle)
		return
	}
	if role == "CodeBuildEventsFunction" {
		lambda.Start(codebuildEventsHandler.Handle)
		return
	}
	log.Fatalf("invalid ROLE=%s", role)
}
