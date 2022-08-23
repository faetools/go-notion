package notion

import "encoding/json"

// UnmarshalJSON overrides the default JSON handling for PageOrDatabase.
func (pdb *PageOrDatabase) UnmarshalJSON(b []byte) error {
	page := &Page{}
	if err := json.Unmarshal(b, page); err == nil && page.Object == "page" {
		pdb.Page = page
		return nil
	}

	db := &Database{}
	if err := json.Unmarshal(b, db); err != nil {
		return err
	}

	pdb.Database = db
	return nil
}

// MarshalJSON overrides the default JSON handling for PageOrDatabase.
func (pdb PageOrDatabase) MarshalJSON() ([]byte, error) {
	if pdb.Page != nil {
		return json.Marshal(pdb.Page)
	}

	return json.Marshal(pdb.Database)
}
