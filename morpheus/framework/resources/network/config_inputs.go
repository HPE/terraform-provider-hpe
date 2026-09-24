// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// configInputs are the write-only configuration values as configured, resolved
// into what is sent to Morpheus. They are read from the configuration, never
// from the plan or state: a write-only value is null in both.
type configInputs struct {
	// all is the config map to send; nil when nothing is configured.
	all map[string]any
	// immutable is the subset Morpheus will not change after create.
	immutable map[string]any
	// block is the typed block in use, or nil, and blockValues its set
	// attributes by name.
	block       *networkconfig.Block
	blockValues map[string]any
	// generic is the free-form config map as configured; nil when unset.
	generic map[string]any
	// typ is the network type the inputs were routed for, or nil when it is not
	// known (or not in the catalogue).
	typ *networkconfig.TypeSpec
	// digests are the sha256 digests of all and immutable.
	digests networkconfig.Digests
}

// resolveConfigInputs reads the write-only configuration from config: the
// free-form "config" map, the typed block if one is set, and the top-level
// connected_gateway, which the NSX-T types read from the config map. typ is the
// network's type when known; it decides which values are immutable. Unknown
// values (a block wired to another resource's output) resolve to nothing, so a
// plan with unknowns neither sends nor digests them; apply sees the known values.
func resolveConfigInputs(
	ctx context.Context,
	config NetworkModel,
	typ *networkconfig.TypeSpec,
) (configInputs, diag.Diagnostics) {
	var diags diag.Diagnostics

	generic, gdiags := genericConfig(ctx, config.Config)
	diags.Append(gdiags...)
	if diags.HasError() {
		return configInputs{}, diags
	}

	blocks, bdiags := configuredBlocks(ctx, config)
	diags.Append(bdiags...)
	if diags.HasError() {
		return configInputs{}, diags
	}

	if len(blocks) > 1 {
		names := make([]string, 0, len(blocks))
		for _, b := range blocks {
			names = append(names, b.block.Attr)
		}

		diags.AddError("network configuration",
			"only one typed configuration block may be set, got "+strings.Join(names, ", "))

		return configInputs{}, diags
	}

	var (
		block  *networkconfig.Block
		values map[string]any
	)

	if len(blocks) == 1 {
		block, values = &blocks[0].block, blocks[0].values
	}

	all, immutable := networkconfig.Route(typ, block, generic, values)

	// connected_gateway is a top-level attribute for compatibility; Morpheus
	// reads it from the config map.
	if !config.ConnectedGateway.IsNull() && !config.ConnectedGateway.IsUnknown() {
		all["connectedGateway"] = config.ConnectedGateway.ValueString()
	}

	digests, err := networkconfig.Compute(all, immutable)
	if err != nil {
		diags.AddError("network configuration", err.Error())

		return configInputs{}, diags
	}

	if len(all) == 0 {
		all = nil
	}

	return configInputs{
		all: all, immutable: immutable, block: block, blockValues: values, generic: generic, typ: typ, digests: digests,
	}, diags
}

// genericConfig converts the free-form config value to a map; nil when unset
// or unknown.
func genericConfig(ctx context.Context, value types.Dynamic) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() || value.IsUnderlyingValueNull() || value.IsUnderlyingValueUnknown() {
		return nil, diags
	}

	raw, err := convert.ValueToAny(ctx, value.UnderlyingValue())
	if err != nil {
		diags.AddError("network configuration", "config: failed to convert: "+err.Error())

		return nil, diags
	}

	m, ok := raw.(map[string]any)
	if !ok {
		diags.AddError("network configuration", "config must be an object of keys to values")

		return nil, diags
	}

	return m, diags
}

// configuredBlock is a typed block that is set in the configuration, with its
// known, non-null attribute values by attribute name.
type configuredBlock struct {
	block  networkconfig.Block
	values map[string]any
}

// configuredBlocks finds the typed blocks set in config. The model's fields
// are matched to the declared blocks by their tfsdk tag, so a block added to
// the schema and to networkconfig.Blocks needs no code here.
func configuredBlocks(ctx context.Context, config NetworkModel) ([]configuredBlock, diag.Diagnostics) {
	var (
		diags diag.Diagnostics
		out   []configuredBlock
	)

	rv := reflect.ValueOf(config)
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("tfsdk")

		b, ok := networkconfig.BlockFor(tag)
		if !ok {
			continue
		}

		value, ok := rv.Field(i).Interface().(basetypes.ObjectValuable)
		if !ok {
			diags.AddError("network configuration",
				fmt.Sprintf("%s: the schema attribute is not an object; the resource and its block declaration disagree", tag))

			continue
		}

		if value.IsNull() || value.IsUnknown() {
			continue
		}

		obj, odiags := value.ToObjectValue(ctx)
		diags.Append(odiags...)
		if odiags.HasError() {
			continue
		}

		values, vdiags := blockValues(ctx, b, obj.Attributes())
		diags.Append(vdiags...)

		out = append(out, configuredBlock{block: b, values: values})
	}

	return out, diags
}

// blockValues converts a block's known, non-null attribute values to plain Go
// values, keyed by attribute name. Attributes the declaration does not know
// are an inconsistency between schema and declaration, reported once.
func blockValues(ctx context.Context, b networkconfig.Block, attrs map[string]attr.Value) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := make(map[string]any, len(attrs))

	for name, v := range attrs {
		if _, ok := b.Key(name); !ok {
			diags.AddError("network configuration",
				fmt.Sprintf("%s.%s: the schema has this attribute but the block declaration does not", b.Attr, name))

			continue
		}

		if v.IsNull() || v.IsUnknown() {
			continue
		}

		raw, err := convert.ValueToAny(ctx, v)
		if err != nil {
			diags.AddError("network configuration", fmt.Sprintf("%s.%s: failed to convert: %s", b.Attr, name, err.Error()))

			continue
		}

		values[name] = raw
	}

	return values, diags
}

// configFromAPI converts the config Morpheus returned into the read-only
// config_from_api value: the whole map, as the server holds it, including keys
// it added and values it masked. Null when the network has no config.
func configFromAPI(ctx context.Context, id int64, cfg *sdk.GetNetwork200ResponseNetworkConfig) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	if cfg == nil {
		return types.DynamicNull(), diags
	}

	// Through JSON rather than ToMap: the SDK's config struct types the keys
	// the spec knows and keeps the rest as additional properties, and ToMap
	// leaves the typed ones as pointers. The wire form is uniform.
	raw, err := json.Marshal(cfg)
	if err != nil {
		diags.AddError("populate network resource",
			fmt.Sprintf("network %d: failed to read config from the API response: %s", id, err.Error()))

		return types.DynamicNull(), diags
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		diags.AddError("populate network resource",
			fmt.Sprintf("network %d: failed to decode config from the API response: %s", id, err.Error()))

		return types.DynamicNull(), diags
	}

	if len(m) == 0 {
		return types.DynamicNull(), diags
	}

	value, err := convert.MapToDynamic(ctx, m)
	if err != nil {
		diags.AddError("populate network resource",
			fmt.Sprintf("network %d: failed to convert config_from_api: %s", id, err.Error()))

		return types.DynamicNull(), diags
	}

	return value, diags
}
