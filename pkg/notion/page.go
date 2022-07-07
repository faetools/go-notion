package notion

// Title returns the page title.
func (p Page) Title() string {
	return p.Properties.title()
}

// Title returns the title of the page.
func (props PropertyValueMap) title() string {
	for _, prop := range props {
		if prop.Title != nil {
			return prop.Title.Content()
		}
	}

	return ""
}
