// Package schemadefaults fills schema-declared defaults into resource state
// after a Read.
//
// The plugin framework applies an attribute's Default only during planning, when
// the configuration value is null. A resource's Read path, by contrast, maps the
// API response straight into state, and when the API omits an optional field the
// mapped value is null. Terraform then compares that null in state against the
// default the plan produces and sees a change, so an imported resource plans an
// update that can never settle -- and where the attribute forces replacement, a
// destroy/recreate of a resource nobody touched.
//
// Apply closes that gap generically: after Read, any attribute that declares a
// default and is null in state is set to its default, exactly as the plan would.
// It walks the schema for defaults and rewrites only null leaves, so a value the
// API did populate is left untouched, and an omitted nested block (null in state)
// keeps its children null -- which is correct, because the plan sees that block
// null too.
//
// It resolves the null-in-state class (MORPH-16192). It does not, and cannot,
// fix a Read that panics or returns an error before it completes: Apply runs
// after Read, so those must be guarded in the Read itself.
package schemadefaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Apply sets schema-declared defaults into any null leaf of state that has one.
// It is a no-op when state is null (e.g. the resource is gone) or the schema
// declares no defaults. Call it at the end of a resource's Read, after state has
// been populated:
//
//	schemadefaults.Apply(ctx, MyResourceSchema(ctx), &resp.State)
func Apply(ctx context.Context, s schema.Schema, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if state == nil || state.Raw.IsNull() || !state.Raw.IsKnown() {
		return diags
	}

	defs := collectDefaults(ctx, s)
	if len(defs) == 0 {
		return diags
	}

	newRaw, err := fillDefaults(ctx, state.Raw, defs)
	if err != nil {
		diags.AddError(
			"apply schema defaults",
			"failed to fill schema defaults into state: "+err.Error(),
		)

		return diags
	}

	state.Raw = newRaw

	return diags
}

// collectDefaults walks a resource schema and returns a map from an attribute's
// collapsed name-path (attribute names joined by "/", element indices dropped) to
// its declared default value, for every String, Bool and Int64 attribute that
// declares one, descending through Single/List/Set nested attributes.
//
// Keys retain the full sequence of attribute names, so a leaf that appears under
// two different parents (config_hvm/no_agent vs config_vmware/no_agent) gets two
// distinct keys and cannot collide. Attribute names cannot contain "/", so a
// nested path can never collide with a literal attribute name either.
func collectDefaults(ctx context.Context, s schema.Schema) map[string]attr.Value {
	out := map[string]attr.Value{}
	walkAttributes(ctx, "", s.Attributes, out)

	return out
}

func walkAttributes(
	ctx context.Context,
	prefix string,
	attrs map[string]schema.Attribute,
	out map[string]attr.Value,
) {
	for name, a := range attrs {
		key := name
		if prefix != "" {
			key = prefix + "/" + name
		}

		switch t := a.(type) {
		case schema.StringAttribute:
			if t.Default != nil {
				var r defaults.StringResponse
				t.Default.DefaultString(ctx, defaults.StringRequest{}, &r)
				out[key] = r.PlanValue
			}
		case schema.BoolAttribute:
			if t.Default != nil {
				var r defaults.BoolResponse
				t.Default.DefaultBool(ctx, defaults.BoolRequest{}, &r)
				out[key] = r.PlanValue
			}
		case schema.Int64Attribute:
			if t.Default != nil {
				var r defaults.Int64Response
				t.Default.DefaultInt64(ctx, defaults.Int64Request{}, &r)
				out[key] = r.PlanValue
			}
		case schema.SingleNestedAttribute:
			walkAttributes(ctx, key, t.Attributes, out)
		case schema.ListNestedAttribute:
			walkAttributes(ctx, key, t.NestedObject.Attributes, out)
		case schema.SetNestedAttribute:
			walkAttributes(ctx, key, t.NestedObject.Attributes, out)
		}
	}
}

// fillDefaults returns a copy of raw in which every null String, Bool or Number
// leaf whose collapsed name-path is present in defs is replaced by that default.
// tftypes.Transform visits every value with its full concrete path; children of a
// null container are never visited, so an omitted nested block stays null.
func fillDefaults(
	ctx context.Context,
	raw tftypes.Value,
	defs map[string]attr.Value,
) (tftypes.Value, error) {
	return tftypes.Transform(raw, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if !v.IsNull() {
			return v, nil
		}

		ty := v.Type()
		if !ty.Is(tftypes.String) && !ty.Is(tftypes.Bool) && !ty.Is(tftypes.Number) {
			return v, nil
		}

		def, ok := defs[normalizePath(p)]
		if !ok {
			return v, nil
		}

		return def.ToTerraformValue(ctx)
	})
}

// normalizePath collapses a concrete attribute path to the sequence of attribute
// names joined by "/", dropping element steps (list index, set value, map key) so
// every element of a collection maps to the same schema default.
func normalizePath(p *tftypes.AttributePath) string {
	var b []byte
	for _, step := range p.Steps() {
		name, ok := step.(tftypes.AttributeName)
		if !ok {
			continue
		}
		if len(b) > 0 {
			b = append(b, '/')
		}
		b = append(b, string(name)...)
	}

	return string(b)
}
