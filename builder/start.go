package builder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
)

func Start(ctx context.Context, input *codebuild.StartBuildInput) (*codebuild.StartBuildOutput, error) {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load SDK config: %w", err)
	}
	b := codebuild.NewFromConfig(cfg)
	out, err := b.StartBuild(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("could not start a CodeBuild build: %w", err)
	}
	return out, nil
}
