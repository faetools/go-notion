package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/faetools/cgtools"
	_client "github.com/faetools/client"
	"github.com/faetools/go-notion/pkg/client"
	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/faetools/kit/terminal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()

	fakes := afero.NewBasePathFs(afero.NewOsFs(), "../pkg/fake/")

	files, err := client.NewFSClient(afero.NewIOFS(fakes), nil)
	checkErr(err)

	g := &fakesGenerator{
		gen:   cgtools.NewGenerator(fakes),
		files: client.NewCachingClient(files),
		cli: NewRequestValidator(http.DefaultClient, func(req *http.Request) error {
			if req.URL.Path == "/v1/databases/d105edb4-586a-4dcc-aaa6-ea944eb8d864" {
				return docs.Skip
			}

			return nil
		}),
		wg: &errgroup.Group{},
	}

	cli, err := notion.NewDefaultClient(os.Getenv("NOTION_TOKEN"), _client.WithHTTPClient(g))
	checkErr(err)

	v := docs.NewVisitor(
		cli,

		// don't do anything after having fetched the document, just continue
		func(p *notion.Page) error { return nil },
		func(blocks notion.Blocks) error { return nil },
		func(db *notion.Database) error { return nil },
		func(entries notion.Pages) error { return nil })

	checkErr(docs.Walk(ctx, v, docs.TypePage, fake.PageID))

	// create files
	checkErr(g.wg.Wait())

	// remove unneccessary files
	for _, path := range files.Unseen() {
		if filepath.Dir(path) == "." {
			continue
		}

		checkErr(fakes.Remove(path))

		terminal.Printf(aurora.Red, "  • %v was removed\n", path)
	}

	// remove unneccessary folders
	checkErr(afero.Walk(fakes, ".", func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			return nil
		}

		files, err := afero.ReadDir(fakes, path)
		checkErr(err)

		if len(files) == 0 {
			checkErr(fakes.RemoveAll(path))

			terminal.Printf(aurora.Red, "  • %v was removed\n", path)
			return fs.SkipDir
		}

		return nil
	}))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
