package notion

import (
	"fmt"
)

// TitleProperty represents a Title property.
var TitleProperty = PropertyMeta{
	Id:    "title", // must be this
	Type:  PropertyTypeTitle,
	Title: &map[string]interface{}{},
}

func NewPage(title string, parent *Parent) Page {
	return Page{
		Object: "page",
		Properties: PropertyValueMap{
			"title": PropertyValue{
				Type:  PropertyTypeTitle,
				Title: NewRichTextsP(title),
			},
		},
		Parent: parent,
	}
}

func NewParagraph(txt string) *Paragraph {
	return &Paragraph{
		Color:    ColorDefault,
		RichText: NewRichTexts(txt),
	}
}

// NewRichTexts creates a RichTexts object with the desired content.
func NewRichTexts(content string) RichTexts {
	return RichTexts{NewRichText(content)}
}

// NewRichTextsP creates a pointer to a RichTexts object with the desired content.
func NewRichTextsP(content string) *RichTexts {
	rt := NewRichTexts(content)
	return &rt
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

func mapSlice[T, R any](collection []T, iteratee func(T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item)
	}

	return result
}

// GetNames returns names of all selected options.
func (vals SelectValues) GetNames() []string {
	return mapSlice(vals, func(val SelectValue) string { return val.Name })
}

// GetByName returns the select value with the respective name.
func (vals SelectValues) GetByName(name string) (SelectValue, bool) {
	for _, val := range vals {
		if val.Name == name {
			return val, true
		}
	}

	return SelectValue{}, false
}

// GetIDs returns the UUIDs of all references.
func (refs References) GetIDs() []UUID {
	return mapSlice(refs, func(ref Reference) UUID { return ref.Id })
}

// ID return the ID of the page or database that it is linked to.
func (l LinkToPage) ID() UUID {
	switch l.Type {
	case LinkToPageTypePageId:
		return *l.PageId
	case LinkToPageTypeDatabaseId:
		return *l.DatabaseId
	default:
		panic("invalid LinkToPage of type " + l.Type)
	}
}

// ID returns the ID of the object that was mentioned.
func (m Mention) ID() UUID {
	switch m.Type {
	case MentionTypeDatabase:
		return m.Database.Id
	case MentionTypePage:
		return m.Page.Id
	case MentionTypeUser:
		return m.User.Id
	default:
		return "<no ID>"
	}
}

// ID returns the ID of the parent.
func (p Parent) ID() UUID {
	switch p.Type {
	case ParentTypeBlockId:
		return *p.BlockId
	case ParentTypeDatabaseId:
		return *p.DatabaseId
	case ParentTypePageId:
		return *p.PageId
	case ParentTypeWorkspace:
		return "workspace"
	default:
		return UUID(fmt.Sprintf("<invalid parent type %s>", p.Type))
	}
}

func (s SelectValue) GetColor() Color {
	if s.Color == nil {
		return ColorDefault
	}

	return *s.Color
}

func (s SelectValue) GetID() string {
	if s.Id == nil {
		return ""
	}

	return *s.Id
}
