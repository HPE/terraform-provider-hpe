package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable                   = (*NormalizedLineEndingsString)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*NormalizedLineEndingsString)(nil)
)

// NormalizedLineEndingsString is a custom String value that treats strings which
// differ only in their line endings (CRLF or CR vs LF) as semantically equal.
// This suits values such as script bodies, which a practitioner may author with
// Windows (CRLF) line endings but the API stores/returns with LF — avoiding both
// the "inconsistent result after apply" and perpetual-diff problems without
// rewriting the planned value (which Terraform Core rejects for a configured
// attribute).
type NormalizedLineEndingsString struct {
	basetypes.StringValue
}

// normalizeLineEndings folds CRLF and CR line endings to LF.
func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")

	return strings.ReplaceAll(s, "\r", "\n")
}

// Type returns a NormalizedLineEndingsStringType.
func (v NormalizedLineEndingsString) Type(_ context.Context) attr.Type {
	return NormalizedLineEndingsStringType{}
}

// Equal returns true if the given value is equivalent.
func (v NormalizedLineEndingsString) Equal(o attr.Value) bool {
	other, ok := o.(NormalizedLineEndingsString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the two strings are equal once their line
// endings are normalized to LF.
func (v NormalizedLineEndingsString) StringSemanticEquals(
	_ context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(NormalizedLineEndingsString)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return normalizeLineEndings(v.ValueString()) ==
		normalizeLineEndings(newValue.ValueString()), diags
}

// NewNormalizedLineEndingsStringNull creates a NormalizedLineEndingsString with
// a null value. Determine whether the value is null via the IsNull method.
func NewNormalizedLineEndingsStringNull() NormalizedLineEndingsString {
	return NormalizedLineEndingsString{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewNormalizedLineEndingsStringUnknown creates a NormalizedLineEndingsString
// with an unknown value. Determine whether the value is unknown via the
// IsUnknown method.
func NewNormalizedLineEndingsStringUnknown() NormalizedLineEndingsString {
	return NormalizedLineEndingsString{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewNormalizedLineEndingsStringValue creates a NormalizedLineEndingsString with
// a known value. Access the value via the ValueString method.
func NewNormalizedLineEndingsStringValue(value string) NormalizedLineEndingsString {
	return NormalizedLineEndingsString{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewNormalizedLineEndingsStringPointerValue creates a NormalizedLineEndingsString
// with a null value if nil, or a known value. Access the value via the
// ValueStringPointer method.
func NewNormalizedLineEndingsStringPointerValue(value *string) NormalizedLineEndingsString {
	return NormalizedLineEndingsString{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
