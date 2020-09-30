TARGET := codebuild-runner
VERSION := v0.0.0

STACK_NAME := codebuild-runner
SAM_S3_BUCKET_NAME ?= sam-codebuild-runner
AWS_REGION ?= ap-northeast-1

$(TARGET):
	GOOS=linux GOARCH=amd64 go build -o $@ .

.PHONY: clean
clean:
	-rm $(TARGET) packaged.yaml

.PHONY: run
run: $(TARGET)
	sam local start-api

packaged.yaml: $(TARGET) template.yaml
	sam package --template-file template.yaml --output-template-file $@ --s3-bucket $(SAM_S3_BUCKET_NAME)

.PHONY: deploy
deploy: packaged.yaml
	sam deploy --template-file $< --stack-name $(STACK_NAME) --capabilities CAPABILITY_IAM --region $(AWS_REGION)

.PHONY: create-bucket
create-bucket:
	aws s3 mb s3://$(SAM_S3_BUCKET_NAME) --region $(AWS_REGION)
