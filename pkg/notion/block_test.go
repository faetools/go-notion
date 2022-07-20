package notion_test

import (
	"errors"
	"fmt"
	"net/url"
	"unicode/utf8"

	. "github.com/faetools/go-notion/pkg/notion"
	"github.com/google/uuid"
)

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

var allNumberConfigFormats = []NumberConfigFormat{
	NumberConfigFormatBaht,
	NumberConfigFormatCanadianDollar,
	NumberConfigFormatChileanPeso,
	NumberConfigFormatColombianPeso,
	NumberConfigFormatDanishKrone,
	NumberConfigFormatDirham,
	NumberConfigFormatDollar,
	NumberConfigFormatEuro,
	NumberConfigFormatForint,
	NumberConfigFormatFranc,
	NumberConfigFormatHongKongDollar,
	NumberConfigFormatKoruna,
	NumberConfigFormatKrona,
	NumberConfigFormatLeu,
	NumberConfigFormatLira,
	NumberConfigFormatMexicanPeso,
	NumberConfigFormatNewTaiwanDollar,
	NumberConfigFormatNewZealandDollar,
	NumberConfigFormatNorwegianKrone,
	NumberConfigFormatNumber,
	NumberConfigFormatNumberWithCommas,
	NumberConfigFormatPercent,
	NumberConfigFormatPhilippinePeso,
	NumberConfigFormatPound,
	NumberConfigFormatRand,
	NumberConfigFormatReal,
	NumberConfigFormatRinggit,
	NumberConfigFormatRiyal,
	NumberConfigFormatRuble,
	NumberConfigFormatRupee,
	NumberConfigFormatRupiah,
	NumberConfigFormatShekel,
	NumberConfigFormatWon,
	NumberConfigFormatYen,
	NumberConfigFormatYuan,
	NumberConfigFormatZloty,
}

var allPropertyTypes = []PropertyType{
	PropertyTypeCheckbox,
	PropertyTypeCreatedBy,
	PropertyTypeCreatedTime,
	PropertyTypeDate,
	PropertyTypeEmail,
	PropertyTypeFiles,
	PropertyTypeFormula,
	PropertyTypeLastEditedBy,
	PropertyTypeLastEditedTime,
	PropertyTypeMultiSelect,
	PropertyTypeNumber,
	PropertyTypePeople,
	PropertyTypePhoneNumber,
	PropertyTypeRelation,
	PropertyTypeRichText,
	PropertyTypeRollup,
	PropertyTypeSelect,
	PropertyTypeTitle,
	PropertyTypeUrl,
}

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
	if _, err := uuid.Parse(string(b.Id)); err != nil {
		return err
	}

	if b.Object != "block" {
		return fmt.Errorf("block object field is %q", b.Object)
	}

	if err := validateUser(b.CreatedBy); err != nil {
		return err
	}

	if err := validateUser(b.LastEditedBy); err != nil {
		return err
	}

	if err := validateParent(b.Parent); err != nil {
		return err
	}

	switch b.Type {
	case BlockTypeAudio:
		return validateFile(b.Audio)
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
		return validateFile(b.File)
	case BlockTypeHeading1:
		return validateParagraph(b.Heading1)
	case BlockTypeHeading2:
		return validateParagraph(b.Heading2)
	case BlockTypeHeading3:
		return validateParagraph(b.Heading3)
	case BlockTypeImage:
		return validateFile(b.Image)
	case BlockTypeLinkPreview:
		return validateURL(b.LinkPreview.Url)
	case BlockTypeLinkToPage:
		return validateLinkToPage(b.LinkToPage)
	case BlockTypeNumberedListItem:
		return validateParagraph(b.NumberedListItem)
	case BlockTypeParagraph:
		return validateParagraph(b.Paragraph)
	case BlockTypePdf:
		return validateFile(b.Pdf)
	case BlockTypeQuote:
		return validateParagraph(b.Quote)
	case BlockTypeSyncedBlock:
		from := b.SyncedBlock.SyncedFrom
		if from == nil {
			return nil
		}

		switch from.Type {
		case SyncedFromTypeBlockId:
			_, err := uuid.Parse(string(*from.BlockId))
			return err
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
		return validateVideo(b.Video)
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
	for _, item := range collection {
		if err := validate(item); err != nil {
			return err
		}
	}

	return nil
}

func validateFile(f *FileWithCaption) error {
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
	if _, err := uuid.Parse(string(u.Id)); err != nil {
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
		_, err := uuid.Parse(string(*l.DatabaseId))
		return err
	case LinkToPageTypePageId:
		_, err := uuid.Parse(string(*l.PageId))
		return err
	default:
		return fmt.Errorf("unknown link to page type %q", l.Type)
	}
}

func validateParent(p Parent) error {
	switch p.Type {
	case ParentTypeBlockId:
		_, err := uuid.Parse(string(*p.BlockId))
		return err
	case ParentTypeDatabaseId:
		_, err := uuid.Parse(string(*p.DatabaseId))
		return err
	case ParentTypePageId:
		_, err := uuid.Parse(string(*p.PageId))
		return err
	case ParentTypeWorkspace:
		if p.Workspace == nil || !*p.Workspace {
			return errors.New("workspace parent does not have flag set")
		}

		return nil
	default:
		return fmt.Errorf("unknown parent type %q", p.Type)
	}
}

func validateVideo(v *Video) error {
	if err := validateRichTexts(v.Caption); err != nil {
		return err
	}

	switch v.Type {
	case VideoTypeExternal:
		return validateURL(v.External.Url)
	case VideoTypeFile:
		return validateURL(v.File.Url)
	default:
		return fmt.Errorf("unknown video type %q", v.Type)
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
		// TODO continue here
	case MentionTypeDate:
	case MentionTypeLinkPreview:
	case MentionTypePage:
	case MentionTypeUser:
	default:
		return fmt.Errorf("unknown mention type %q", m.Type)
	}

	return nil
}
