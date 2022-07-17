package fake

import (
	"github.com/faetools/client"
	"github.com/faetools/go-notion/pkg/notion"
)

// NewClient returns a new notion client returning fake results.
func NewClient() (*notion.Client, error) {
	return notion.NewDefaultClient("", client.WithHTTPClient(NewDoer()))
}
