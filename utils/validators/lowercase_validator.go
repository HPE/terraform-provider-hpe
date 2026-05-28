// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = LowercaseValidator{}

type LowercaseValidator struct{}

func (v LowercaseValidator) Description(context.Context) string {
	return "string value must be lowercase"
}

func (v LowercaseValidator) MarkdownDescription(context.Context) string {
	return "string value must be lowercase"
}

func (v LowercaseValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if value != strings.ToLower(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"String Value Must Be Lowercase",
			"The value must contain only lowercase characters.",
		)
	}
}

func Lowercase() LowercaseValidator {
	return LowercaseValidator{}
}
