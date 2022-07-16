package notion_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

func TestClient(t *testing.T) {
	t.Parallel()
	t.Cleanup(gock.Off)

	cli, err := notion.NewDefaultClient("")
	require.NoError(t, err)

	clientTester{Client: cli}.getPage(t, fake.PageID)
}

type clientTester struct{ *notion.Client }

func (c clientTester) getPage(t *testing.T, id notion.Id) {
	t.Helper()

	path := fmt.Sprintf("v1/pages/%s", id)
	fake.MockResponseTo(t, path)

	p, err := c.GetNotionPage(context.Background(), id)
	require.NoError(t, err)

	assertResponseWellParsed(t, path, p)

	c.getBlocks(t, id)
}

func cleanUUID(id notion.Id) notion.Id {
	return notion.Id(uuid.MustParse(string(id)).String())
}

func (c clientTester) getBlocks(t *testing.T, id notion.Id) {
	t.Helper()

	id = cleanUUID(id)

	path := fmt.Sprintf("v1/blocks/%s/children", cleanUUID(id))
	fake.MockResponseTo(t, path)

	resp, err := c.GetBlocks(context.Background(), id, &notion.GetBlocksParams{})
	require.NoError(t, err)

	assertResponseWellParsed(t, path, resp.JSON200)

	for _, b := range resp.JSON200.Results {
		switch b.Type {
		case notion.BlockTypeChildPage:
			c.getPage(t, notion.Id(b.Id))
		case notion.BlockTypeChildDatabase:
			// unfortunately, notion does not tell us
			// if this child database has the same ID as the block ID
			// or if this child database is referenced
			if b.Id == "d105edb4-586a-4dcc-aaa6-ea944eb8d864" {
				continue
			}

			c.getDatabase(t, notion.Id(b.Id))
		default:
			if b.HasChildren {
				c.getBlocks(t, notion.Id(b.Id))
			}
		}
	}
}

func (c clientTester) getDatabase(t *testing.T, id notion.Id) {
	t.Helper()

	path := fmt.Sprintf("v1/databases/%s", id)
	fake.MockResponseTo(t, path)

	resp, err := c.GetDatabase(context.Background(), id)
	require.NoError(t, err)

	assertResponseWellParsed(t, path, resp.JSON200)

	// TODO database entries
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
