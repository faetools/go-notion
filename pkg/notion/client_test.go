package notion_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/faetools/go-notion/pkg/fake"
	. "github.com/faetools/go-notion/pkg/notion"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

func TestClient(t *testing.T) {
	t.Parallel()
	t.Cleanup(gock.Off)

	cli, err := NewDefaultClient("")
	require.NoError(t, err)

	// Walk(context.Background(), &clientTester{Client: cli, t: t}, ObjectTypePage, fake.PageID)

	clientTester{Client: cli}.getPage(t, fake.PageID)
}

type clientTester struct {
	*Client

	t *testing.T
}

// func (g *fakesGenerator) VisitPage(ctx context.Context, id notion.Id) error {
// 	_ = g.getResponse("/v1/pages/%s", id, func(id notion.Id) (*http.Response, []byte) {
// 		resp, err := g.cli.GetPage(ctx, id)
// 		checkErr(err)

// 		return resp.HTTPResponse, resp.Body
// 	})

// 	return nil
// }

// func (g *fakesGenerator) VisitBlocks(ctx context.Context, id notion.Id) (notion.Blocks, error) {
// 	body := g.getResponse("/v1/blocks/%s/children", id, func(id notion.Id) (*http.Response, []byte) {
// 		resp, err := g.cli.GetBlocks(ctx, id, &notion.GetBlocksParams{})
// 		checkErr(err)

// 		return resp.HTTPResponse, resp.Body
// 	})

// 	var list notion.BlocksList
// 	checkErr(json.Unmarshal(body, &list))

// 	return list.Results, nil
// }

// func (g *fakesGenerator) VisitDatabase(ctx context.Context, id notion.Id) error {
// 	switch id {
// 	case "d105edb4-586a-4dcc-aaa6-ea944eb8d864":
// 		// not the ID of the actual database
// 		return notion.SkipDatabase
// 	}

// 	_ = g.getResponse("/v1/databases/%s", id, func(id notion.Id) (*http.Response, []byte) {
// 		resp, err := g.cli.GetDatabase(ctx, id)
// 		checkErr(err)

// 		return resp.HTTPResponse, resp.Body
// 	})

// 	return nil
// }

// func (g *fakesGenerator) VisitDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
// 	body := g.getResponse("/v1/databases/%s/query", id, func(id notion.Id) (*http.Response, []byte) {
// 		resp, err := g.cli.QueryDatabase(ctx, id, notion.QueryDatabaseJSONRequestBody{})
// 		checkErr(err)

// 		return resp.HTTPResponse, resp.Body
// 	})

// 	var list notion.PagesList
// 	checkErr(json.Unmarshal(body, &list))

// 	return list.Results, nil
// }

func (c clientTester) getPage(t *testing.T, id Id) {
	t.Helper()

	path := fmt.Sprintf("v1/pages/%s", id)
	fake.MockResponseTo(t, path)

	p, err := c.GetNotionPage(context.Background(), id)
	require.NoError(t, err)

	assertResponseWellParsed(t, path, p)

	c.getBlocks(t, id)
}

func cleanUUID(id Id) Id {
	return Id(uuid.MustParse(string(id)).String())
}

func (c clientTester) getBlocks(t *testing.T, id Id) {
	t.Helper()

	path := fmt.Sprintf("v1/blocks/%s/children", id)
	fake.MockResponseTo(t, path)

	resp, err := c.GetBlocks(context.Background(), id, &GetBlocksParams{})
	require.NoError(t, err)

	assertResponseWellParsed(t, path, resp.JSON200)

	for _, b := range resp.JSON200.Results {
		switch b.Type {
		case BlockTypeChildPage:
			c.getPage(t, Id(b.Id))
		case BlockTypeChildDatabase:
			// unfortunately, notion does not tell us
			// if this child database has the same ID as the block ID
			// or if this child database is referenced
			if b.Id == "d105edb4-586a-4dcc-aaa6-ea944eb8d864" {
				continue
			}

			c.getDatabase(t, Id(b.Id))
		default:
			if b.HasChildren {
				c.getBlocks(t, Id(b.Id))
			}
		}
	}
}

func (c clientTester) getDatabase(t *testing.T, id Id) {
	t.Helper()

	path := fmt.Sprintf("v1/databases/%s", id)
	fake.MockResponseTo(t, path)

	resp, err := c.GetDatabase(context.Background(), id)
	require.NoError(t, err)

	assertResponseWellParsed(t, path, resp.JSON200)

	c.queryDatabase(t, id)
}

func (c clientTester) queryDatabase(t *testing.T, id Id) {
	t.Helper()

	path := fmt.Sprintf("v1/databases/%s/query", id)
	fake.MockResponseTo(t, path)

	resp, err := c.QueryDatabase(context.Background(), id, QueryDatabaseJSONRequestBody{})
	require.NoError(t, err)

	assertResponseWellParsed(t, path, resp.JSON200)

	// TODO each entry
	// c.queryDatabase(t, id)
}

func assertResponseWellParsed(t *testing.T, path string, resp any) {
	t.Helper()

	respMarshalled, err := json.Marshal(resp)
	require.NoError(t, err)

	var expected, actual any
	assert.NoError(t, json.Unmarshal(fake.ResponseTo(t, path), &expected))
	assert.NoError(t, json.Unmarshal(respMarshalled, &actual))

	assert.Equal(t, expected, cleanTimestamps(actual),
		"result of GET %s was not well parsed", path)
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
