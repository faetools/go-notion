package docs

import (
	"context"

	"github.com/faetools/go-notion/pkg/notion"
)

// getterWithCache gets and caches notion documents.
type getterWithCache struct {
	cli   notion.Getter
	cache *Cache
}

// NewGetterWithCache returns a new Getter.
func NewGetterWithCache(cli notion.Getter, cache *Cache) notion.Getter {
	if cache == nil {
		cache = NewCache()
	}

	return &getterWithCache{cli: cli, cache: cache}
}

// GetNotionPage returns the notion page or an error.
// An internal cache is used, if available.
func (c *getterWithCache) GetNotionPage(ctx context.Context, id notion.Id) (*notion.Page, error) {
	var err error

	return c.cache.LoadOrStorePage(id,
		func(id notion.Id) (p *notion.Page) {
			p, err = c.cli.GetNotionPage(ctx, id)
			return p
		}), err
}

// GetAllBlocks returns all blocks of a given page or block.
// An internal cache is used, if available.
func (c *getterWithCache) GetAllBlocks(ctx context.Context, parentID notion.Id) (notion.Blocks, error) {
	var err error

	return c.cache.LoadOrStoreBlocks(parentID,
		func(id notion.Id) (bs notion.Blocks) {
			bs, err = c.cli.GetAllBlocks(ctx, parentID)
			return bs
		}), err
}

// GetNotionDatabase returns the notion database or an error.
// An internal cache is used, if available.
func (c *getterWithCache) GetNotionDatabase(ctx context.Context, id notion.Id) (*notion.Database, error) {
	var err error

	return c.cache.LoadOrStoreDatabase(id,
		func(id notion.Id) (db *notion.Database) {
			db, err = c.cli.GetNotionDatabase(ctx, id)
			return db
		}), err
}

// GetAllDatabaseEntries returns all database entries or an error.
// An internal cache is used, if available.
func (c *getterWithCache) GetAllDatabaseEntries(ctx context.Context, parentID notion.Id) (notion.Pages, error) {
	var err error

	return c.cache.LoadOrStoreDatabaseEntries(parentID,
		func(parentID notion.Id) (entries notion.Pages) {
			entries, err = c.cli.GetAllDatabaseEntries(ctx, parentID)
			return entries
		}), err
}
