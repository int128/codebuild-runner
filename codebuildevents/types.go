package codebuildevents

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type CodeBuildEvent struct {
	events.CodeBuildEvent

	// CodeBuildEvent type has key `detail-type` but actual SNS message has key `detailType`.
	// Here we use `detailType`.
	DetailType string `json:"detailType"`

	Detail CodeBuildEventDetail `json:"detail"`
}

var regexpSourceVersionRefSHA = regexp.MustCompile(`^.+?\^{(.+?)}$`)
var regexpSourceVersionSHA = regexp.MustCompile(`^[0-9a-fA-F]+$`)

func (e CodeBuildEvent) GetCommitSHA() string {
	v := e.Detail.AdditionalInformation.SourceVersion
	m := regexpSourceVersionRefSHA.FindStringSubmatch(v)
	if m != nil {
		return m[1]
	}
	if regexpSourceVersionSHA.MatchString(v) {
		return v
	}
	return ""
}

func (e CodeBuildEvent) GetEnvironmentVariable(key string) string {
	for _, v := range e.Detail.AdditionalInformation.Environment.EnvironmentVariables {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

type CodeBuildEventDetail struct {
	events.CodeBuildEventDetail

	AdditionalInformation CodeBuildEventAdditionalInformation `json:"additional-information"`
}

type CodeBuildEventAdditionalInformation struct {
	events.CodeBuildEventAdditionalInformation

	// actual SNS message has key source-version
	SourceVersion string `json:"source-version"`
}

func ParseCodeBuildEvent(m string) (CodeBuildEvent, error) {
	var e CodeBuildEvent
	if err := json.Unmarshal([]byte(m), &e); err != nil {
		return CodeBuildEvent{}, fmt.Errorf("invalid json: %w", err)
	}
	return e, nil
}
