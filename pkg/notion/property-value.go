package notion

import (
	"encoding/json"
	"time"
)

type (
	propertyValue PropertyValue

	date struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type"`
		Date *Date  `json:"date"`
	}

	number struct {
		ID     string   `json:"id,omitempty"`
		Type   string   `json:"type"`
		Number *float64 `json:"number"`
	}

	propertyValueSelect struct {
		ID     string       `json:"id"`
		Type   PropertyType `json:"type"`
		Select *SelectValue `json:"select"`
	}

	propertyValueStatus struct {
		ID     string       `json:"id"`
		Type   PropertyType `json:"type"`
		Status *SelectValue `json:"status"`
	}

	propertyValueURL struct {
		ID   string       `json:"id"`
		Type PropertyType `json:"type"`
		URL  *string      `json:"url"`
	}

	propertyValuePhone struct {
		ID          string       `json:"id"`
		Type        PropertyType `json:"type"`
		PhoneNumber *string      `json:"phone_number"`
	}

	propertyValueEmail struct {
		ID    string       `json:"id"`
		Type  PropertyType `json:"type"`
		Email *string      `json:"email"`
	}
)

// MarshalJSON fulfils json.Marshaler.
func (v PropertyValue) MarshalJSON() ([]byte, error) {
	switch v.Type {
	case PropertyTypeDate:
		return json.Marshal(date{
			ID:   v.Id,
			Type: string(v.Type),
			Date: v.Date,
		})
	case PropertyTypeNumber:
		return json.Marshal(number{
			ID:     v.Id,
			Type:   string(v.Type),
			Number: v.Number,
		})
	case PropertyTypeSelect:
		return json.Marshal(propertyValueSelect{
			ID:     v.Id,
			Type:   v.Type,
			Select: v.Select,
		})
	case PropertyTypeStatus:
		return json.Marshal(propertyValueStatus{
			ID:     v.Id,
			Type:   v.Type,
			Status: v.Status,
		})
	case PropertyTypeUrl:
		return json.Marshal(propertyValueURL{
			ID:   v.Id,
			Type: v.Type,
			URL:  v.Url,
		})
	case PropertyTypePhoneNumber:
		return json.Marshal(propertyValuePhone{
			ID:          v.Id,
			Type:        v.Type,
			PhoneNumber: v.PhoneNumber,
		})
	case PropertyTypeEmail:
		return json.Marshal(propertyValueEmail{
			ID:    v.Id,
			Type:  v.Type,
			Email: v.Email,
		})
	default:
		return json.Marshal(propertyValue(v))
	}
}

// GetCheckbox returns the checkbox value.
func (v *PropertyValue) GetCheckbox() bool {
	return v != nil && v.Checkbox != nil && *v.Checkbox
}

// GetCreatedBy returns the user that created the object.
func (v *PropertyValue) GetCreatedBy() User {
	if v == nil || v.CreatedBy == nil {
		return User{}
	}

	return *v.CreatedBy
}

// GetCreatedTime returns the time the object was created.
func (v *PropertyValue) GetCreatedTime() time.Time {
	if v == nil || v.CreatedTime == nil {
		return time.Time{}
	}

	return *v.CreatedTime
}

// GetDate returns the date value.
func (v *PropertyValue) GetDate() Date {
	if v == nil || v.Date == nil {
		return Date{}
	}

	return *v.Date
}

// GetEmail returns the email value.
func (v *PropertyValue) GetEmail() string {
	if v == nil || v.Email == nil {
		return ""
	}

	return *v.Email
}

// GetFiles returns the files value.
func (v *PropertyValue) GetFiles() Files {
	if v == nil || v.Files == nil {
		return nil
	}

	return *v.Files
}

// GetFormula returns the formula.
func (v *PropertyValue) GetFormula() Formula {
	if v == nil || v.Formula == nil {
		return Formula{}
	}

	return *v.Formula
}

// GetLastEditedBy returns the user that last edited the object.
func (v *PropertyValue) GetLastEditedBy() User {
	if v == nil || v.LastEditedBy == nil {
		return User{}
	}

	return *v.LastEditedBy
}

// GetMultiSelect returns the multiselect value.
func (v *PropertyValue) GetMultiSelect() SelectValues {
	if v == nil || v.MultiSelect == nil {
		return SelectValues{}
	}

	return *v.MultiSelect
}

// GetNumber returns the number value.
func (v *PropertyValue) GetNumber() float64 {
	if v == nil || v.Number == nil {
		return 0
	}

	return *v.Number
}

// GetPeople returns the people value.
func (v *PropertyValue) GetPeople() Users {
	if v == nil || v.People == nil {
		return nil
	}

	return *v.People
}

// GetPhoneNumber returns the phone number.
func (v *PropertyValue) GetPhoneNumber() string {
	if v == nil || v.PhoneNumber == nil {
		return ""
	}

	return *v.PhoneNumber
}

// GetRelation returns the relation value.
func (v *PropertyValue) GetRelation() References {
	if v == nil || v.Relation == nil {
		return nil
	}

	return *v.Relation
}

// GetRichText returns the rich text value.
func (v *PropertyValue) GetRichText() RichTexts {
	if v == nil || v.RichText == nil {
		return nil
	}

	return *v.RichText
}

// GetRollup returns the rollup value.
func (v *PropertyValue) GetRollup() Rollup {
	if v == nil || v.Rollup == nil {
		return Rollup{}
	}

	return *v.Rollup
}

// GetSelect returns the value that was selected.
func (v *PropertyValue) GetSelect() SelectValue {
	if v == nil || v.Select == nil {
		return SelectValue{}
	}

	return *v.Select
}

// GetStatus returns the status that was selected.
func (v *PropertyValue) GetStatus() SelectValue {
	if v == nil || v.Status == nil {
		return SelectValue{}
	}

	return *v.Status
}

// GetTitle returns the title value.
func (v *PropertyValue) GetTitle() RichTexts {
	if v == nil || v.Title == nil {
		return nil
	}

	return *v.Title
}

// GetURL returns the URL of the object.
func (v *PropertyValue) GetURL() string {
	if v == nil || v.Url == nil {
		return ""
	}

	return *v.Url
}
