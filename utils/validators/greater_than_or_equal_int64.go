// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Int64 = GreaterThanOrEqualValidator{}

// GreaterThanOrEqualValidator validates that this int64 attribute is greater
// than or equal to the value of a referenced root-level int64 attribute.
//
// The comparison is skipped when either this attribute or the referenced
// attribute is null or unknown, so pairing this validator with an
// AlsoRequires validator on the referenced attribute is recommended when the
// referenced attribute must be present.
type GreaterThanOrEqualValidator struct {
	AttributeName string
}

func (v GreaterThanOrEqualValidator) Description(context.Context) string {
	return fmt.Sprintf("value must be greater than or equal to %q", v.AttributeName)
}

func (v GreaterThanOrEqualValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("value must be greater than or equal to `%s`", v.AttributeName)
}

func (v GreaterThanOrEqualValidator) ValidateInt64(
	ctx context.Context,
	request validator.Int64Request,
	response *validator.Int64Response,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	var refValue types.Int64
	diags := request.Config.GetAttribute(ctx, path.Root(v.AttributeName), &refValue)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	if refValue.IsNull() || refValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueInt64() < refValue.ValueInt64() {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Invalid attribute value",
				fmt.Sprintf(
					"Attribute %q must be greater than or equal to %q, got: %d.",
					request.Path,
					v.AttributeName,
					request.ConfigValue.ValueInt64(),
				),
			),
		)
	}
}

// GreaterThanOrEqual returns a validator that ensures this int64 attribute is
// greater than or equal to the specified root-level int64 attribute. The
// comparison is skipped when either attribute is null or unknown.
func GreaterThanOrEqual(attributeName string) validator.Int64 {
	return GreaterThanOrEqualValidator{
		AttributeName: attributeName,
	}
}
