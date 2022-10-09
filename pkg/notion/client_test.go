package notion_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/faetools/go-notion-example/fake"
	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/notion"
	. "github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
)

var reBlock = regexp.MustCompile("^v1/blocks/[a-z0-9-]+.json$")

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	v := docs.NewVisitor(
		// get the document and at the same time check if the response has been parsed correctly
		&responseTester{cli: fake.NotionClient},

		// don't do anything after having fetched the document, just continue
		func(p *Page) error { return nil },
		func(blocks Blocks) error { return nil },
		func(db *Database) error { return nil },
		func(entries Pages) error { return nil })

	if err := docs.Walk(ctx, v, docs.TypePage, fake.ExamplePageID); err != nil {
		t.Fatal(err)
	}
}

type responseTester struct{ cli *Client }

// GetNotionPage implements Getter
func (rt *responseTester) GetNotionPage(ctx context.Context, id Id) (*Page, error) {
	resp, err := rt.cli.GetPage(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validatePage(resp.JSON200); err != nil {
		return nil, err
	}

	return resp.JSON200, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

var (
	maxPageSizeInt          = 100
	maxPageSize    PageSize = PageSize(maxPageSizeInt)
)

// GetAllBlocks implements Getter
func (rt *responseTester) GetAllBlocks(ctx context.Context, id Id) (Blocks, error) {
	resp, err := rt.cli.GetBlocks(ctx, id, &GetBlocksParams{
		PageSize: &maxPageSize,
	})
	if err != nil {
		return nil, err
	}

	if err := validateBlocks(resp.JSON200.Results); err != nil {
		return nil, err
	}

	return resp.JSON200.Results, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

// GetNotionDatabase implements Getter
func (rt *responseTester) GetNotionDatabase(ctx context.Context, id Id) (*Database, error) {
	if id == "d105edb4-586a-4dcc-aaa6-ea944eb8d864" {
		// not the ID of the actual database
		return nil, &notion.Error{Status: http.StatusNotFound}
	}

	resp, err := rt.cli.GetDatabase(ctx, id)
	if err != nil {
		return nil, err
	}

	db := resp.JSON200

	if err := validateDatabase(db); err != nil {
		return nil, err
	}

	return db, validateResponseParsing(resp.HTTPResponse, resp.Body, db)
}

// GetAllDatabaseEntries implements Getter
func (rt *responseTester) GetAllDatabaseEntries(ctx context.Context, id Id) (Pages, error) {
	resp, err := rt.cli.QueryDatabase(ctx, id, QueryDatabaseJSONRequestBody{
		PageSize: maxPageSizeInt,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range resp.JSON200.Results {
		if err := validatePage(&p); err != nil {
			return nil, err
		}
	}

	return resp.JSON200.Results, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

// validateResponseParsing check whether we have parsed the response correctly,
// or if we missed or added fields.
func validateResponseParsing(resp *http.Response, body []byte, parsed any) error {
	path := resp.Request.URL.Path

	gotMarshalled, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("marshalling the parsed response for %q: %w", path, err)
	}

	// unmarshal both as general as we can
	var want, got any

	if err := json.Unmarshal(body, &want); err != nil {
		return fmt.Errorf("unmarshalling response body for %q: %w", path, err)
	}

	if err := json.Unmarshal(gotMarshalled, &got); err != nil {
		return fmt.Errorf("unmarshalling back the response we marshalled for %q: %w", path, err)
	}

	t := &testBuffer{}
	assert.Equal(t, want, cleanTimestamps(got), "result of GET %s was not well parsed", path)
	return t.Err()
}

type testBuffer struct{ bytes.Buffer }

func (t *testBuffer) Errorf(format string, args ...interface{}) {
	t.WriteString(fmt.Sprintf(format, args...))
}

func (t testBuffer) Err() error {
	if t.Len() == 0 {
		return nil
	}

	return errors.New(t.String())
}

func cleanTimestamps(a any) any {
	switch v := a.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return v
		}

		return t.Format(layoutTime)
	case map[string]any:
		for key, val := range v {
			v[key] = cleanTimestamps(val)
		}

		return v
	case []any:
		out := make([]any, len(v))

		for i, elem := range v {
			out[i] = cleanTimestamps(elem)
		}

		return out
	default:
		return a
	}
}
