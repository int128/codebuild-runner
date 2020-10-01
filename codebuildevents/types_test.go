package codebuildevents

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_parseCodeBuildEvent(t *testing.T) {
	e, err := ParseCodeBuildEvent(
		// actual SNS message received via email
		`{
  "account": "123456789012",
  "detailType": "CodeBuild Build Phase Change",
  "region": "ap-northeast-1",
  "source": "aws.codebuild",
  "time": "2020-10-01T00:04:31Z",
  "notificationRuleArn": "arn:aws:codestar-notifications:ap-northeast-1:123456789012:notificationrule/97cd247ce6257610f632ce4a53e95e8e8b8123b3",
  "detail": {
    "completed-phase": "SUBMITTED",
    "project-name": "codebuild-runner",
    "build-id": "arn:aws:codebuild:ap-northeast-1:123456789012:build/codebuild-runner:de5b2e03-a44f-4217-b684-c7b03aede7b7",
    "completed-phase-context": "[]",
    "additional-information": {
      "cache": {
        "type": "NO_CACHE"
      },
      "build-number": 7.0,
      "timeout-in-minutes": 60.0,
      "build-complete": false,
      "initiator": "codebuild-runner-GitHubWebhookFunctionRole-1234567890/codebuild-runner-GitHubWebhookFunction-1234567890",
      "build-start-time": "Oct 1, 2020 12:04:31 AM",
      "source": {
        "report-build-status": false,
        "location": "https://github.com/int128/codebuild-runner.git",
        "git-clone-depth": 1.0,
        "type": "GITHUB",
        "git-submodules-config": {
          "fetch-submodules": false
        }
      },
      "source-version": "refs/heads/master^{75614758eace7c1d9d032b8980854ab486068f3b}",
      "artifact": {
        "location": ""
      },
      "environment": {
        "image": "aws/codebuild/amazonlinux2-x86_64-standard:3.0",
        "privileged-mode": false,
        "image-pull-credentials-type": "CODEBUILD",
        "compute-type": "BUILD_GENERAL1_SMALL",
        "type": "LINUX_CONTAINER",
        "environment-variables": [
          {
            "name": "GITHUB_WEBHOOK_REF",
            "type": "PLAINTEXT",
            "value": "refs/heads/master"
          },
          {
            "name": "GITHUB_WEBHOOK_HEADCOMMIT_ID",
            "type": "PLAINTEXT",
            "value": "75614758eace7c1d9d032b8980854ab486068f3b"
          }
        ]
      },
      "logs": {
        "deep-link": "https://console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#logEvent:group=null;stream=null"
      },
      "phases": [
        {
          "phase-context": [],
          "start-time": "Oct 1, 2020 12:04:31 AM",
          "end-time": "Oct 1, 2020 12:04:31 AM",
          "duration-in-seconds": 0.0,
          "phase-type": "SUBMITTED",
          "phase-status": "SUCCEEDED"
        },
        {
          "start-time": "Oct 1, 2020 12:04:31 AM",
          "phase-type": "QUEUED"
        }
      ],
      "queued-timeout-in-minutes": 480.0
    },
    "completed-phase-status": "SUCCEEDED",
    "completed-phase-duration-seconds": 0.0,
    "version": "1",
    "completed-phase-start": "Oct 1, 2020 12:04:31 AM",
    "completed-phase-end": "Oct 1, 2020 12:04:31 AM"
  },
  "resources": [
    "arn:aws:codebuild:ap-northeast-1:123456789012:build/codebuild-runner:de5b2e03-a44f-4217-b684-c7b03aede7b7"
  ],
  "additionalAttributes": {}
}`)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if e.DetailType != events.CodeBuildPhaseChangeDetailType {
		t.Errorf("DetailType wants `%s` but was `%s`", events.CodeBuildPhaseChangeDetailType, e.DetailType)
	}
	if e.Detail.AdditionalInformation.SourceVersion != "refs/heads/master^{75614758eace7c1d9d032b8980854ab486068f3b}" {
		t.Errorf("SourceVersion wants %s but was %s",
			"refs/heads/master^{75614758eace7c1d9d032b8980854ab486068f3b}",
			e.Detail.AdditionalInformation.SourceVersion)
	}
}
