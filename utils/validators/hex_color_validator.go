// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// hexColorRegexp matches CSS hex color codes with 3, 4, 6, or 8 hex digits,
// for example #fff, #abcd, #1a73e8, or #1a73e8ff.
var hexColorRegexp = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var _ validator.String = HexColorValidator{}

// HexColorValidator rejects a configured string that is not a valid hex color
// code, so an invalid color is reported at validate/plan time instead of being
// silently accepted by the API.
type HexColorValidator struct{}

func (v HexColorValidator) Description(_ context.Context) string {
	return "string value must be a valid hex color code, for example #1a73e8, #fff, or #1a73e8ff"
}

func (v HexColorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v HexColorValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if !hexColorRegexp.MatchString(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Hex Color",
			"Attribute "+request.Path.String()+" must be a valid hex color code, "+
				"for example #1a73e8, #fff, or #1a73e8ff, got: "+value,
		)
	}
}

// HexColor returns a validator that ensures a string is a valid hex color code.
func HexColor() HexColorValidator {
	return HexColorValidator{}
}
