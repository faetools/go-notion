package notion

func (rs References) GetUUIDs() UUIDs {
	if rs == nil {
		return nil
	}

	urls := make(UUIDs, len(rs))
	for i, r := range rs {
		urls[i] = r.Id
	}

	return urls
}
