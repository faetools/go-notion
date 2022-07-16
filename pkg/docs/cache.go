package docs

import (
	"sync"

	"github.com/faetools/go-notion/pkg/notion"
)

// Cache caches notion documents.
type Cache struct {
	mu sync.Mutex

	pages           cacheEntry[*notion.Page]
	blockChildren   cacheEntry[notion.Blocks]
	databases       cacheEntry[*notion.Database]
	databaseEntries cacheEntry[notion.Pages]
}

// NewCache returns a new cache.
// It can be filled using its method LoadOrStore methods.
func NewCache() *Cache {
	return &Cache{
		pages:           cacheEntry[*notion.Page]{},
		blockChildren:   cacheEntry[notion.Blocks]{},
		databases:       cacheEntry[*notion.Database]{},
		databaseEntries: cacheEntry[notion.Pages]{},
	}
}

type cacheEntry[T any] map[notion.Id]T

func (e *cacheEntry[T]) loadOrStore(id notion.Id, get func(notion.Id) T) T {
	if p, ok := (*e)[id]; ok {
		return p
	}

	(*e)[id] = get(id)

	return (*e)[id]
}

// LoadOrStorePage returns the page if present.
// Otherwise, it stores and returns the value it gets from calling get.
func (c *Cache) LoadOrStorePage(
	id notion.Id, get func(id notion.Id) *notion.Page,
) *notion.Page {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.pages.loadOrStore(id, get)
}

// LoadOrStoreBlocks returns the blocks if present.
// Otherwise, it stores and returns the value it gets from calling get.
func (c *Cache) LoadOrStoreBlocks(
	parentID notion.Id, get func(id notion.Id) notion.Blocks,
) notion.Blocks {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.blockChildren.loadOrStore(parentID, get)
}

// LoadOrStoreDatabase returns the database if present.
// Otherwise, it stores and returns the value it gets from calling get.
func (c *Cache) LoadOrStoreDatabase(
	id notion.Id, get func(id notion.Id) *notion.Database,
) *notion.Database {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.databases.loadOrStore(id, get)
}

// LoadOrStoreDatabaseEntries returns the database entries if present.
// Otherwise, it stores and returns the value it gets from calling get.
func (c *Cache) LoadOrStoreDatabaseEntries(
	parentID notion.Id, get func(id notion.Id) notion.Pages,
) notion.Pages {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries := c.databaseEntries.loadOrStore(parentID, get)

	for _, p := range entries {
		if c.pages[notion.Id(p.Id)] == nil {
			c.pages[notion.Id(p.Id)] = &p
		}
	}

	return entries
}
