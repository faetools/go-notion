# go-notion

![Devtool version](https://img.shields.io/badge/Devtool-0.0.18-brightgreen.svg)
![Maintainer](https://img.shields.io/badge/team-firestarters-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/faetools/go-notion)](https://goreportcard.com/report/github.com/faetools/go-notion)
[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/faetools/go-notion)

## About

This repository contains an [OpenAPI definition](api/openapi.yaml) of Notion's API based on [their documentation](https://developers.notion.com/) as well as a go library to use the API.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/faetools/go-notion/pkg/notion"
)


func main() {
	ctx := context.Background()

	cl, err := notion.NewDefaultClient("[bearer token]")
	if err != nil {
		log.Fatal(err)
	}

	p, err := cl.GetNotionPage(ctx, "[page ID]")
	if err != nil {
		log.Fatal(err)
	}

	...
}
```

## Expansions

There are several expansions (work in progress):

- [ ] [go-notion-codegen](https://github.com/faetools/go-notion-codegen): Generates go code for your databases.
- [ ] [notion-to-goldmark](https://github.com/faetools/notion-to-goldmark): Transforms notion blocks into [goldmark](https://github.com/yuin/goldmark) nodes.
- [ ] [notion-to-md](https://github.com/faetools/notion-to-goldmark): Transforms notion blocks into markdown.

## Contribution

We use a code generator to generate go code based on the OpenAPI. In addition to the auto generated code, we added a number of convenience methods.

Feel free to contribute to this repo by making a PR that changes the OpenAPI. We will then run the code generator to generate respective go code.

Alternatively, feel free to add to the manually written convenience methods.
