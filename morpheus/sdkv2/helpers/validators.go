// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package helpers

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringLengthAtMost returns a ValidateDiagFunc that rejects a string whose
// UTF-8 byte length exceeds maxLength. Length is measured in bytes (Go's
// len(string)). It mirrors the plugin framework's stringvalidator.LengthAtMost
// so the two validator families emit consistent diagnostics. A negative
// maxLength is a provider bug and yields an Invalid Validator Usage diagnostic.
func StringLengthAtMost(maxLength int) schema.SchemaValidateDiagFunc {
	return func(i any, _ cty.Path) diag.Diagnostics {
		if maxLength < 0 {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid Validator Usage",
				Detail: fmt.Sprintf(
					"When validating the schema, an implementation issue was found. "+
						"This is always an issue with the provider and should be reported "+
						"to the provider developers.\n\nAn invalid usage of the %q validator "+
						"was found: maxLength cannot be less than zero - maxLength: %d",
					"StringLengthAtMost", maxLength,
				),
			}}
		}

		v, ok := i.(string)
		if !ok {
			return diag.FromErr(TypeAssertFailError("string value", i))
		}

		if l := len(v); l > maxLength {
			return diag.Diagnostics{{
				Severity: diag.Error,
				Summary:  "Invalid Attribute Value Length",
				Detail:   fmt.Sprintf("string length must be at most %d, got: %d", maxLength, l),
			}}
		}

		return nil
	}
}
