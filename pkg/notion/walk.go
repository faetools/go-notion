package notion

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

// ObjectType defines the type of object.
type ObjectType string

// Defines values for ObjectType.
const (
	ObjectTypePage            ObjectType = "page"
	ObjectTypeBlocks          ObjectType = "blocks"
	ObjectTypeDatabase        ObjectType = "database"
	ObjectTypeDatabaseEntries ObjectType = "database_entries"
)

// SkipPage is used as a return value from a Visitor to indicate that
// the page named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipPage = errors.New("skip this page") //nolint:go-lint

// SkipDatabase is used as a return value from a Visitor to indicate that
// the database named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDatabase = errors.New("skip this database") //nolint:go-lint

// Visitor traverses through notion documents.
// If the first result not empty, Walk visits each of the children (blocks or entries).
type Visitor interface {
	VisitPage(ctx context.Context, id Id) error

	// VisitBlocks visits any blocks associated with the ID.
	// Only return blocks you want to be walked.
	VisitBlocks(context.Context, Id) (Blocks, error)

	VisitDatabase(context.Context, Id) error

	VisitDatabaseEntries(context.Context, Id) (Pages, error)
}

// Walk traverses notion documents.
func Walk(ctx context.Context, v Visitor, tp ObjectType, id Id) error {
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		return nil
	})

	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		return nil
	})

	switch tp {
	case ObjectTypePage:
		if err := v.VisitPage(ctx, id); err != nil {
			if errors.Is(err, SkipPage) {
				return nil
			}

			return err
		}

		return Walk(ctx, v, ObjectTypeBlocks, id)
	case ObjectTypeBlocks:
		blocks, err := v.VisitBlocks(ctx, id)
		if err != nil {
			return err
		}

		for _, b := range blocks {
			switch b.Type {
			case BlockTypeChildPage:
				if err := Walk(ctx, v, ObjectTypePage, Id(b.Id)); err != nil {
					return err
				}
			case BlockTypeChildDatabase:
				// Unfortunately, notion does not tell us if this child database
				// has the same ID as the block ID or if a child database was just referenced.
				//
				// We're still calling Walk, the user will need to filter out such references
				// in their VisitDatabase method.
				if err := Walk(ctx, v, ObjectTypeDatabase, Id(b.Id)); err != nil {
					if errors.Is(err, SkipDatabase) {
						return nil
					}

					return err
				}
			default:
				if b.HasChildren {
					if err := Walk(ctx, v, ObjectTypeBlocks, Id(b.Id)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	case ObjectTypeDatabase:
		if err := v.VisitDatabase(ctx, id); err != nil {
			if errors.Is(err, SkipDatabase) {
				return nil
			}

			return err
		}

		return Walk(ctx, v, ObjectTypeDatabaseEntries, id)
	case ObjectTypeDatabaseEntries:
		entries, err := v.VisitDatabaseEntries(ctx, id)
		if err != nil {
			if errors.Is(err, SkipDatabase) {
				return nil
			}

			return err
		}

		for _, p := range entries {
			if err := Walk(ctx, v, ObjectTypePage, Id(p.Id)); err != nil {
				return err
			}
		}

		return nil
	default:
		return fmt.Errorf("unknown object type %q", tp)
	}
}
