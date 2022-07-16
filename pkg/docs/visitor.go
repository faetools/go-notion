package docs

import (
	"context"

	"github.com/faetools/go-notion/pkg/notion"
)

type (
	// PageVisit defines what is to be done when visiting a page.
	PageVisit func(p *notion.Page) error

	// BlocksVisit defines what is to be done when visiting blocks.
	BlocksVisit func(blocks notion.Blocks) error

	// DatabaseVisit defines what is to be done when visiting a database.
	DatabaseVisit func(db *notion.Database) error

	// DatabaseEntriesVisit defines what is to be done when visiting database entries.
	DatabaseEntriesVisit func(entries notion.Pages) error
)

// Visitor traverses through notion documents.
// If the first result not empty, Walk visits each of the children (blocks or entries).
// The implementor is responsible for fetching and caching.
type Visitor interface {
	VisitPage(context.Context, notion.Id) error
	VisitBlocks(context.Context, notion.Id) (notion.Blocks, error)
	VisitDatabase(context.Context, notion.Id) error
	VisitDatabaseEntries(context.Context, notion.Id) (notion.Pages, error)
}

type visitor struct {
	*Getter

	atPage            PageVisit
	atBlocks          BlocksVisit
	atDatabase        DatabaseVisit
	atDatabaseEntries DatabaseEntriesVisit
}

// NewVisitor returns a visitor.
func NewVisitor(
	g *Getter,
	atPage PageVisit, atBlocks BlocksVisit,
	atDatabase DatabaseVisit, atDatabaseEntries DatabaseEntriesVisit,
) Visitor {
	return &visitor{
		Getter:            g,
		atPage:            atPage,
		atBlocks:          atBlocks,
		atDatabase:        atDatabase,
		atDatabaseEntries: atDatabaseEntries,
	}
}

func (v *visitor) VisitPage(ctx context.Context, id notion.Id) error {
	if v.atPage == nil {
		return nil
	}

	p, err := v.GetPage(ctx, id)
	if err != nil {
		return err
	}

	return v.atPage(p)
}

func (v *visitor) VisitBlocks(ctx context.Context, id notion.Id) (notion.Blocks, error) {
	if v.atBlocks == nil {
		return nil, nil
	}

	bs, err := v.GetBlocks(ctx, id)
	if err != nil {
		return nil, err
	}

	return bs, v.atBlocks(bs)
}

func (v *visitor) VisitDatabase(ctx context.Context, id notion.Id) error {
	if v.atDatabase == nil {
		return nil
	}

	db, err := v.GetDatabase(ctx, id)
	if err != nil {
		return err
	}

	return v.atDatabase(db)
}

func (v *visitor) VisitDatabaseEntries(ctx context.Context, id notion.Id) (notion.Pages, error) {
	if v.atDatabaseEntries == nil {
		return nil, nil
	}

	entries, err := v.GetDatabaseEntries(ctx, id)
	if err != nil {
		return nil, err
	}

	return entries, v.atDatabaseEntries(entries)
}
