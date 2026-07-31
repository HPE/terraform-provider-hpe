// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setMustContainAndAllowValidator struct {
	required []string
	allowed  []string
}

func SetMustContainAndAllow(required []string, allowed []string) validator.Set {
	return setMustContainAndAllowValidator{
		required: required,
		allowed:  allowed,
	}
}

func (v setMustContainAndAllowValidator) Description(ctx context.Context) string {
	return "Validates required and allowed values in set"
}

func (v setMustContainAndAllowValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v setMustContainAndAllowValidator) ValidateSet(
	ctx context.Context,
	req validator.SetRequest,
	resp *validator.SetResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Convert set → map[string]bool
	values := map[string]bool{}
	for _, elem := range req.ConfigValue.Elements() {
		str := elem.(types.String).ValueString()
		values[str] = true
	}

	// ✅ Check required values
	for _, r := range v.required {
		if !values[r] {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing required value",
				fmt.Sprintf("Value '%s' must be present in the set.", r),
			)
		}
	}

	// ✅ Check allowed values
	for val := range values {
		if !contains(v.allowed, val) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid value",
				fmt.Sprintf("Value '%s' is not allowed. Allowed values: %v", val, v.allowed),
			)
		}
	}
}

// helper
func contains(list []string, v string) bool {
	for _, e := range list {
		if e == v {
			return true
		}
	}

	return false
}
