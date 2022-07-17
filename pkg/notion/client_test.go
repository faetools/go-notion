package notion_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cli, fsClient, err := fake.NewClient()
	assert.NoError(t, err)

	v := docs.NewVisitor(
		// get the document and at the same time check if the response has been parsed correctly
		&responseTester{cli: cli},

		// don't do anything after having fetched the document, just continue
		func(p *notion.Page) error { return nil },
		func(blocks notion.Blocks) error { return nil },
		func(db *notion.Database) error { return nil },
		func(entries notion.Pages) error { return nil })

	assert.NoError(t, docs.Walk(ctx, v, docs.TypePage, fake.PageID))

	assert.Empty(t, fsClient.Unseen())
}

type responseTester struct{ cli *notion.Client }

// GetNotionPage implements notion.Getter
func (rt *responseTester) GetNotionPage(ctx context.Context, id notion.Id) (*notion.Page, error) {
	resp, err := rt.cli.GetPage(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

// GetAllBlocks implements notion.Getter
func (rt *responseTester) GetAllBlocks(ctx context.Context, id notion.Id) (notion.Blocks, error) {
	resp, err := rt.cli.GetBlocks(ctx, id, &notion.GetBlocksParams{})
	if err != nil {
		return nil, err
	}

	return resp.JSON200.Results, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

// GetNotionDatabase implements notion.Getter
func (rt *responseTester) GetNotionDatabase(ctx context.Context, id notion.Id) (*notion.Database, error) {
	resp, err := rt.cli.GetDatabase(ctx, id)
	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode() == http.StatusNotFound:
		// not the ID of the actual database
		return nil, resp.JSON404
	}

	return resp.JSON200, validateResponseParsing(resp.HTTPResponse, resp.Body, resp.JSON200)
}

// GetAllDatabaseEntries implements notion.Getter
func (rt *responseTester) GetAllDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
	resp, err := rt.cli.QueryDatabase(ctx, id, notion.QueryDatabaseJSONRequestBody{})
	if err != nil {
		return nil, err
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
