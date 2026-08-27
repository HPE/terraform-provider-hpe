// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func StrToType(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}

	return types.StringValue(*s)
}

// TimeToType formats a nullable timestamp as an RFC 3339 string, returning a
// null value when the timestamp is nil.
func TimeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

func StrSliceToSet(items []string) types.Set {
	if len(items) == 0 {
		return types.SetNull(types.StringType)
	}

	var vals []attr.Value
	for _, i := range items {
		vals = append(vals, types.StringValue(i))
	}

	set, diags := types.SetValue(types.StringType, vals)
	if diags.HasError() {
		return types.SetNull(types.StringType)
	}

	return set
}

func SetToStrSlice(set types.Set) ([]string, error) {
	var items []string

	for _, elem := range set.Elements() {
		switch val := elem.(type) {
		case basetypes.StringValue:
			items = append(items, val.ValueString())
		default:
			return nil, fmt.Errorf("value %v is not a string", val)
		}
	}

	return items, nil
}

func SetToSlice[T any](ctx context.Context, set types.Set) ([]T, error) {
	var items []T
	for _, elem := range set.Elements() {
		a, err := ValueToAny(ctx, elem)
		if err != nil {
			return nil, err
		}

		v, ok := a.(T)
		if !ok {
			return nil, fmt.Errorf("converting Set to %T not possible, underlying type does not match", items)
		}

		items = append(items, v)
	}

	return items, nil
}

func BoolToType(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*b)
}

func StringToBool(ctx context.Context, s string) types.Bool {
	switch strings.ToLower(s) {
	case "on", "true", "yes", "1":
		return types.BoolValue(true)
	case "off", "false", "no", "0":
		return types.BoolValue(false)
	default:
		tflog.Debug(ctx, fmt.Sprintf("converting string to BoolNull: %s", s))

		return types.BoolNull()
	}
}

func BoolToStringOnOff(b bool) types.String {
	switch b {
	case true:
		return types.StringValue("on")
	case false:
		return types.StringValue("off")
	default:
		return types.StringNull()
	}
}

func BoolTypeToStringPointerOnOff(val types.Bool) *string {
	return BoolToStringOnOff(val.ValueBool()).ValueStringPointer()
}

func BoolToStringYesNo(b bool) types.String {
	switch b {
	case true:
		return types.StringValue("yes")
	case false:
		return types.StringValue("no")
	default:
		return types.StringNull()
	}
}

func BoolToStringTrueFalse(b bool) types.String {
	switch b {
	case true:
		return types.StringValue("true")
	case false:
		return types.StringValue("false")
	default:
		return types.StringNull()
	}
}

func Int64ToType(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*i)
}

// Int32ToType converts an optional int32 to a Terraform Int64 value, preserving
// null. Terraform has no 32-bit number type, so int32 fields (which the
// generated SDK uses for several models) widen to Int64.
func Int32ToType(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

func NumToType(f *float32) types.Number {
	if f == nil {
		return types.NumberNull()
	}

	// Convert float32 to big.Float
	bf := big.NewFloat(float64(*f))

	return types.NumberValue(bf)
}

func Int64SliceToSet(items []int64) types.Set {
	if len(items) == 0 {
		return types.SetNull(types.Int64Type)
	}

	var vals []attr.Value
	for _, i := range items {
		vals = append(vals, types.Int64Value(i))
	}

	set, diags := types.SetValue(types.Int64Type, vals)
	if diags.HasError() {
		return types.SetNull(types.Int64Type)
	}

	return set
}

type MappingFunc[I any, O any] func(in I) O

// Map objects in a slice into a Terraform Set Type according to the mapping
// function.
// Mapping function must take in an arbitrary Go value (e.g. API response
// struct) and return the mapped Terraform value.
func ToSetType[S any, O basetypes.ObjectValuable](
	ctx context.Context,
	slice []S,
	mapper MappingFunc[S, O],
) (basetypes.SetValue, diag.Diagnostics) {
	values := []attr.Value{}
	var obj O
	if len(slice) == 0 {
		return basetypes.NewSetNull(obj.Type(ctx)), nil
	}

	for _, i := range slice {
		v := mapper(i)

		values = append(values, v)
	}

	return types.SetValue(obj.Type(ctx), values)
}

// Map objects in a slice into a Terraform List Type according to the mapping
// function.
// Mapping function must take in an arbitrary Go value (e.g. API response
// struct) and return the mapped Terraform value.
func ToListType[S any, O basetypes.ObjectValuable](
	ctx context.Context,
	slice []S,
	mapper MappingFunc[S, O],
) (basetypes.ListValue, diag.Diagnostics) {
	values := []attr.Value{}
	var obj O

	if len(slice) == 0 {
		return basetypes.NewListNull(obj.Type(ctx)), nil
	}

	for _, i := range slice {
		v := mapper(i)

		values = append(values, v)
	}

	return types.ListValue(obj.Type(ctx), values)
}

// Map List objects into a slice of objects according to the mapping function.
// Mapping function must take in a Terraform value and return the mapped object
// (e.g. fill out an API request struct)
func FromSetType[S attr.Value, O any](
	ctx context.Context,
	set types.Set,
	mapper MappingFunc[S, O],
) ([]O, diag.Diagnostics) {
	var out []O
	var elems []S

	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}

	if diags := set.ElementsAs(ctx, &elems, false); diags.HasError() {
		return nil, diags
	}

	for _, el := range elems {
		out = append(out, mapper(el))
	}

	return out, nil
}

// Map List objects into a slice of objects according to the mapping function.
// Mapping function must take in a Terraform value and return the mapped object
// (e.g. fill out an API request struct)
func FromListType[S attr.Value, O any](
	ctx context.Context,
	list types.List,
	mapper MappingFunc[S, O],
) ([]O, diag.Diagnostics) {
	var out []O
	var elems []S

	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	if diags := list.ElementsAs(ctx, &elems, false); diags.HasError() {
		return nil, diags
	}

	for _, el := range elems {
		out = append(out, mapper(el))
	}

	return out, nil
}

// CookieTypeToAPI converts the Terraform cookie_type value ("session"/"persistence")
// to the API representation ("LBSessionCookieTime"/"LBPersistenceCookieTime").
func CookieTypeToAPI(val types.String) *string {
	if val.IsNull() || val.IsUnknown() {
		return nil
	}

	switch val.ValueString() {
	case "session":
		return PtrString("LBSessionCookieTime")
	case "persistence":
		return PtrString("LBPersistenceCookieTime")
	default:
		return nil
	}
}

// CookieTypeFromAPI converts the API cookie_type value ("LBSessionCookieTime"/"LBPersistenceCookieTime")
// back to the Terraform representation ("session"/"persistence").
func CookieTypeFromAPI(val *string) types.String {
	if val == nil {
		return types.StringNull()
	}

	switch *val {
	case "LBSessionCookieTime":
		return types.StringValue("session")
	case "LBPersistenceCookieTime":
		return types.StringValue("persistence")
	default:
		return types.StringNull()
	}
}

// PtrString is a helper to get a pointer to a string value.
func PtrString(s string) *string {
	return &s
}
