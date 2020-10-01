package usecases

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func phaseChangeEvent(ctx context.Context, e events.CodeBuildEvent) error {
	log.Printf("Detail=%+v", e.Detail)
	return nil
}

func findEnvironmentVariable(e events.CodeBuildEvent, key string) string {
	for _, v := range e.Detail.AdditionalInformation.Environment.EnvironmentVariables {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}
