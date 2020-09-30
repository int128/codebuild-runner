TARGETS := codebuild_events github_webhook
VERSION := v0.0.0

STACK_NAME := codebuild-runner
SAM_S3_BUCKET_NAME ?= sam-codebuild-runner
AWS_REGION ?= ap-northeast-1

.PHONY: all
all: $(TARGETS)
codebuild_events:
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/codebuild_events
github_webhook:
	GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/github_webhook

.PHONY: clean
clean:
	-rm $(TARGETS) packaged.yaml

.PHONY: run
run: github_webhook
	sam local start-api

packaged.yaml: $(TARGETS) template.yaml
	sam package --template-file template.yaml --output-template-file $@ --s3-bucket $(SAM_S3_BUCKET_NAME)

.PHONY: deploy
deploy: packaged.yaml
	sam deploy --template-file $< --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --region $(AWS_REGION)

.PHONY: create-bucket
create-bucket:
	aws s3 mb s3://$(SAM_S3_BUCKET_NAME) --region $(AWS_REGION)
