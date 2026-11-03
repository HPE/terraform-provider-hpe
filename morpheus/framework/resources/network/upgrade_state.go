// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// schemaVersion is the network resource's schema version. Version 1 made
// config write-only; state written by version 0 carries the configured config
// map and must be upgraded, because Terraform rejects a state that holds a
// value for a write-only attribute.
const schemaVersion = 1

// upgradeStashKey is the single key under which the version-0 config map is
// carried, inside config_from_api, from the state upgrade to the first Read.
//
// The upgrade cannot write private state, and the framework nulls write-only
// attributes after the upgrader runs, so the old value would be gone before
// Read could turn it into the digests that let the next plan recognise an
// unchanged configuration. Read consumes the stash — seeding the digests and
// replacing config_from_api with the API's view — so it exists only between
// the upgrade and the next refresh.
const upgradeStashKey = "_upgraded_from_state"

var _ resource.ResourceWithUpgradeState = &Resource{}

// UpgradeState upgrades version-0 state (config stored in state) to version 1.
func (r *Resource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	prior := networkSchemaV0(ctx)

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &prior,
			StateUpgrader: upgradeNetworkStateV0toV1,
		},
	}
}

// attributesAddedInV1 are the attributes version 1 added; none of them is in a
// version-0 state. The version-0 shape is pinned by testdata/state_v0.json, a
// state written by the previous release, so a change here or to the schema
// that is not reflected there fails the unit test.
var attributesAddedInV1 = []string{
	"availability_zone",
	"config_azure",
	"config_from_api",
	"config_gcp",
	"config_hvm_data",
	"config_hvm_overlay",
	"config_hvm_private",
	"config_hvm_sriov",
	"config_nsx_cloud",
	"config_nsxt",
	"config_openstack",
	"config_ovs_port_group",
	"network_server_id",
	"parent_network_id",
	"scope_id",
	"type_code",
	"vxlan_id",
}

// networkSchemaV0 is the version-0 schema: the current one without the
// attributes version 1 added, with config as the stored, computed dynamic it
// was, and type_id required as it was. The framework decodes a stored
// version-0 state against it before calling the upgrader.
func networkSchemaV0(ctx context.Context) schema.Schema {
	current := NetworkResourceSchema(ctx)
	attrs := make(map[string]schema.Attribute, len(current.Attributes))

	added := make(map[string]bool, len(attributesAddedInV1))
	for _, name := range attributesAddedInV1 {
		added[name] = true
	}

	for name, attr := range current.Attributes {
		switch {
		case added[name]:
			continue
		case name == "config":
			attrs[name] = schema.DynamicAttribute{Optional: true, Computed: true}
		case name == "type_id":
			attrs[name] = schema.Int64Attribute{Required: true}
		default:
			attrs[name] = attr
		}
	}

	return schema.Schema{Version: 0, Attributes: attrs, Blocks: current.Blocks}
}

// upgradeNetworkStateV0toV1 copies every attribute across, nulls the
// write-only ones, and stashes the version-0 config in config_from_api.
func upgradeNetworkStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old map[string]tftypes.Value
	if err := req.State.Raw.As(&old); err != nil {
		resp.Diagnostics.AddError("upgrade network resource state",
			"reading the version-0 state: "+err.Error())

		return
	}

	newType, ok := NetworkResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		resp.Diagnostics.AddError("upgrade network resource state", "the current schema is not an object")

		return
	}

	upgraded := make(map[string]tftypes.Value, len(newType.AttributeTypes))

	for name, typ := range newType.AttributeTypes {
		switch name {
		case "config":
			// Write-only now: never in state. Its version-0 value is stashed
			// in config_from_api below.
			upgraded[name] = tftypes.NewValue(typ, nil)
		case "config_from_api":
			upgraded[name] = stashConfig(old["config"], typ)
		default:
			// Carried across when version 0 had it; the attributes version 1
			// added (write-only blocks, type_code, network_server_id, ...)
			// start null and are filled by the next Read.
			if v, present := old[name]; present {
				upgraded[name] = v
			} else {
				upgraded[name] = tftypes.NewValue(typ, nil)
			}
		}
	}

	resp.State.Raw = tftypes.NewValue(newType, upgraded)
}

// stashConfig wraps the version-0 config value as {_upgraded_from_state: value}
// for the config_from_api slot, or null when there was no config.
func stashConfig(oldConfig tftypes.Value, slot tftypes.Type) tftypes.Value {
	if !oldConfig.IsKnown() || oldConfig.IsNull() {
		return tftypes.NewValue(slot, nil)
	}

	wrapper := tftypes.Object{AttributeTypes: map[string]tftypes.Type{upgradeStashKey: oldConfig.Type()}}

	return tftypes.NewValue(wrapper, map[string]tftypes.Value{upgradeStashKey: oldConfig})
}

// upgradeStash returns the version-0 config carried in config_from_api by the
// state upgrade, if that is what the value is. Anything else — the API's view
// of the config, or null — is not a stash.
func upgradeStash(ctx context.Context, fromAPI types.Dynamic) (map[string]any, bool, error) {
	if fromAPI.IsNull() || fromAPI.IsUnknown() || fromAPI.IsUnderlyingValueNull() || fromAPI.IsUnderlyingValueUnknown() {
		return nil, false, nil
	}

	raw, err := convert.ValueToAny(ctx, fromAPI.UnderlyingValue())
	if err != nil {
		return nil, false, fmt.Errorf("reading config_from_api: %w", err)
	}

	outer, ok := raw.(map[string]any)
	if !ok || len(outer) != 1 {
		return nil, false, nil
	}

	inner, ok := outer[upgradeStashKey].(map[string]any)
	if !ok {
		return nil, false, nil
	}

	return inner, true, nil
}
