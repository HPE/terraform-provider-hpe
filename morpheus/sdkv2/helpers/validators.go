// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateDependentFieldNotSelf is a CustomizeDiff function for option type
// resources that expose both field_name and dependent_field. It rejects a
// configuration where dependent_field equals field_name: a field cannot depend
// on itself, which produces a circular dependsOnCode and an unstable form
// reload. Morpheus stores a submitted dependsOnCode verbatim and provides no
// guard against self-reference, so the provider rejects it at plan time.
//
// The check reads the raw config, so it never fires on a computed read-back of
// dependent_field.
func ValidateDependentFieldNotSelf(
	_ context.Context, d *schema.ResourceDiff, _ any,
) error {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}

	fieldNameVal := raw.GetAttr("field_name")
	dependentFieldVal := raw.GetAttr("dependent_field")
	if fieldNameVal.IsNull() || !fieldNameVal.IsKnown() ||
		dependentFieldVal.IsNull() || !dependentFieldVal.IsKnown() {
		return nil
	}

	fieldName := fieldNameVal.AsString()
	dependentField := dependentFieldVal.AsString()
	if dependentField != "" && dependentField == fieldName {
		return fmt.Errorf(
			"dependent_field must not equal field_name (%q): a field cannot depend "+
				"on itself, which creates a circular dependsOnCode; point "+
				"dependent_field at a different field",
			fieldName,
		)
	}

	return nil
}
