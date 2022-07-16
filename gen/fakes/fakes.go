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
	"github.com/faetools/go-notion/pkg/docs"
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

	if err := docs.Walk(ctx, g, docs.TypePage, fake.PageID); err != nil {
		log.Fatal(err)
	}
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

func (g *fakesGenerator) VisitPage(ctx context.Context, id notion.Id) error {
	_ = g.getResponse("/v1/pages/%s", id, func(id notion.Id) (*http.Response, []byte) {
		resp, err := g.cli.GetPage(ctx, id)
		checkErr(err)

		return resp.HTTPResponse, resp.Body
	})

	return nil
}

func (g *fakesGenerator) VisitBlocks(ctx context.Context, id notion.Id) (notion.Blocks, error) {
	body := g.getResponse("/v1/blocks/%s/children", id, func(id notion.Id) (*http.Response, []byte) {
		resp, err := g.cli.GetBlocks(ctx, id, &notion.GetBlocksParams{})
		checkErr(err)

		return resp.HTTPResponse, resp.Body
	})

	var list notion.BlocksList
	checkErr(json.Unmarshal(body, &list))

	return list.Results, nil
}

func (g *fakesGenerator) VisitDatabase(ctx context.Context, id notion.Id) error {
	switch id {
	case "d105edb4-586a-4dcc-aaa6-ea944eb8d864":
		// not the ID of the actual database
		return docs.Skip
	}

	_ = g.getResponse("/v1/databases/%s", id, func(id notion.Id) (*http.Response, []byte) {
		resp, err := g.cli.GetDatabase(ctx, id)
		checkErr(err)

		return resp.HTTPResponse, resp.Body
	})

	return nil
}

func (g *fakesGenerator) VisitDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
	body := g.getResponse("/v1/databases/%s/query", id, func(id notion.Id) (*http.Response, []byte) {
		resp, err := g.cli.QueryDatabase(ctx, id, notion.QueryDatabaseJSONRequestBody{})
		checkErr(err)

		return resp.HTTPResponse, resp.Body
	})

	var list notion.PagesList
	checkErr(json.Unmarshal(body, &list))

	return list.Results, nil
}

func (g *fakesGenerator) getResponse(urlPathFormat string, id notion.Id,
	getResponse func(id notion.Id) (*http.Response, []byte),
) []byte {
	urlPath := fmt.Sprintf(urlPathFormat, id)
	filePath := filepath.Join("../pkg/fake/", urlPath+".json")

	content, err := os.ReadFile(filePath)
	if err == nil {
		return content
	}

	resp, body := getResponse(id)

	if want := resp.Request.URL.Path; want != urlPath {
		log.Fatalf("mismatched paths: %q vs %q", want, urlPath)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status for %q not OK but %s", urlPath, resp.Status)
	}

	checkErr(g.WriteBytes(filePath, body))

	return body
}
