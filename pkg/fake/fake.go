package fake

import (
	_ "embed" // fake responses

	"github.com/faetools/go-notion/pkg/notion"
)

// PageID is the page ID of our example page.
const PageID notion.Id = "96245c8f178444a482ad1941127c3ec3"

// GetPageResponse is the response we got by calling GetPage on the example page.
//
//go:embed get-page.json
var GetPageResponse string

// GetBlocksResponse is the response we got by calling GetBlocks on the example page.
//
//go:embed get-blocks.json
var GetBlocksResponse string
