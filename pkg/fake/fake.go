package fake

import (
	"embed" // fake responses
	"net/http"
	"testing"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/kjk/common/require"
	"gopkg.in/h2non/gock.v1"
)

// PageID is the ID of our example page.
// It can be viewed here: https://ancient-gibbon-2cd.notion.site/Example-Page-96245c8f178444a482ad1941127c3ec3
const PageID notion.Id = "96245c8f178444a482ad1941127c3ec3"

// Responses contains a number of responses we have generated.
//
//go:embed v1
var Responses embed.FS

// ResponseTo returns the body of a fake response to a GET request for the stated path.
func ResponseTo(t *testing.T, path string) []byte {
	t.Helper()

	body, err := Responses.ReadFile(path + ".json")
	require.NoError(t, err)

	return body
}

// MockResponseTo mocks a response to a particular path.
func MockResponseTo(t *testing.T, path string) {
	t.Helper()

	gock.New(notion.DefaultServer).
		Path(path).
		Reply(http.StatusOK).
		SetHeader("Content-Type", "json").
		BodyString(string(ResponseTo(t, path)))
}
