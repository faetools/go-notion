package notion_test

import (
	"testing"

	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
)

func TestMaps(t *testing.T) {
	t.Parallel()

	// these get set differently by the code generator
	// unfortunately we have to set them manually after running it
	assert.Equal(t, notion.PropertyValueMap(nil), notion.Page{}.Properties)
	assert.Equal(t, notion.PropertyMetaMap(nil), notion.Database{}.Properties)
}
