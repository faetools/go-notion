package main

import (
	"bytes"
	"log"
	"os"

	"github.com/faetools/cgtools"
)

// We are fixing auto-generated code.

const typesPath = "../pkg/notion/types.gen.go"

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
}
