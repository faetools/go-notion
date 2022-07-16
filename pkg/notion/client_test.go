package notion_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/faetools/go-notion/pkg/docs"
	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

func TestClient(t *testing.T) {
	t.Parallel()
	t.Cleanup(gock.Off)

	ctx := context.Background()

	cli, err := notion.NewDefaultClient("")
	require.NoError(t, err)

	v := docs.NewVisitor(
		// get the document and at the same time check if the response has been parsed correctly
		docs.NewGetter(&responseTester{cli: cli, t: t}, nil),

		// don't do anything after having fetched the document, just continue
		func(p *notion.Page) error { return nil },
		func(blocks notion.Blocks) error { return nil },
		func(db *notion.Database) error { return nil },
		func(entries notion.Pages) error { return nil })

	assert.NoError(t, docs.Walk(ctx, v, docs.TypePage, fake.PageID))
}

type responseTester struct {
	cli *notion.Client

	t *testing.T
}

// GetNotionPage implements notion.Getter
func (g *responseTester) GetNotionPage(ctx context.Context, id notion.Id) (*notion.Page, error) {
	g.t.Helper()

	path := fmt.Sprintf("v1/pages/%s", id)
	fake.MockResponseTo(g.t, path)

	p, err := g.cli.GetNotionPage(ctx, id)
	require.NoError(g.t, err)

	assertResponseWellParsed(g.t, path, p)

	return p, nil
}

// GetAllBlocks implements notion.Getter
func (rt *responseTester) GetAllBlocks(ctx context.Context, id notion.Id) (notion.Blocks, error) {
	rt.t.Helper()

	path := fmt.Sprintf("v1/blocks/%s/children", id)
	fake.MockResponseTo(rt.t, path)

	resp, err := rt.cli.GetBlocks(ctx, id, &notion.GetBlocksParams{})
	require.NoError(rt.t, err)

	assertResponseWellParsed(rt.t, path, resp.JSON200)

	return resp.JSON200.Results, nil
}

// GetNotionDatabase implements notion.Getter
func (rt *responseTester) GetNotionDatabase(ctx context.Context, id notion.Id) (*notion.Database, error) {
	rt.t.Helper()

	switch id {
	case "d105edb4-586a-4dcc-aaa6-ea944eb8d864":
		// not the ID of the actual database
		return nil, docs.Skip
	}

	path := fmt.Sprintf("v1/databases/%s", id)
	fake.MockResponseTo(rt.t, path)

	resp, err := rt.cli.GetDatabase(ctx, id)
	require.NoError(rt.t, err)

	assertResponseWellParsed(rt.t, path, resp.JSON200)

	return resp.JSON200, nil
}

// GetAllDatabaseEntries implements notion.Getter
func (rt *responseTester) GetAllDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
	rt.t.Helper()

	path := fmt.Sprintf("v1/databases/%s/query", id)
	fake.MockResponseTo(rt.t, path)

	resp, err := rt.cli.QueryDatabase(ctx, id, notion.QueryDatabaseJSONRequestBody{})
	require.NoError(rt.t, err)

	assertResponseWellParsed(rt.t, path, resp.JSON200)

	return resp.JSON200.Results, nil
}

func assertResponseWellParsed(t *testing.T, path string, resp any) {
	t.Helper()

	respMarshalled, err := json.Marshal(resp)
	require.NoError(t, err)

	var expected, actual any
	assert.NoError(t, json.Unmarshal(fake.ResponseTo(t, path), &expected))
	assert.NoError(t, json.Unmarshal(respMarshalled, &actual))

	assert.Equal(t, expected, cleanTimestamps(actual), "result of GET %s was not well parsed", path)
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
