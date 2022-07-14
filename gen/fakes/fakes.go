package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/faetools/cgtools"
	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
)

func main() {
	ctx := context.Background()

	g := cgtools.NewOsGenerator()

	cl, err := notion.NewDefaultClient(os.Getenv("NOTION_TOKEN"))
	checkErr(err)

	for _, v := range []struct {
		getContent func() []byte
		fileName   string
	}{
		{func() []byte {
			resp, err := cl.GetPage(ctx, fake.PageID)
			checkErr(err)

			return resp.Body
		}, "get-page.json"},
		{func() []byte {
			resp, err := cl.GetBlocks(ctx, fake.PageID, &notion.GetBlocksParams{})
			checkErr(err)

			return resp.Body
		}, "get-blocks.json"},
	} {
		path := filepath.Join("../pkg/fake/", v.fileName)

		if err := g.WriteBytes(path, v.getContent()); err != nil {
			log.Fatal(err)
		}

	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
