# Harborman

Harborman is used to transform and forward Harbor's webhook messages to the
corresponding webhook services.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/scholar7r/harborman)
[![reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/scholar7r/harborman)](https://pkg.go.dev/github.com/scholar7r/harborman)
[![license](https://img.shields.io/badge/license-GPLv2-brightgreen.svg)](https://github.com/scholar7r/harborman/blob/HEAD/LICENSE)
[![codecov](https://codecov.io/gh/scholar7r/harborman/graph/badge.svg?token=QTS2wjp3pb)](https://codecov.io/gh/scholar7r/harborman)

## Configuration

Below is a simple configuration file template. You can name it
`harborman.yaml` and place it in the same directory as the
`docker-compose.yaml` file; it will be mounted as a runtime configuration file.

```yaml
debug: false
listen: ":80"
authorization: "..."
notifiers:
  - type: lark
    url: "https://..."
    authorization: "..."
  - type: discord
    url: "https://..."
```

## Make Contributions

Learn [Contributing to a project](https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project).

This project following the [conventional branch](https://conventional-branch.github.io) and [conventional commit](https://conventionalcommits.org).

## Coverage Trace

![Coverage Trace](https://codecov.io/gh/scholar7r/harborman/graphs/sunburst.svg?token=QTS2wjp3pb)

## License

This software is released under the GPLv2 license.
