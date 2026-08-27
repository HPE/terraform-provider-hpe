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
	_ basetypes.StringValuable                   = (*TrimmedString)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*TrimmedString)(nil)
)

type TrimmedString struct {
	basetypes.StringValue
}

// Type returns a TrimmedStringType.
func (v TrimmedString) Type(_ context.Context) attr.Type {
	return TrimmedStringType{}
}

// Equal returns true if the given value is equivalent.
func (v TrimmedString) Equal(o attr.Value) bool {
	other, ok := o.(TrimmedString)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true if the given strings are equal once the
// leading and trailing whtiespace is removed
func (v TrimmedString) StringSemanticEquals(
	_ context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(TrimmedString)
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

	newTrimmed := strings.TrimSpace(newValue.ValueString())
	vTrimmed := strings.TrimSpace(v.ValueString())

	return newTrimmed == vTrimmed, diags
}

// NewNormalizedNull creates a Normalized with a null value. Determine whether
// the value is null via IsNull method.
func NewTrimmedStringNull() TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewTrimmedStringUnknown creates a Normalized with an unknown value. Determine
// whether the value is unknown via IsUnknown method.
func NewTrimmedStringUnknown() TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewTrimmedStringValue creates a Normalized with a known value. Access the
// value via ValueString method.
func NewTrimmedStringValue(value string) TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewTrimmedStringPointerValue creates a Normalized with a null value if nil or
// a known value. Access the value via ValueStringPointer method.
func NewTrimmedStringPointerValue(value *string) TrimmedString {
	return TrimmedString{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
