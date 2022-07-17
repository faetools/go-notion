package fake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/faetools/client"
	"github.com/faetools/go-notion/pkg/notion"
)

type doer struct{}

// NewDoer returns a new client that returns fake responses.
// They are based on real responses we once got from the example page, blocks, subpages and subdatabases.
// Example page: https://www.notion.so/Example-Page-96245c8f178444a482ad1941127c3ec3
func NewDoer() client.HTTPRequestDoer { return &doer{} }

var header = map[string][]string{
	client.ContentType: {client.MIMEApplicationJSON},
}

func (c doer) Do(req *http.Request) (*http.Response, error) {
	status := http.StatusOK

	body, err := responses.ReadFile(req.URL.Path[1:] + ".json")
	if err != nil {
		status = http.StatusNotFound

		body, _ = json.Marshal(notion.ErrorResponse{
			Status:  http.StatusNotFound,
			Code:    fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound)),
			Message: err.Error(),
			Object:  "error",
		})
	}

	return &http.Response{
		StatusCode:    status,
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(body)),
		Request:       req,
		ContentLength: int64(len(body)),
	}, nil
}
