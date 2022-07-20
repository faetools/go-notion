package notion_test

import . "github.com/faetools/go-notion/pkg/notion"

func validateDatabase(db *Database) error {
	if err := validateRichTexts(db.Title); err != nil {
		return err
	}

	if err := validateRichTexts(db.Description); err != nil {
		return err
	}

	return nil
}
