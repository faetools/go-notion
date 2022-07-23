package notion_test

import (
	"errors"
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

	if err := validateUUID(db.Id); err != nil {
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

func errIfNil(tp PropertyType, config *map[string]any) error {
	if config == nil {
		return fmt.Errorf("%s is empty", tp)
	}

	return nil
}

func validateShortID(name, id string) error {
	if 3 < len(id) && len(id) < 9 {
		return nil
	}

	return fmt.Errorf("%s %q has length %d", name, id, len(id))
}

func validateUUID(id UUID) error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return fmt.Errorf("UUID %q: %w", id, err)
	}

	return nil
}

func validatePropertyMeta(p PropertyMeta) error {
	switch p.Type {
	case PropertyTypeCheckbox:
		return errIfNil(p.Type, p.Checkbox)
	case PropertyTypePeople:
		return errIfNil(p.Type, p.People)
	case PropertyTypeStatus:
		return errIfNil(p.Type, p.Status)
	case PropertyTypeCreatedBy:
		return errIfNil(p.Type, p.CreatedBy)
	case PropertyTypeCreatedTime:
		return errIfNil(p.Type, p.CreatedTime)
	case PropertyTypeDate:
		return errIfNil(p.Type, p.Date)
	case PropertyTypeEmail:
		return errIfNil(p.Type, p.Email)
	case PropertyTypeFiles:
		return errIfNil(p.Type, p.Files)
	case PropertyTypeFormula:
		if p.Formula.Expression == "" {
			return errors.New("formula expression is empty")
		}
	case PropertyTypeLastEditedBy:
		return errIfNil(p.Type, p.LastEditedBy)
	case PropertyTypeLastEditedTime:
		return errIfNil(p.Type, p.LastEditedTime)
	case PropertyTypeMultiSelect:
		if p.MultiSelect == nil {
			return fmt.Errorf("%s is empty", p.Type)
		}
	case PropertyTypeNumber:
		for _, format := range allNumberConfigFormats {
			if format == p.Number.Format {
				return nil
			}
		}

		return fmt.Errorf("unkown number format %q", p.Number.Format)
	case PropertyTypePhoneNumber:
		return errIfNil(p.Type, p.PhoneNumber)
	case PropertyTypeRelation:
		if p.Relation.SyncedPropertyName == "" {
			return fmt.Errorf("SyncedPropertyName is empty in %#v", p)
		}

		if err := validateShortID("synced_property_id", *p.Relation.SyncedPropertyId); err != nil {
			return err
		}

		if err := validateUUID(p.Relation.DatabaseId); err != nil {
			return fmt.Errorf("database_id %q: %w", p.Relation.DatabaseId, err)
		}
	case PropertyTypeRichText:
		return errIfNil(p.Type, p.RichText)
	case PropertyTypeRollup:
		if err := validateShortID("rollup_property_id", p.Rollup.RollupPropertyId); err != nil {
			return err
		}

		if err := validateShortID("relation_property_id", p.Rollup.RelationPropertyId); err != nil {
			return err
		}

		if p.Rollup.Function == "" {
			return fmt.Errorf("function is empty in %#v", p)
		}

		if p.Rollup.RelationPropertyName == "" {
			return fmt.Errorf("RelationPropertyName is empty in %#v", p)
		}

		if p.Rollup.RollupPropertyName == "" {
			return fmt.Errorf("RollupPropertyName is empty in %#v", p)
		}
	case PropertyTypeSelect:
		if p.Select == nil {
			return fmt.Errorf("%s is empty", p.Type)
		}
	case PropertyTypeTitle:
		return errIfNil(p.Type, p.Title)
	case PropertyTypeUrl:
		return errIfNil(p.Type, p.Url)
	default:
		return fmt.Errorf("unknown property type %q", p.Type)
	}

	return nil
}
