# envparse

[![GoDev][godev-image]][godev-url]
[![Build Status][build-image]][build-url]
[![Coverage Status][coverage-image]][coverage-url]
[![Go Report Card][goreport-image]][goreport-url]

Envparse is library for parse ENV variables written in [Go][go].

Taken from moby/docker Dockerfile parser.

See original [moby parser][moby-parser].

Parses env variables in form of:
```
<key> <value>
<key>=<value> ...
```

See also original [ENV documentation][docs-docker-env]

## Quick-start

```go
package main

import (
  "bytes"
  "fmt"

  "github.com/go-x-pkg/envparse"
)

func main() {
  env, err := envparse.Parse("LD_LIBRARY_PATH=/usr/lib:/usr/nvidia/cuda/lib64 PATH=/usr/bin:/bin")

  fmt.Println(env, err)
}
```

[godev-image]: https://img.shields.io/badge/go.dev-reference-5272B4?logo=go&logoColor=white
[godev-url]: https://pkg.go.dev/github.com/go-x-pkg/envparse

[build-image]: https://travis-ci.com/go-x-pkg/envparse.svg?branch=master
[build-url]: https://travis-ci.com/go-x-pkg/envparse

[coverage-image]: https://coveralls.io/repos/github/go-x-pkg/envparse/badge.svg?branch=master
[coverage-url]: https://coveralls.io/github/go-x-pkg/envparse?branch=master

[goreport-image]: https://goreportcard.com/badge/github.com/go-x-pkg/envparse
[goreport-url]: https://goreportcard.com/report/github.com/go-x-pkg/envparse

[go]: http://golang.org/
[moby-parser]: https://github.com/moby/buildkit/tree/master/frontend/dockerfile/parser
[docs-docker-env]: https://docs.docker.com/engine/reference/builder/#env
