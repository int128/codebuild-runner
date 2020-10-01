package usecases

import (
	"context"
	"log"

	"github.com/int128/codebuild-runner/codebuildevents"
)

func phaseChangeEvent(ctx context.Context, e codebuildevents.CodeBuildEvent) error {
	log.Printf("Detail=%+v", e.Detail)
	return nil
}

func findEnvironmentVariable(e codebuildevents.CodeBuildEvent, key string) string {
	for _, v := range e.Detail.AdditionalInformation.Environment.EnvironmentVariables {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}
