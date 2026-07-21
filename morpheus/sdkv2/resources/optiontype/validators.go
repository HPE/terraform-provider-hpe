// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

// validateDependentFieldNotSelf is a CustomizeDiff function for option type
// resources that expose both field_name and dependent_field. It rejects a
// configuration where dependent_field equals field_name: a field cannot depend
// on itself, which produces a circular dependsOnCode and an unstable form
// reload. Morpheus stores a submitted dependsOnCode verbatim and provides no
// guard against self-reference, so the provider rejects it at plan time.
//
// The check reads the raw config, so it never fires on a computed read-back of
// dependent_field.
func validateDependentFieldNotSelf(
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

// validateOptionTypeRows is a ValidateDiagFunc for the textarea option type's
// rows attribute. rows is submitted verbatim into the option type's config JSON,
// which Morpheus stores without validating, so a non-numeric value such as
// "invalid" is silently accepted by the API. The provider therefore rejects a
// non-integer (or negative) value at plan time. An empty value is treated as
// unset and left to the Optional/Computed default.
func validateOptionTypeRows(i any, _ cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("rows", i))
	}

	if v == "" {
		return nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return diag.Errorf("rows must be a whole number, got %q", v)
	}
	if n < 0 {
		return diag.Errorf("rows must be a non-negative integer (>= 0), got %d", n)
	}

	return nil
}

// optionTypeDescriptionMaxBytes is the maximum UTF-8 byte length the Morpheus
// API accepts for an option type description: OptionTypesController.save guards
// on description.getBytes('UTF-8').length > 255, and the OptionType domain
// constrains description to maxSize: 255. A longer value is rejected with an
// opaque HTTP 400, so the provider enforces the same limit at plan time.
const optionTypeDescriptionMaxBytes = 255

// validateOptionTypeDescription is a ValidateDiagFunc for the option type
// description attribute. Go's len(string) is the UTF-8 byte length, so this
// mirrors the API guard exactly and surfaces a clear error at plan time.
func validateOptionTypeDescription(i any, _ cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("description", i))
	}

	if len(v) > optionTypeDescriptionMaxBytes {
		return diag.Errorf(
			"description must be %d bytes or fewer (Morpheus API limit); got %d",
			optionTypeDescriptionMaxBytes, len(v),
		)
	}

	return nil
}
