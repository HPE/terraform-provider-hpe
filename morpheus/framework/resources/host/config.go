// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// genericConfig converts the free-form config dynamic value (a map) into a
// HostConfigGeneric. Keys that match known fields are mapped; everything else is
// carried through AdditionalProperties so the type's option types receive it.
func genericConfig(
	ctx context.Context,
	plan HostModel,
) (*sdk.HostConfigGeneric, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := &sdk.HostConfigGeneric{
		AdditionalProperties: map[string]interface{}{},
	}

	underlying := plan.Config.UnderlyingValue()
	obj, ok := underlying.(types.Object)
	if !ok {
		// ValidObjectMap validator guarantees this is an object map; be defensive.
		diags.AddError(
			"invalid config",
			"config must be a map/object of configuration keys",
		)

		return nil, diags
	}

	for k, v := range obj.Attributes() {
		gv, gd := attrToGo(ctx, v)
		diags.Append(gd...)
		if diags.HasError() {
			return nil, diags
		}

		switch k {
		case "userData":
			if s, ok := gv.(string); ok {
				cfg.UserData = &s
			}
		default:
			cfg.AdditionalProperties[k] = gv
		}
	}

	return cfg, diags
}

// attrToGo converts a framework attr.Value to a plain Go value for JSON
// encoding. Nulls become nil, objects and maps become map[string]interface{}
// and lists, sets and tuples become []interface{}, so a type's option types
// receive structure rather than its textual rendering.
func attrToGo(ctx context.Context, v attr.Value) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil, diags
	}

	switch tv := v.(type) {
	case types.String:
		return tv.ValueString(), diags
	case types.Bool:
		return tv.ValueBool(), diags
	case types.Int64:
		return tv.ValueInt64(), diags
	case types.Number:
		f, _ := tv.ValueBigFloat().Float64()

		return f, diags
	case types.Dynamic:
		return attrToGo(ctx, tv.UnderlyingValue())
	case types.Object:
		return attrsToGo(ctx, tv.Attributes())
	case types.Map:
		return attrsToGo(ctx, tv.Elements())
	case types.List:
		return elemsToGo(ctx, tv.Elements())
	case types.Set:
		return elemsToGo(ctx, tv.Elements())
	case types.Tuple:
		return elemsToGo(ctx, tv.Elements())
	default:
		diags.AddError(
			"invalid config",
			fmt.Sprintf("config value of unsupported type %T", v),
		)

		return nil, diags
	}
}

func attrsToGo(ctx context.Context, attrs map[string]attr.Value) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	out := make(map[string]interface{}, len(attrs))
	for k, v := range attrs {
		gv, gd := attrToGo(ctx, v)
		diags.Append(gd...)
		if diags.HasError() {
			return nil, diags
		}

		out[k] = gv
	}

	return out, diags
}

func elemsToGo(ctx context.Context, elems []attr.Value) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	out := make([]interface{}, 0, len(elems))
	for _, v := range elems {
		gv, gd := attrToGo(ctx, v)
		diags.Append(gd...)
		if diags.HasError() {
			return nil, diags
		}

		out = append(out, gv)
	}

	return out, diags
}
