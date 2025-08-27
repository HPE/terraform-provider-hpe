// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheusvalidators

// This validator allows us to verify, at plan time
// that a dynamic attribute looks like an object, eg
// this is ok
//   config = {
//     foo = "bar"
//   }
// This would not pass validation
//   config = "foo"
import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
)

var _ validator.Dynamic = DynamicAttributeObjectValidator{}

type DynamicAttributeObjectValidator struct{}

func (v DynamicAttributeObjectValidator) Description(context.Context) string {
	return "verify that the dynamic attribute can be converted to a valid object/map"
}

func (v DynamicAttributeObjectValidator) MarkdownDescription(context.Context) string {
	return "verify that the dynamic attribute can be converted to a valid object/map"
}

func (v DynamicAttributeObjectValidator) ValidateDynamic(
	ctx context.Context,
	request validator.DynamicRequest,
	response *validator.DynamicResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	configValue := request.ConfigValue.UnderlyingValue()
	configMap, err := convert.ValueToAny(ctx, configValue)
	if err != nil {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Invalid format",
				"attribute must be a valid object/map: "+err.Error(),
			),
		)

		return
	}

	// Check if it can be converted to map[string]any
	_, ok := configMap.(map[string]any)
	if !ok {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Invalid type",
				"attribute must be a valid object/map",
			),
		)
	}
}

// ValidObjectMap returns a validator that ensures the dynamic attribute can be
// converted to a valid object/map using the same logic as the apply-time
// validation
func ValidObjectMap() validator.Dynamic {
	return DynamicAttributeObjectValidator{}
}
