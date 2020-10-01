package usecases

import "testing"

func Test_computeCodeBuildURL(t *testing.T) {
	codeBuildURL, err := computeCodeBuildURL("arn:aws:codebuild:ap-northeast-1:123456789012:build/codebuild-runner:de5b2e03-a44f-4217-b684-c7b03aede7b7")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	want := "https://ap-northeast-1.console.aws.amazon.com/codesuite/codebuild/123456789012/projects/codebuild-runner/build/codebuild-runner:de5b2e03-a44f-4217-b684-c7b03aede7b7/config?region=ap-northeast-1"
	if want != codeBuildURL {
		t.Errorf("CodeBuildURL wants %s but was %s", want, codeBuildURL)
	}
}
