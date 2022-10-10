package fake

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/faetools/client"
	"github.com/faetools/clueo"
	"github.com/faetools/go-notion-example/fake/export"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/spf13/afero"
)

//go:embed responses.zip
var responses []byte

// ExamplePageID is the ID of our example page.
// It can be viewed here: https://ancient-gibbon-2cd.notion.site/Example-Page-96245c8f178444a482ad1941127c3ec3
const ExamplePageID notion.Id = "96245c8f-1784-44a4-82ad-1941127c3ec3"

// ExamplePageIDWithUnsupported is the ID of our example page that has unsupported elements.
// It can be viewed here: https://ancient-gibbon-2cd.notion.site/Example-Page-with-unsupported-elements-5de9f82ff2284386937fb3490c62f6e5
const ExamplePageIDWithUnsupported notion.Id = "5de9f82f-f228-4386-937f-b3490c62f6e5"

var ExamplePageName = func() string {
	p, err := NotionClient.GetNotionPage(context.Background(), ExamplePageID)
	if err != nil {
		panic(err)
	}

	return p.Title()
}()

var Client = func() *http.Client {
	fs := afero.NewMemMapFs()

	if err := export.UnzipInto(responses, fs); err != nil {
		panic(err)
	}

	return &http.Client{
		Transport: clueo.NewFSReader(afero.NewIOFS(fs)),
	}
}()

var NotionClient = func() *notion.Client {
	cli, err := notion.NewDefaultClient("", client.WithHTTPClient(Client))
	if err != nil {
		panic(err)
	}

	return cli
}()
