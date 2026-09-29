// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = RequiredWhenSiblingStringValidator{}

// RequiredWhenSiblingStringValidator validates that this string attribute is
// set to a non-empty value whenever a sibling root-level string attribute holds
// one of MatchValues (or, when MatchWhenNull is true, is itself unset). This
// mirrors conditional-required behaviour that the Morpheus API enforces
// server-side but does not surface until a request is rejected (for example,
// option_list.source_url is required when type == "rest", and the server
// defaults an unset type to "rest").
type RequiredWhenSiblingStringValidator struct {
	SiblingName   string
	MatchValues   []string
	MatchWhenNull bool
}

func (v RequiredWhenSiblingStringValidator) Description(context.Context) string {
	if v.MatchWhenNull {
		return fmt.Sprintf(
			"must be set when %q is one of %v or is unset", v.SiblingName, v.MatchValues,
		)
	}

	return fmt.Sprintf(
		"must be set when %q is one of %v", v.SiblingName, v.MatchValues,
	)
}

func (v RequiredWhenSiblingStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v RequiredWhenSiblingStringValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	var sibling types.String
	diags := request.Config.GetAttribute(ctx, path.Root(v.SiblingName), &sibling)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// The sibling value is not yet resolved; defer — Terraform re-validates
	// once it is known.
	if sibling.IsUnknown() {
		return
	}

	matched := false
	if sibling.IsNull() {
		matched = v.MatchWhenNull
	} else {
		for _, m := range v.MatchValues {
			if sibling.ValueString() == m {
				matched = true

				break
			}
		}
	}

	if !matched {
		return
	}

	// This attribute may legitimately be unknown at plan time (e.g. sourced
	// from another resource). Let Terraform re-validate once it resolves.
	if request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.IsNull() || strings.TrimSpace(request.ConfigValue.ValueString()) == "" {
		siblingDesc := fmt.Sprintf("%q is one of %v", v.SiblingName, v.MatchValues)
		if v.MatchWhenNull {
			siblingDesc = fmt.Sprintf("%q is one of %v or is unset", v.SiblingName, v.MatchValues)
		}
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				"Missing required attribute",
				fmt.Sprintf(
					"Attribute %q must be set to a non-empty value when %s.",
					request.Path, siblingDesc,
				),
			),
		)
	}
}

// RequiredWhenSiblingEquals returns a validator that requires this string
// attribute to be non-empty when the named sibling attribute equals one of the
// provided values.
func RequiredWhenSiblingEquals(siblingName string, values ...string) validator.String {
	return RequiredWhenSiblingStringValidator{
		SiblingName:   siblingName,
		MatchValues:   values,
		MatchWhenNull: false,
	}
}

// RequiredWhenSiblingEqualsOrNull behaves like RequiredWhenSiblingEquals but
// also triggers when the sibling attribute is unset (null). This is useful when
// the server applies a default to the sibling that itself makes this attribute
// required.
func RequiredWhenSiblingEqualsOrNull(siblingName string, values ...string) validator.String {
	return RequiredWhenSiblingStringValidator{
		SiblingName:   siblingName,
		MatchValues:   values,
		MatchWhenNull: true,
	}
}
