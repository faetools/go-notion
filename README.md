# go-notion

![Devtool version](https://img.shields.io/badge/Devtool-0.0.17-brightgreen.svg)
![Maintainer](https://img.shields.io/badge/team-firestarters-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/faetools/go-notion)](https://goreportcard.com/report/github.com/faetools/go-notion)
[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/faetools/go-notion)

## About

This repository contains an [OpenAPI definition](api/openapi.yaml) of Notion's API based on [their documentation](https://developers.notion.com/).

Faetools uses a code generator to generate go code based on that OpenAPI.

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

	resp, err := cl.GetNotionPage(ctx, "[page ID]")
	if err != nil {
		log.Fatal(err)
	}

	...
}
```

## Contribution

Feel free to contribute to this repo by making a PR that changes the OpenAPI. We will then run the code generator.

Alternatively, feel free to add go code that will make the API easier to use.
