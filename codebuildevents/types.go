package codebuildevents

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type CodeBuildEvent struct {
	events.CodeBuildEvent

	// CodeBuildEvent type has key `detail-type` but actual SNS message has key `detailType`.
	// Here we use `detailType`.
	DetailType string `json:"detailType"`
}

func ParseCodeBuildEvent(m string) (CodeBuildEvent, error) {
	var e CodeBuildEvent
	if err := json.Unmarshal([]byte(m), &e); err != nil {
		return CodeBuildEvent{}, fmt.Errorf("invalid json: %w", err)
	}
	return e, nil
}
