# go-notion

![Devtool version](https://img.shields.io/badge/Devtool-0.0.18-brightgreen.svg)
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

## Generate Database Property Values

You can use the `pkg/gen` package to generate helpers to get the values of a database entry.

To do this, create a package for each database and define the properties of each database.

For example, you could have three databases, `foo`, `bar`, and `blub`. For each, you create a package in the folder `databases`. Each package has a public variable called `Properties` of type `notion.PropertyMetaMap`

Then you just need to create a file with the following content in `databases`:

```go
package main

import (
	"log"

	"github.com/faetools/go-notion/pkg/gen"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
	"github.com/user/myrepo/databases/bar"
	"github.com/user/myrepo/databases/blub"
	"github.com/user/myrepo/databases/foo"
)

//go:generate go run gen.go

func main() {
	fs := afero.NewOsFs()

	for pkgName, props := range map[string]notion.PropertyMetaMap{
		"foo":  foo.Properties,
		"bar":  bar.Properties,
		"blub": blub.Properties,
	} {
		if err := gen.PropertyValues(fs, pkgName, props); err != nil {
			log.Fatal(err)
		}
	}
}
```

Run `go generate ./...` and your code will get generated.

## Contribution

Feel free to contribute to this repo by making a PR that changes the OpenAPI. We will then run the code generator.

Alternatively, feel free to add go code that will make the API easier to use.
