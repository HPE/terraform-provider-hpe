package helpers

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringMaxLength returns a ValidateDiagFunc that rejects a string whose UTF-8
// byte length exceeds maxLength. Length is measured in bytes (Go's len(string)),
// matching the Morpheus API's getBytes("UTF-8").length guards; exceeding such a
// limit otherwise yields an opaque HTTP 400, so validating at plan time surfaces
// a clear error up front. The framework attaches the offending attribute path to
// the diagnostic automatically.
func StringMaxLength(maxLength int) schema.SchemaValidateDiagFunc {
	return func(i any, _ cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.FromErr(TypeAssertFailError("string value", i))
		}

		if len(v) > maxLength {
			return diag.Errorf("value must be %d bytes or fewer, got %d", maxLength, len(v))
		}

		return nil
	}
}
