package notion

import (
	"errors"
	"fmt"
)

var errNoType = errors.New("could not determine type of block")

// Title returns the title of the block.
func (b Block) Title() string {
	switch b.Type {
	case BlockTypeChildDatabase:
		return b.ChildDatabase.Title
	case BlockTypeChildPage:
		return b.ChildPage.Title
	default:
		return fmt.Sprintf("<no title defined for block type %q>", b.Type)
	}
}

func (b *Block) Validate() error {
	b.Object = "block"

	seenType := false

	for tp, isNil := range map[BlockType]bool{
		BlockTypeAudio:            b.Audio == nil,
		BlockTypeBookmark:         b.Bookmark == nil,
		BlockTypeBreadcrumb:       b.Breadcrumb == nil,
		BlockTypeBulletedListItem: b.BulletedListItem == nil,
		BlockTypeCallout:          b.Callout == nil,
		BlockTypeChildDatabase:    b.ChildDatabase == nil,
		BlockTypeChildPage:        b.ChildPage == nil,
		BlockTypeCode:             b.Code == nil,
		BlockTypeColumn:           b.Column == nil,
		BlockTypeColumnList:       b.ColumnList == nil,
		BlockTypeDivider:          b.Divider == nil,
		BlockTypeEmbed:            b.Embed == nil,
		BlockTypeEquation:         b.Equation == nil,
		BlockTypeFile:             b.File == nil,
		BlockTypeHeading1:         b.Heading1 == nil,
		BlockTypeHeading2:         b.Heading2 == nil,
		BlockTypeHeading3:         b.Heading3 == nil,
		BlockTypeImage:            b.Image == nil,
		BlockTypeLinkPreview:      b.LinkPreview == nil,
		BlockTypeLinkToPage:       b.LinkToPage == nil,
		BlockTypeNumberedListItem: b.NumberedListItem == nil,
		BlockTypeParagraph:        b.Paragraph == nil,
		BlockTypePdf:              b.Pdf == nil,
		BlockTypeQuote:            b.Quote == nil,
		BlockTypeSyncedBlock:      b.SyncedBlock == nil,
		BlockTypeTable:            b.Table == nil,
		BlockTypeTableOfContents:  b.TableOfContents == nil,
		BlockTypeTableRow:         b.TableRow == nil,
		BlockTypeTemplate:         b.Template == nil,
		BlockTypeToDo:             b.ToDo == nil,
		BlockTypeToggle:           b.Toggle == nil,
		BlockTypeUnsupported:      b.Unsupported == nil,
		BlockTypeVideo:            b.Video == nil,
	} {
		if isNil {
			continue
		}

		if seenType {
			return fmt.Errorf("block has more than one type: %q and %q",
				b.Type, tp)
		}

		switch b.Type {
		case "", tp:
			b.Type = tp
			seenType = true
		default:
			return fmt.Errorf("block should be of type %q but was set to %q",
				tp, b.Type)
		}
	}

	if !seenType {
		return errNoType
	}

	return nil
}
