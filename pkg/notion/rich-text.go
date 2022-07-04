package notion

import (
	"fmt"
	"strings"
)

// RawContent returns the raw content of all the rich texts.
func (ts RichTexts) rawContent() string {
	s := make([]string, len(ts))
	for i, t := range ts {
		s[i] = t.rawContent()
	}

	return strings.Join(s, "\n")
}

// RawContent returns the content of the rich text object.
// NOTE: At the moment, only really implemented for text objects.
func (t RichText) rawContent() string {
	switch t.Type {
	case RichTextTypeText:
		return t.Text.Content
	case RichTextTypeMention:
		return fmt.Sprintf("%#v", *t.Mention)
	case RichTextTypeEquation:
		return fmt.Sprintf("%#v", *t.Equation)
	default:
		return fmt.Sprintf("unknown RichText type %q", t.Type)
	}
}
