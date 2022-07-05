package notion

// TitleProperty represents a Title property.
var TitleProperty = PropertyMeta{
	Id:    "title", // must be this
	Type:  PropertyTypeTitle,
	Title: &map[string]interface{}{},
}

// NewRichTexts creates a RichTexts object with the desired content.
func NewRichTexts(content string) RichTexts {
	return RichTexts{NewRichText(content)}
}

// NewRichText creates a RichText object with the desired content.
func NewRichText(content string) RichText {
	return RichText{
		Type:        RichTextTypeText,
		PlainText:   content,
		Text:        &Text{Content: content},
		Annotations: Annotations{Color: ColorDefault},
	}
}

// GetNames returns names of all selected options.
func (opts PropertyOptions) GetNames() []string {
	names := make([]string, len(opts))

	for i, sel := range opts {
		names[i] = sel.Name
	}

	return names
}
