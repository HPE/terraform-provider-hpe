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
// so the two validator families emit consistent diagnostics.
func StringLengthAtMost(maxLength int) schema.SchemaValidateDiagFunc {
	return func(i any, _ cty.Path) diag.Diagnostics {
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
