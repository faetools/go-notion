package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/faetools/cgtools"
)

// We are generating code for this repository.

//go:generate go run ./gen.go

const (
	typesPath = "../pkg/notion/types.gen.go"
	// testPageID notion.UUID = "96245c8f178444a482ad1941127c3ec3"
)

func main() {
	g := cgtools.NewOsGenerator()

	types, err := os.ReadFile(typesPath)
	if err != nil {
		log.Fatal(err)
	}

	types = bytes.Replace(types,
		[]byte("Properties PropertyMetas"),
		[]byte("Properties PropertyMetaMap"), 1)

	types = bytes.Replace(types,
		[]byte("Properties PropertyValues"),
		[]byte("Properties PropertyValueMap"), 1)

	if err := g.WriteBytes(typesPath, types); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}
