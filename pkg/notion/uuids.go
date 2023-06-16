package notion

type UUIDs []UUID

func (ids UUIDs) Strings() []string {
	if len(ids) == 0 {
		return nil
	}

	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = string(id)
	}

	return strs
}
