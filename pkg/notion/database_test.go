package notion_test

import (
	"fmt"

	. "github.com/faetools/go-notion/pkg/notion"
	"github.com/google/uuid"
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

func validateDatabase(db *Database) error {
	if err := validateFile(db.Cover); err != nil {
		return err
	}

	if err := validateUser(db.CreatedBy); err != nil {
		return err
	}

	if err := validateRichTexts(db.Description); err != nil {
		return err
	}

	if db.Icon != nil {
		if err := validateIcon(*db.Icon); err != nil {
			return err
		}
	}

	if _, err := uuid.Parse(string(db.Id)); err != nil {
		return err
	}

	if err := validateUser(db.LastEditedBy); err != nil {
		return err
	}

	if db.Object != "database" {
		return fmt.Errorf("object field of database is %q", db.Object)
	}

	if err := validatePropertyMetaMap(db.Properties); err != nil {
		return err
	}

	if err := validateRichTexts(db.Title); err != nil {
		return err
	}

	return validateURL(db.Url)
}

func validateFile(f *File) error {
	if f == nil {
		return nil
	}

	switch f.Type {
	case FileTypeExternal:
		return validateURL(f.External.Url)
	case FileTypeFile:
		return validateURL(f.File.Url)
	default:
		return fmt.Errorf("unknown file type %q", f.Type)
	}
}

func validatePropertyMetaMap(m PropertyMetaMap) error {
	for k, prop := range m {
		if prop.Name != k {
			return fmt.Errorf("name is not same as key; name %q and key %q", prop.Name, k)
		}

		if err := validatePropertyMeta(prop); err != nil {
			return err
		}
	}

	return nil
}

func errIfNil[T any](propName string, meta *T) error {
	if meta == nil {
		return fmt.Errorf("%s is empty", propName)
	}

	return nil
}

func validatePropertyMeta(p PropertyMeta) error {
	switch p.Type {
	case PropertyTypeCheckbox:
		return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeCreatedBy:
	// 	return errIfNil("created_by", p.)
	// case PropertyTypeCreatedTime:
	// 	return errIfNil("checkbox", p.Crea)
	case PropertyTypeDate:
		return errIfNil("date", p.Date)
	// case PropertyTypeEmail:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeFiles:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeFormula:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeLastEditedBy:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeLastEditedTime:
	// 	return errIfNil("checkbox", p.Checkbox)
	case PropertyTypeMultiSelect:
		return errIfNil("multi_select", p.MultiSelect)
	case PropertyTypeNumber:
		for _, format := range allNumberConfigFormats {
			if format == p.Number.Format {
				return nil
			}
		}

		return fmt.Errorf("unkown number format %q", p.Number.Format)
	// case PropertyTypePeople:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypePhoneNumber:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeRelation:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeRichText:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeRollup:
	// 	return errIfNil("checkbox", p.Checkbox)
	// case PropertyTypeSelect:
	// 	return errIfNil("checkbox", p.Checkbox)
	case PropertyTypeTitle:
		return errIfNil("title", p.Title)
	case PropertyTypeUrl:
		// return errIfNil("url", p.)
		return nil
	default:
		return nil // TODO delete
		return fmt.Errorf("unknown property type %q", p.Type)
	}

	return nil
}
