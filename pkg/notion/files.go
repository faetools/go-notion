package notion

func (fs Files) GetURLs() []string {
	if fs == nil {
		return nil
	}

	urls := make([]string, len(fs))
	for i, f := range fs {
		urls[i] = f.URL()
	}

	return urls
}
