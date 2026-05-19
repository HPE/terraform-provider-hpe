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

var _ validator.Int64 = RequiresNonZeroInt64AtValidator{}

// RequiresNonZeroInt64AtValidator validates that a referenced root-level int64
// attribute is present and has a non-zero value when this attribute is set.
type RequiresNonZeroInt64AtValidator struct {
	AttributeName string
}

func (v RequiresNonZeroInt64AtValidator) Description(context.Context) string {
	return fmt.Sprintf("requires %q to be set to a non-zero value", v.AttributeName)
}

func (v RequiresNonZeroInt64AtValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("requires `%s` to be set to a non-zero value", v.AttributeName)
}

func (v RequiresNonZeroInt64AtValidator) ValidateInt64(
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
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Missing required attribute",
				fmt.Sprintf(
					"Attribute %q must be set to a non-zero value when %q is configured.",
					v.AttributeName,
					request.Path,
				),
			),
		)

		return
	}

	if refValue.ValueInt64() == 0 {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Invalid attribute value",
				fmt.Sprintf(
					"Attribute %q must be set to a non-zero value when %q is configured, got: 0.",
					v.AttributeName,
					request.Path,
				),
			),
		)
	}
}

// RequiresNonZeroInt64At returns a validator that ensures the specified
// root-level int64 attribute is present and has a non-zero value when this
// int64 attribute is configured.
func RequiresNonZeroInt64At(attributeName string) validator.Int64 {
	return RequiresNonZeroInt64AtValidator{
		AttributeName: attributeName,
	}
}
