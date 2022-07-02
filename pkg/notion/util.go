package notion

// TitleProperty represents a Title property.
var TitleProperty = PropertyMeta{
	Id:    "title", // must be this
	Name:  "Title",
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
