package notion_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/faetools/go-notion/pkg/fake"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetPage(t *testing.T) {
	t.Parallel()
	t.Cleanup(gock.Off)

	ctx := context.Background()

	gock.New("https://api.notion.com").
		Path(fmt.Sprintf("/v1/pages/%s", fake.PageID)).
		Reply(http.StatusOK).
		SetHeader("Content-Type", "json").
		BodyString(fake.GetPageResponse)

	cli, err := notion.NewDefaultClient("")
	require.NoError(t, err)

	p, err := cli.GetNotionPage(ctx, fake.PageID)
	require.NoError(t, err)

	m, err := json.Marshal(p)
	require.NoError(t, err)

	assert.JSONEq(t, fake.GetPageResponse, string(m))
}
