package notion

import "encoding/json"

type relationConfiguration struct {
	DatabaseId         UUID                      `json:"database_id"`
	SyncedPropertyId   string                    `json:"synced_property_id,omitempty"`
	SyncedPropertyName string                    `json:"synced_property_name,omitempty"`
	Type               RelationConfigurationType `json:"type,omitempty"`
}

func (r RelationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(relationConfiguration(r))
}
