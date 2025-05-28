package morpheusvalidators

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = JSONValidator{}

type JSONValidator struct{}

func (c JSONValidator) Description(context.Context) string {
	return "verify that the attribute is a valid JSON blob"
}

func (c JSONValidator) MarkdownDescription(context.Context) string {
	return "verify that the attribute is a valid JSON blob"
}

func (c JSONValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if err := json.Unmarshal([]byte(value), &map[string]any{}); err != nil {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"not valid JSON",
				"attribute must contain a JSON value: "+err.Error(),
			),
		)
	}
}
