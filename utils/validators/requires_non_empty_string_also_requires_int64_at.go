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

var _ validator.String = RequiresNonEmptyStringAlsoRequiresInt64AtValidator{}

// RequiresNonEmptyStringAlsoRequiresInt64AtValidator validates that a sibling
// int64 attribute is present when this string attribute is set to a non-empty
// value.
type RequiresNonEmptyStringAlsoRequiresInt64AtValidator struct {
	SiblingName string
}

func (v RequiresNonEmptyStringAlsoRequiresInt64AtValidator) Description(context.Context) string {
	return fmt.Sprintf("requires %q to be set when this attribute is non-empty", v.SiblingName)
}

func (v RequiresNonEmptyStringAlsoRequiresInt64AtValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("requires `%s` to be set when this attribute is non-empty", v.SiblingName)
}

func (v RequiresNonEmptyStringAlsoRequiresInt64AtValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueString() == "" {
		return
	}

	siblingPath := request.Path.ParentPath().AtName(v.SiblingName)

	var siblingValue types.Int64
	diags := request.Config.GetAttribute(ctx, siblingPath, &siblingValue)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	if siblingValue.IsNull() || siblingValue.IsUnknown() {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Missing required sibling attribute",
				fmt.Sprintf(
					"Attribute %q must be set when %q is %q.",
					path.Root(v.SiblingName),
					request.Path,
					request.ConfigValue.ValueString(),
				),
			),
		)
	}
}

// RequiresNonEmptyStringAlsoRequiresInt64At returns a validator that ensures
// the specified sibling int64 attribute is present when this string attribute
// is set to a non-empty value.
func RequiresNonEmptyStringAlsoRequiresInt64At(siblingName string) validator.String {
	return RequiresNonEmptyStringAlsoRequiresInt64AtValidator{
		SiblingName: siblingName,
	}
}
