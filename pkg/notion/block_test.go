package notion_test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"unicode/utf8"

	. "github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	t.Parallel()
	t.Skip() // TODO: once a new version is released, test if we can get a mentioned block

	ctx := context.Background()

	// https://www.notion.so/fae-tools/B-cher-zum-Deutschlernen-e3af0440bd0b4be3833bae72db427d90?pvs=4#abc1910b59ec49b386e4988d89e33f48

	c, err := NewDefaultClient(os.Getenv("NOTION_TOKEN"))
	assert.NoError(t, err)

	rsp, err := c.GetBlock(ctx, "abc1910b59ec49b386e4988d89e33f48")
	assert.NoError(t, err)

	// NOTE: Skipping right over the mentioned block!
	assert.Equal(t, `{"object":"block","id":"abc1910b-59ec-49b3-86e4-988d89e33f48","parent":{"type":"page_id","page_id":"e3af0440-bd0b-4be3-833b-ae72db427d90"},"created_time":"2023-10-24T19:42:00.000Z","last_edited_time":"2023-10-24T20:38:00.000Z","created_by":{"object":"user","id":"af171d5d-c36f-45bc-a0a3-6086c0dafa45"},"last_edited_by":{"object":"user","id":"af171d5d-c36f-45bc-a0a3-6086c0dafa45"},"has_children":false,"archived":false,"type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"Springe zum Abschnitt: ","link":null},"annotations":{"bold":true,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Springe zum Abschnitt: ","href":null},{"type":"text","text":{"content":" TODO","link":null},"annotations":{"bold":true,"italic":true,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":" TODO","href":null}],"color":"default"},"request_id":"0b532314-b1c6-46c6-bdb1-6a4e759c26db"}`, string(rsp.Body))
}

var (
	allCodeLanguages = []CodeLanguage{
		CodeLanguageAbap,
		CodeLanguageArduino,
		CodeLanguageBash,
		CodeLanguageBasic,
		CodeLanguageC,
		CodeLanguageC1,
		CodeLanguageC2,
		CodeLanguageClojure,
		CodeLanguageCoffeescript,
		CodeLanguageCss,
		CodeLanguageDart,
		CodeLanguageDiff,
		CodeLanguageDocker,
		CodeLanguageElixir,
		CodeLanguageElm,
		CodeLanguageErlang,
		CodeLanguageF,
		CodeLanguageFlow,
		CodeLanguageFortran,
		CodeLanguageGherkin,
		CodeLanguageGlsl,
		CodeLanguageGo,
		CodeLanguageGraphql,
		CodeLanguageGroovy,
		CodeLanguageHaskell,
		CodeLanguageHtml,
		CodeLanguageJava,
		CodeLanguageJavaccc,
		CodeLanguageJavascript,
		CodeLanguageJson,
		CodeLanguageJulia,
		CodeLanguageKotlin,
		CodeLanguageLatex,
		CodeLanguageLess,
		CodeLanguageLisp,
		CodeLanguageLivescript,
		CodeLanguageLua,
		CodeLanguageMakefile,
		CodeLanguageMarkdown,
		CodeLanguageMarkup,
		CodeLanguageMatlab,
		CodeLanguageMermaid,
		CodeLanguageNix,
		CodeLanguageObjectiveC,
		CodeLanguageOcaml,
		CodeLanguagePascal,
		CodeLanguagePerl,
		CodeLanguagePhp,
		CodeLanguagePlainText,
		CodeLanguagePowershell,
		CodeLanguageProlog,
		CodeLanguageProtobuf,
		CodeLanguagePython,
		CodeLanguageR,
		CodeLanguageReason,
		CodeLanguageRuby,
		CodeLanguageRust,
		CodeLanguageSass,
		CodeLanguageScala,
		CodeLanguageScheme,
		CodeLanguageScss,
		CodeLanguageShell,
		CodeLanguageSql,
		CodeLanguageSwift,
		CodeLanguageTypescript,
		CodeLanguageVbNet,
		CodeLanguageVerilog,
		CodeLanguageVhdl,
		CodeLanguageVisualBasic,
		CodeLanguageWebassembly,
		CodeLanguageXml,
		CodeLanguageYaml,
	}

	allColors = []Color{
		ColorBlue,
		ColorBlueBackground,
		ColorBrown,
		ColorBrownBackground,
		ColorDefault,
		ColorGray,
		ColorGrayBackground,
		ColorGreen,
		ColorGreenBackground,
		ColorOrange,
		ColorOrangeBackground,
		ColorPink,
		ColorPinkBackground,
		ColorPurple,
		ColorPurpleBackground,
		ColorRed,
		ColorRedBackground,
		ColorYellow,
		ColorYellowBackground,
	}
)

func validateBlocks(bs Blocks) error {
	return firstError(bs, validateBlock)
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("missing URL")
	}

	_, err := url.Parse(rawURL)

	return err
}

func validateBlock(b Block) error {
	if err := validateUUID(b.Id); err != nil {
		return err
	}

	if b.Object != "block" {
		return fmt.Errorf("block object field is %q", b.Object)
	}

	tp := b.Type
	b.Type = ""

	if err := b.Validate(); err != nil {
		return err
	}

	if b.Type != tp {
		return fmt.Errorf("could not recover type %q", tp)
	}

	if err := validateUser(b.CreatedBy); err != nil {
		return err
	}

	if err := validateUser(b.LastEditedBy); err != nil {
		return err
	}

	if err := validateParent(&b.Parent); err != nil {
		return err
	}

	switch b.Type {
	case BlockTypeAudio:
		return validateFileWithCaption(b.Audio)
	case BlockTypeBookmark:
		if err := validateURL(b.Bookmark.Url); err != nil {
			return err
		}

		return validateRichTexts(b.Bookmark.Caption)
	case BlockTypeBreadcrumb:
		if b.Breadcrumb == nil {
			return errors.New("no breadcrumb object")
		}
	case BlockTypeBulletedListItem:
		return validateParagraph(b.BulletedListItem)
	case BlockTypeCallout:
		return validateCallout(b.Callout)
	case BlockTypeChildDatabase, BlockTypeChildPage:
		// unfortunately, we don't get the ID
		return nil
	case BlockTypeCode:
		return validateCode(b.Code)
	case BlockTypeColumn:
		if b.Column == nil {
			return errors.New("no column object")
		}
	case BlockTypeColumnList:
		if b.ColumnList == nil {
			return errors.New("no column list object")
		}
	case BlockTypeDivider:
		if b.Divider == nil {
			return errors.New("no divider object")
		}
	case BlockTypeEmbed:
		if err := validateURL(b.Embed.Url); err != nil {
			return err
		}

		return validateRichTexts(b.Embed.Caption)
	case BlockTypeEquation:
		return validateEquation(b.Equation)
	case BlockTypeFile:
		return validateFileWithCaption(b.File)
	case BlockTypeHeading1:
		return validateHeading(b.Heading1)
	case BlockTypeHeading2:
		return validateHeading(b.Heading2)
	case BlockTypeHeading3:
		return validateHeading(b.Heading3)
	case BlockTypeImage:
		return validateFileWithCaption(b.Image)
	case BlockTypeLinkPreview:
		return validateURL(b.LinkPreview.Url)
	case BlockTypeLinkToPage:
		return validateLinkToPage(b.LinkToPage)
	case BlockTypeNumberedListItem:
		return validateParagraph(b.NumberedListItem)
	case BlockTypeParagraph:
		return validateParagraph(b.Paragraph)
	case BlockTypePdf:
		return validateFileWithCaption(b.Pdf)
	case BlockTypeQuote:
		return validateParagraph(b.Quote)
	case BlockTypeSyncedBlock:
		from := b.SyncedBlock.SyncedFrom
		if from == nil {
			return nil
		}

		switch from.Type {
		case SyncedFromTypeBlockId:
			return validateUUID(*from.BlockId)
		default:
			return fmt.Errorf("synced block was synced from %q", from.Type)
		}
	case BlockTypeTable:
		if b.Table.TableWidth < 1 {
			return fmt.Errorf("table has width %d", b.Table.TableWidth)
		}
	case BlockTypeTableOfContents:
		return validateColor(b.TableOfContents.Color)
	case BlockTypeTableRow:
		return firstError(b.TableRow.Cells, validateRichTexts)
	case BlockTypeTemplate:
		return validateRichTexts(b.Template.RichText)
	case BlockTypeToDo:
		if err := validateRichTexts(b.ToDo.RichText); err != nil {
			return err
		}

		return validateColor(b.ToDo.Color)
	case BlockTypeToggle:
		return validateParagraph(b.Toggle)
	case BlockTypeUnsupported:
		if b.Unsupported == nil {
			return errors.New("no unsupported object")
		}
	case BlockTypeVideo:
		return validateFileWithCaption(b.Video)
	default:
		return fmt.Errorf("unknown block type %q", b.Type)
	}

	return nil
}

func validateCode(c *Code) error {
	if c.Caption != nil {
		if err := validateRichTexts(*c.Caption); err != nil {
			return err
		}
	}

	if err := validateRichTexts(c.RichText); err != nil {
		return err
	}

	for _, lang := range allCodeLanguages {
		if lang == c.Language {
			return nil
		}
	}

	return fmt.Errorf("unknown code language %q", c.Language)
}

func validateRichTexts(ts RichTexts) error {
	return firstError(ts, validateRichText)
}

func validateRichText(t RichText) error {
	if err := validateAnnotations(t.Annotations); err != nil {
		return err
	}

	if t.Href != nil {
		if err := validateURL(*t.Href); err != nil {
			return err
		}
	}

	if t.PlainText == "" {
		return errors.New("no text in rich text")
	}

	switch t.Type {
	case RichTextTypeEquation:
		return validateEquation(t.Equation)
	case RichTextTypeMention:
		return validateMention(t.Mention)
	case RichTextTypeText:
		if t.Text.Content == "" {
			return errors.New("no content in text")
		}

		if t.Text.Link != nil {
			return validateURL(t.Text.Link.Url)
		}

		return nil
	default:
		return fmt.Errorf("unkown rich text type %q", t.Type)
	}
}

func validateColor(c Color) error {
	for _, validColor := range allColors {
		if c == validColor {
			return nil
		}
	}

	return fmt.Errorf("unknown color %q", c)
}

func validateAnnotations(a Annotations) error {
	return validateColor(a.Color)
}

func firstError[T any](collection []T, validate func(T) error) error {
	for i, item := range collection {
		if err := validate(item); err != nil {
			return fmt.Errorf("@%d: %w", i, err)
		}
	}

	return nil
}

func validateFileWithCaption(f *FileWithCaption) error {
	if f.Caption != nil {
		if err := validateRichTexts(*f.Caption); err != nil {
			return err
		}
	}

	switch f.Type {
	case FileWithCaptionTypeExternal:
		return validateURL(f.External.Url)
	case FileWithCaptionTypeFile:
		return validateURL(f.File.Url)
	default:
		return fmt.Errorf("unknown file type %q", f.Type)
	}
}

func validateParagraph(p *Paragraph) error {
	if err := validateColor(p.Color); err != nil {
		return err
	}

	return validateRichTexts(p.RichText)
}

func validateHeading(p *Heading) error {
	if err := validateColor(p.Color); err != nil {
		return err
	}

	return validateRichTexts(p.RichText)
}

func validateCallout(c *Callout) error {
	if err := validateColor(c.Color); err != nil {
		return err
	}

	if err := validateIcon(c.Icon); err != nil {
		return err
	}

	return validateRichTexts(c.RichText)
}

func validateIcon(c Icon) error {
	switch c.Type {
	case IconTypeEmoji:
		if c.Emoji == nil || *c.Emoji == "" {
			return errors.New("no emoji provided")
		}

		if count := utf8.RuneCountInString(*c.Emoji); count > 2 {
			return fmt.Errorf("emoji %q has %d runes", *c.Emoji, count)
		}

		return nil

	case IconTypeExternal:
		return validateURL(c.External.Url)
	case IconTypeFile:
		return validateURL(c.File.Url)
	default:
		return fmt.Errorf("unknown icon type %q", c.Type)
	}
}

func validateUser(u *User) error {
	if err := validateUUID(u.Id); err != nil {
		return err
	}

	if u.AvatarUrl != nil {
		if err := validateURL(*u.AvatarUrl); err != nil {
			return err
		}
	}

	if u.Object != "user" {
		return fmt.Errorf("user has %q as object field", u.Object)
	}

	if u.Type != nil {
		switch *u.Type {
		case UserTypeBot:
			if u.Bot == nil {
				return errors.New("no bot object")
			}
		case UserTypePerson:
			if u.Person == nil {
				return errors.New("no person object")
			}

			if u.Person.Email == "" {
				return errors.New("no email in person object")
			}
		default:
			return fmt.Errorf("unknown user type %q", *u.Type)
		}
	}

	return nil
}

func validateLinkToPage(l *LinkToPage) error {
	switch l.Type {
	case LinkToPageTypeDatabaseId:
		return validateUUID(*l.DatabaseId)
	case LinkToPageTypePageId:
		return validateUUID(*l.PageId)
	default:
		return fmt.Errorf("unknown link to page type %q", l.Type)
	}
}

func validateParent(p *Parent) error {
	if p == nil {
		return nil
	}

	switch p.Type {
	case ParentTypeBlockId:
		return validateUUID(*p.BlockId)
	case ParentTypeDatabaseId:
		return validateUUID(*p.DatabaseId)
	case ParentTypePageId:
		return validateUUID(*p.PageId)
	case ParentTypeWorkspace:
		if p.Workspace == nil || !*p.Workspace {
			return errors.New("workspace parent does not have flag set")
		}

		return nil
	default:
		return fmt.Errorf("unknown parent type %q", p.Type)
	}
}

func validateEquation(eq *Equation) error {
	if eq.Expression == "" {
		return errors.New("no expression in equation")
	}

	return nil
}

func validateMention(m *Mention) error {
	switch m.Type {
	case MentionTypeDatabase:
		return validateUUID(m.Database.Id)
	case MentionTypeDate:
		return nil
	case MentionTypeLinkPreview:
		return validateURL(m.LinkPreview.Url)
	case MentionTypePage:
		return validateUUID(m.Page.Id)
	case MentionTypeUser:
		return validateUser(m.User)
	default:
		return fmt.Errorf("unknown mention type %q", m.Type)
	}
}
