package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/faetools/cgtools"
	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
)

func main() {
	ctx := context.Background()

	cli, err := notion.NewDefaultClient(os.Getenv("NOTION_TOKEN"))
	checkErr(err)

	g := &fakesGenerator{
		Generator: cgtools.NewOsGenerator(),
		cli:       cli,
	}

	g.genGetPageResponse(ctx, fake.PageID)
	g.genGetBlocksResponse(ctx, fake.PageID)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type fakesGenerator struct {
	*cgtools.Generator
	cli *notion.Client
}

func (g *fakesGenerator) genResponseFactory(
	urlPathFormat string,
	getResponse func(context.Context, notion.Id) (*http.Response, []byte),
	followUp func(context.Context, []byte),
) func(context.Context, notion.Id) {
	return func(ctx context.Context, id notion.Id) {
		urlPath := fmt.Sprintf(urlPathFormat, id)
		filePath := filepath.Join("../pkg/fake/", urlPath+".json")

		content, err := os.ReadFile(filePath)
		if err != nil {
			resp, body := getResponse(ctx, id)

			if want := resp.Request.URL.Path; want != urlPath {
				log.Fatalf("mismatched paths: %q vs %q", want, urlPath)
			}

			if resp.StatusCode != http.StatusOK {
				log.Fatalf("status for %q not OK but %s", urlPath, resp.Status)
			}

			content = body

			checkErr(g.WriteBytes(filePath, content))
		}

		if followUp != nil {
			followUp(ctx, content)
		}
	}
}

func (g *fakesGenerator) genGetPageResponse(ctx context.Context, id notion.Id) {
	g.genResponseFactory(
		"/v1/pages/%s",
		func(ctx context.Context, id notion.Id) (*http.Response, []byte) {
			resp, err := g.cli.GetPage(ctx, id)
			checkErr(err)

			return resp.HTTPResponse, resp.Body
		},
		nil,
	)(ctx, id)
}

func (g *fakesGenerator) genGetBlocksResponse(ctx context.Context, id notion.Id) {
	g.genResponseFactory(
		"/v1/blocks/%s/children",
		func(ctx context.Context, id notion.Id) (*http.Response, []byte) {
			resp, err := g.cli.GetBlocks(ctx, id, &notion.GetBlocksParams{})
			checkErr(err)

			return resp.HTTPResponse, resp.Body
		},
		func(ctx context.Context, b []byte) {
			var blockList notion.BlocksList
			checkErr(json.Unmarshal(b, &blockList))

			for _, block := range blockList.Results {
				if block.HasChildren {
					g.genGetBlocksResponse(ctx, notion.Id(block.Id))
				}
			}
		},
	)(ctx, id)
}
