# codebuild-runner

This is a Lambda function to integrate GitHub Webhook and AWS CodeBuild.

You can define workflows in your repository.

```yaml
# /.codebuild/workflows/build.yaml
on:
  push:
    paths:
      - .codebuild/workflows/build.yaml
      - '**.go'

jobs:
  build:
    buildspec: /buildspec.yml
```
