// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheusvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Bool = BoolAttributesSetterValidator{}

func BoolAttributesSetter(expressions ...path.Expression) validator.Bool {
	return BoolAttributesSetterValidator{
		PathExpressions: expressions,
	}
}

type BoolAttributesSetterValidator struct {
	PathExpressions path.Expressions
}

func (c BoolAttributesSetterValidator) Description(context.Context) string {
	return "Verify the attributes are non-null when Bool is true and null when Bool is false"
}

func (c BoolAttributesSetterValidator) MarkdownDescription(ctx context.Context) string {
	return c.Description(ctx)
}

func (c BoolAttributesSetterValidator) ValidateBool(
	ctx context.Context,
	request validator.BoolRequest,
	response *validator.BoolResponse,
) {
	expressions := request.PathExpression.MergeExpressions(c.PathExpressions...)

	// returns false if unknown, null or false
	// Note this assumes the attribute values are not computed and are always optional
	boolVal := request.ConfigValue.ValueBool()

	for _, expression := range expressions {
		// finding all paths that match the expression
		// e.g path.MatchRoot("list_of_objects").AtListIndex(path.MatchAll()).AtName("id")
		matchedPaths, diags := request.Config.PathMatches(ctx, expression)

		response.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(request.Path) {
				continue
			}

			var mpVal attr.Value

			diags := request.Config.GetAttribute(ctx, mp, &mpVal)
			response.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attributes have a known value
			if mpVal.IsUnknown() {
				return
			}

			if boolVal && mpVal.IsNull() {
				response.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					request.Path,
					fmt.Sprintf("Attribute %q must be specified when %q is true", mp, request.Path),
				))
			}

			if !boolVal && !mpVal.IsNull() {
				response.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					request.Path,
					fmt.Sprintf("Attribute %q cannot be specified when %q is not provided or false",
						mp, request.Path),
				))
			}
		}
	}
}
