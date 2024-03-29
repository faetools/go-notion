package notion

// GetFormula returns the formula configuration.
func (v PropertyMeta) GetFormula() FormulaConfig {
	if v.Formula == nil {
		return FormulaConfig{}
	}

	return *v.Formula
}

// GetMultiSelect returns the multiselect configuration.
func (v PropertyMeta) GetMultiSelect() SelectValuesWrapper {
	if v.MultiSelect == nil {
		return SelectValuesWrapper{}
	}

	return *v.MultiSelect
}

// GetNumber returns the number configuration.
func (v PropertyMeta) GetNumber() NumberConfig {
	if v.Number == nil {
		return NumberConfig{}
	}

	return *v.Number
}

// GetRelation returns the relation configuration.
func (v PropertyMeta) GetRelation() RelationConfiguration {
	if v.Relation == nil {
		return RelationConfiguration{}
	}

	return *v.Relation
}

// GetRollup returns the rollup configuration.
func (v PropertyMeta) GetRollup() RollupConfig {
	if v.Rollup == nil {
		return RollupConfig{}
	}

	return *v.Rollup
}

// GetSelect returns the select configuration.
func (v PropertyMeta) GetSelect() SelectValuesWrapper {
	if v.Select == nil {
		return SelectValuesWrapper{}
	}

	return *v.Select
}
