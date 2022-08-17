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
