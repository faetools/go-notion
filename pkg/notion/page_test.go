package notion_test

import (
	"errors"
	"fmt"

	. "github.com/faetools/go-notion/pkg/notion"
)

func validatePage(p *Page) error {
	if err := validateUUID(p.Id); err != nil {
		return err
	}

	if err := validateFile(p.Cover); err != nil {
		return err
	}

	if err := validateUser(p.CreatedBy); err != nil {
		return err
	}

	if err := validateUser(p.LastEditedBy); err != nil {
		return err
	}

	if p.Icon != nil {
		if err := validateIcon(*p.Icon); err != nil {
			return err
		}
	}

	if p.Object != "page" {
		return fmt.Errorf("object field of page is %q", p.Object)
	}

	if err := validateParent(p.Parent); err != nil {
		return err
	}

	if err := validatePropertyValueMap(p.Properties); err != nil {
		return err
	}

	return validateURL(p.Url)
}

func validatePropertyValueMap(m PropertyValueMap) error {
	for k, prop := range m {
		if err := validateShortID(k, prop.Id); err != nil {
			return err
		}

		if err := validatePropertyValue(prop); err != nil {
			return fmt.Errorf("property with key %q: %w", k, err)
		}
	}

	return nil
}

func validatePropertyValue(p PropertyValue) error {
	switch p.Type {
	case PropertyTypeTitle:
		return validateRichTexts(*p.Title)
	case PropertyTypeUrl:
		if p.Url != nil {
			return validateURL(*p.Url)
		}
	case PropertyTypeRollup:
		return validateRollup(p.Rollup)
	case PropertyTypeCheckbox:
		if p.Checkbox == nil {
			return errors.New("checkbox value is empty")
		}
	case PropertyTypeCreatedTime:
		if p.CreatedTime == nil {
			return errors.New("created time value is empty")
		}
	case PropertyTypeStatus:
		return validateSelect(p.Status)
	case PropertyTypeRichText:
		return validateRichTexts(*p.RichText)
	case PropertyTypeRelation:
		return validateReferences(p.Relation)
	case PropertyTypeCreatedBy:
		return validateUser(p.CreatedBy)
	case PropertyTypeLastEditedBy:
		return validateUser(p.LastEditedBy)
	case PropertyTypePeople:
		for _, p := range *p.People {
			if err := validateUser(&p); err != nil {
				return err
			}
		}
	case PropertyTypeFiles:
		for _, f := range *p.Files {
			if err := validateFile(&f); err != nil {
				return err
			}
		}
	case PropertyTypeFormula:
		return validateFormula(p.Formula)
	case PropertyTypeSelect:
		return validateSelect(p.Select)
	case PropertyTypeMultiSelect:
		for _, sel := range *p.MultiSelect {
			if err := validateSelect(&sel); err != nil {
				return err
			}
		}
	case PropertyTypeNumber, PropertyTypePhoneNumber, PropertyTypeDate, PropertyTypeEmail:
		// it's fine for these values to be empty in case the user has not filled the cell
		return nil
	default:
		return fmt.Errorf("unknown property type %q", p.Type)
	}

	return nil
}

func validateRollup(r *Rollup) error {
	switch r.Type {
	case RollupTypeString:
		if r.String == nil {
			return errors.New("rollup string is empty")
		}
	case RollupTypeArray:
		return validateRollupArray(r.Array)
	case RollupTypeNumber:
		if r.Number == nil {
			return errors.New("rollup number is empty")
		}
	default:
		return fmt.Errorf("unknown rollup type %q", r.Type)
	}

	return nil
}

func validateSelect(s *SelectValue) error {
	if s == nil {
		return nil
	}

	// ID is either UUID or short string
	if err := validateUUID(UUID(s.Id)); err != nil {
		if err := validateShortID("select value", s.Id); err != nil {
			return err
		}
	}

	if s.Name == "" {
		return fmt.Errorf("select value %q has no name", s.Id)
	}

	return validateColor(s.Color)
}

func validateReferences(refs *References) error {
	for _, ref := range *refs {
		if err := validateUUID(ref.Id); err != nil {
			return err
		}
	}

	return nil
}

func validateFormula(f *Formula) error {
	switch f.Type {
	case FormulaTypeString:
		if f.String == nil {
			return errors.New("formula without string")
		}
	case FormulaTypeBoolean:
		if f.Boolean == nil {
			return errors.New("formula without boolean")
		}
	case FormulaTypeDate, FormulaTypeNumber:
		// can be nil
		return nil
	default:
		return fmt.Errorf("unknown formula type %q", f.Type)
	}
	return nil
}

func validateRollupArray(array *RollupArray) error {
	if array == nil {
		return errors.New("rollup array is empty")
	}

	if len(*array) == 0 {
		return nil
	}

	tp := (*array)[0].Type

	for _, v := range *array {
		if v.Type != tp {
			return fmt.Errorf("conflicting types in rollup array: %q vs. %q", tp, v.Type)
		}

		switch v.Type {
		case RollupArrayItemTypeDate,
			RollupArrayItemTypeNumber,
			RollupArrayItemTypeString: // Ok
		case RollupArrayItemTypeTitle:
			if v.Title == nil {
				continue
			}

			if err := validateRichTexts(*v.Title); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unkown rollup array type %q", v.Type)

		}
	}

	return nil
}
