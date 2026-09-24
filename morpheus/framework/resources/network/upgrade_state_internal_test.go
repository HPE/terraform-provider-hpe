// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// The version-0 → 1 upgrade must null the write-only attributes and carry the
// old config to Read inside config_from_api, under the stash key.
func TestUnitUpgradeStateV0toV1StashesConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	prior := networkSchemaV0(ctx)
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	// A version-0 state: every attribute null except the ones we care about.
	oldAttrs := map[string]tftypes.Value{}
	for name, typ := range priorType.AttributeTypes {
		oldAttrs[name] = tftypes.NewValue(typ, nil)
	}

	oldAttrs["id"] = tftypes.NewValue(tftypes.Number, 252537)
	oldAttrs["name"] = tftypes.NewValue(tftypes.String, "n")

	cfgType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"hvmVlanTrunks": tftypes.String, "hvmTrunkOnly": tftypes.Bool}}
	oldAttrs["config"] = tftypes.NewValue(cfgType, map[string]tftypes.Value{
		"hvmVlanTrunks": tftypes.NewValue(tftypes.String, "100,200"),
		"hvmTrunkOnly":  tftypes.NewValue(tftypes.Bool, true),
	})

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: tftypes.NewValue(priorType, oldAttrs)}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: NetworkResourceSchema(ctx)}}

	upgradeNetworkStateV0toV1(ctx, req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var upgraded NetworkModel
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatal(diags)
	}

	if upgraded.Id.ValueInt64() != 252537 || upgraded.Name.ValueString() != "n" {
		t.Errorf("attributes must be carried across: %+v", upgraded)
	}

	if !upgraded.Config.IsNull() || !upgraded.ConfigOvsPortGroup.IsNull() {
		t.Error("write-only attributes must be null after upgrade")
	}

	if !upgraded.NetworkServerId.IsNull() {
		t.Error("attributes new in v1 must be null after upgrade")
	}

	stash, found, err := upgradeStash(ctx, upgraded.ConfigFromApi)
	if err != nil || !found {
		t.Fatalf("expected the old config stashed in config_from_api; found=%v err=%v value=%v", found, err, upgraded.ConfigFromApi)
	}

	if stash["hvmVlanTrunks"] != "100,200" || stash["hvmTrunkOnly"] != true {
		t.Errorf("stash content wrong: %v", stash)
	}

	// Seeding from the stash must digest identically to the same values coming
	// from the configuration — otherwise the first plan after upgrading would
	// show a change nobody made.
	ovs := typeSpecFor("ovsPortGroup")
	if ovs == nil {
		t.Fatal("the catalogue must know ovsPortGroup")
	}

	fromStash, diags := resolveConfigInputs(ctx, seedModel(t, ctx, stash), ovs)
	if diags.HasError() {
		t.Fatal(diags)
	}

	configured := seedModel(t, ctx, map[string]any{"hvmVlanTrunks": "100,200", "hvmTrunkOnly": true})

	fromConfig, diags := resolveConfigInputs(ctx, configured, ovs)
	if diags.HasError() {
		t.Fatal(diags)
	}

	if fromStash.digests != fromConfig.digests {
		t.Errorf("stash and configuration must digest equally: %+v vs %+v", fromStash.digests, fromConfig.digests)
	}

	// The same values through the typed block digest identically too: the
	// digest is over the routed config map, not over how it was expressed.
	typed := seedModel(t, ctx, nil)
	typed.ConfigOvsPortGroup = NewConfigOvsPortGroupValueMust(ConfigOvsPortGroupValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"vlan_trunks": types.StringValue("100,200"),
		"trunk_only":  types.BoolValue(true),
	})

	fromBlock, diags := resolveConfigInputs(ctx, typed, ovs)
	if diags.HasError() {
		t.Fatal(diags)
	}

	if fromBlock.digests != fromConfig.digests || fromBlock.block == nil || fromBlock.block.Attr != "config_ovs_port_group" {
		t.Errorf("typed block and generic map must digest equally: %+v vs %+v (block %v)",
			fromBlock.digests, fromConfig.digests, fromBlock.block)
	}
}

// A config_from_api holding the API's view is not a stash, and neither is null.
func TestUnitUpgradeStashRecognisesOnlyTheStash(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	if _, found, _ := upgradeStash(ctx, types.DynamicNull()); found {
		t.Error("null is not a stash")
	}

	apiView := seedModel(t, ctx, map[string]any{"hvmVlanTrunks": "100,200", "hvmNetworkLayout": "ovs"}).Config
	if _, found, _ := upgradeStash(ctx, apiView); found {
		t.Error("the API's view of config must not be mistaken for a stash")
	}

	// Even a single-key map is not a stash unless the key is the marker.
	single := seedModel(t, ctx, map[string]any{"hvmTrunkOnly": false}).Config
	if _, found, _ := upgradeStash(ctx, single); found {
		t.Error("a single arbitrary key is not a stash")
	}
}

func TestUnitStashConfigNullWhenNoConfig(t *testing.T) {
	t.Parallel()

	slot := tftypes.DynamicPseudoType
	if v := stashConfig(tftypes.NewValue(tftypes.DynamicPseudoType, nil), slot); !v.IsNull() {
		t.Error("no version-0 config → null stash")
	}

	if v := stashConfig(tftypes.NewValue(tftypes.DynamicPseudoType, tftypes.UnknownValue), slot); !v.IsNull() {
		t.Error("unknown version-0 config → null stash")
	}
}

// seedModel builds a model whose write-only config holds m, as
// resolveConfigInputs would see it from the configuration.
func seedModel(t *testing.T, ctx context.Context, m map[string]any) NetworkModel {
	t.Helper()

	dyn, err := convert.MapToDynamic(ctx, m)
	if err != nil {
		t.Fatal(err)
	}

	model := NetworkModel{Config: dyn, ConnectedGateway: types.StringNull()}
	model.ConfigOvsPortGroup = NewConfigOvsPortGroupValueNull()

	return model
}

// The version-0 schema must describe what the previous release actually wrote,
// and the upgrader must accept it. testdata/state_v0.json is the attributes
// object of an hpe_morpheus_network instance from a state file written by the
// 2.1 line (schema_version 0), which is exactly what Terraform hands the
// provider as the raw state to upgrade. Decoding uses the same options the
// framework uses (fwserver/server_upgraderesourcestate.go).
func TestUnitUpgradeStateAcceptsAStateWrittenByTheOldSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	raw, err := os.ReadFile("testdata/state_v0.json")
	if err != nil {
		t.Fatal(err)
	}

	var stored map[string]json.RawMessage
	if err := json.Unmarshal(raw, &stored); err != nil {
		t.Fatal(err)
	}

	prior := networkSchemaV0(ctx)
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	// The prior schema has exactly the attributes the old state has: nothing
	// the old release never wrote, nothing it wrote missing.
	for name := range stored {
		if _, ok := priorType.AttributeTypes[name]; !ok {
			t.Errorf("the version-0 state has %q, which the version-0 schema lacks", name)
		}
	}

	for name := range priorType.AttributeTypes {
		if _, ok := stored[name]; !ok {
			t.Errorf("the version-0 schema has %q, which no version-0 state holds; add it to attributesAddedInV1", name)
		}
	}

	// And every attribute version 1 added is one the current schema has.
	current := NetworkResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	for _, name := range attributesAddedInV1 {
		if _, ok := current.AttributeTypes[name]; !ok {
			t.Errorf("attributesAddedInV1 names %q, which the current schema does not have", name)
		}
	}

	// Decode as the framework does — and also strictly, which proves the
	// prior schema has no key the state lacks a value for.
	rawState := tfprotov6.RawState{JSON: raw}

	value, err := rawState.UnmarshalWithOpts(priorType, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		t.Fatalf("the framework's decode of the old state against the version-0 schema failed: %v", err)
	}

	if _, err := rawState.Unmarshal(priorType); err != nil {
		t.Fatalf("strict decode of the old state against the version-0 schema failed: %v", err)
	}

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: value}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: NetworkResourceSchema(ctx)}}

	upgradeNetworkStateV0toV1(ctx, req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var upgraded NetworkModel
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatal(diags)
	}

	if upgraded.Id.ValueInt64() != 252799 || upgraded.TypeId.ValueInt64() != 1 || upgraded.CloudId.ValueInt64() != 17 ||
		upgraded.Cidr.ValueString() != "10.208.1.0/24" {
		t.Errorf("attributes must be carried across: id=%v type_id=%v cloud_id=%v cidr=%v",
			upgraded.Id, upgraded.TypeId, upgraded.CloudId, upgraded.Cidr)
	}

	if !upgraded.Config.IsNull() || !upgraded.TypeCode.IsNull() || !upgraded.NetworkServerId.IsNull() {
		t.Error("write-only and version-1 attributes must be null after the upgrade")
	}

	stash, found, err := upgradeStash(ctx, upgraded.ConfigFromApi)
	if err != nil || !found || stash["spikeKey"] != "v1" || stash["another"] != "x" {
		t.Errorf("the old config must be stashed in config_from_api: found=%v err=%v stash=%v", found, err, stash)
	}
}
