package docs

import (
	"context"
	"errors"
	"fmt"

	"github.com/faetools/go-notion/pkg/notion"
)

// Type defines the type of document.
type Type string

// Defines values for Type.
const (
	TypePage            Type = "page"
	TypeBlocks          Type = "blocks"
	TypeDatabase        Type = "database"
	TypeDatabaseEntries Type = "database_entries"
)

// Skip is used as a return value from a Visitor to indicate that
// the page or database named in the call is to be skipped.
// It is not returned as an error by any other function.
var Skip = errors.New("skip this page or database") //nolint:go-lint

// Walk traverses notion documents.
func Walk(ctx context.Context, v Visitor, tp Type, id notion.Id) error {
	switch tp {
	case TypePage:
		if err := v.VisitPage(ctx, id); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		return Walk(ctx, v, TypeBlocks, id)
	case TypeBlocks:
		blocks, err := v.VisitBlocks(ctx, id)
		if err != nil {
			return err
		}

		for _, b := range blocks {
			switch b.Type {
			case notion.BlockTypeChildPage:
				if err := Walk(ctx, v, TypePage, notion.Id(b.Id)); err != nil {
					return err
				}
			case notion.BlockTypeChildDatabase:
				// Unfortunately, notion does not tell us if this child database
				// has the same ID as the block ID or if a child database was just referenced.
				//
				// We're still calling Walk, the user will need to filter out such references
				// in their VisitDatabase method.
				if err := Walk(ctx, v, TypeDatabase, notion.Id(b.Id)); err != nil {
					return err
				}
			default:
				if b.HasChildren {
					if err := Walk(ctx, v, TypeBlocks, notion.Id(b.Id)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	case TypeDatabase:
		if err := v.VisitDatabase(ctx, id); err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		return Walk(ctx, v, TypeDatabaseEntries, id)
	case TypeDatabaseEntries:
		entries, err := v.VisitDatabaseEntries(ctx, id)
		if err != nil {
			if errors.Is(err, Skip) {
				return nil
			}

			return err
		}

		for _, p := range entries {
			if err := Walk(ctx, v, TypePage, notion.Id(p.Id)); err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("unknown object type %q", tp)
	}
}
