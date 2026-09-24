// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
)

func ovsBlock(ctx context.Context, trunks, trunkOnly attr.Value) ConfigOvsPortGroupValue {
	return NewConfigOvsPortGroupValueMust(ConfigOvsPortGroupValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"vlan_trunks": trunks,
		"trunk_only":  trunkOnly,
	})
}

func TestUnitConfiguredBlocksAreFoundByTag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Every declared block must be an attribute of the model, or the schema and
	// the declaration have drifted apart.
	for _, b := range networkconfig.Blocks {
		if _, ok := NetworkResourceSchema(ctx).Attributes[b.Attr]; !ok {
			t.Errorf("block %s is declared but the schema has no such attribute", b.Attr)
		}
	}

	model := seedModel(t, ctx, nil)

	blocks, diags := configuredBlocks(ctx, model)
	if diags.HasError() || len(blocks) != 0 {
		t.Fatalf("a null block must not be reported: %v %v", blocks, diags)
	}

	model.ConfigOvsPortGroup = ovsBlock(ctx, types.StringValue("100,200"), types.BoolUnknown())

	blocks, diags = configuredBlocks(ctx, model)
	if diags.HasError() || len(blocks) != 1 || blocks[0].block.Attr != "config_ovs_port_group" {
		t.Fatalf("the set block must be found: %v %v", blocks, diags)
	}

	if blocks[0].values["vlan_trunks"] != "100,200" {
		t.Errorf("known values are read: %v", blocks[0].values)
	}

	if _, present := blocks[0].values["trunk_only"]; present {
		t.Errorf("an unknown value resolves to nothing: %v", blocks[0].values)
	}

	// An unknown block as a whole (wired to another resource) is nothing yet.
	model.ConfigOvsPortGroup = NewConfigOvsPortGroupValueUnknown()

	if blocks, _ = configuredBlocks(ctx, model); len(blocks) != 0 {
		t.Errorf("an unknown block must not be reported: %v", blocks)
	}
}

func TestUnitResolveConfigInputsRoutesTheBlockForTheType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ovs := typeSpecFor("ovsPortGroup")

	model := seedModel(t, ctx, nil)
	model.ConfigOvsPortGroup = ovsBlock(ctx, types.StringValue("100,200"), types.BoolValue(false))
	model.ConnectedGateway = types.StringValue("gw")

	inputs, diags := resolveConfigInputs(ctx, model, ovs)
	if diags.HasError() {
		t.Fatal(diags)
	}

	if inputs.all["hvmVlanTrunks"] != "100,200" || inputs.all["hvmTrunkOnly"] != false || inputs.all["connectedGateway"] != "gw" {
		t.Errorf("routing: %v", inputs.all)
	}

	if len(inputs.immutable) != 0 || inputs.typ != ovs || inputs.block == nil {
		t.Errorf("inputs: immutable=%v typ=%v block=%v", inputs.immutable, inputs.typ, inputs.block)
	}

	// Nothing configured: nothing to send, but digests are still computed so
	// a later plan compares equal.
	empty, diags := resolveConfigInputs(ctx, seedModel(t, ctx, nil), ovs)
	if diags.HasError() || empty.all != nil || empty.digests.All == "" {
		t.Errorf("empty inputs: %+v %v", empty, diags)
	}

	if empty.digests == inputs.digests {
		t.Error("different configurations must digest differently")
	}
}

func TestUnitCheckBlockMatchesType(t *testing.T) {
	t.Parallel()

	ovsBlock, _ := networkconfig.BlockFor("config_ovs_port_group")
	ovs := typeSpecFor("ovsPortGroup")
	azure := typeSpecFor("azureNetwork")

	if azure == nil || ovs == nil {
		t.Fatal("the catalogue must know ovsPortGroup and azureNetwork")
	}

	if d := checkBlockMatchesType(nil, azure); d != nil {
		t.Errorf("no block: no error, got %v", d)
	}

	if d := checkBlockMatchesType(&ovsBlock, nil); d != nil {
		t.Errorf("unknown type: nothing to check, got %v", d)
	}

	if d := checkBlockMatchesType(&ovsBlock, ovs); d != nil {
		t.Errorf("matching type: no error, got %v", d)
	}

	d := checkBlockMatchesType(&ovsBlock, azure)
	if d == nil {
		t.Fatal("a block for another type must be rejected")
	}

	detail := d.Detail()
	wants := []string{
		"config_ovs_port_group configures", "ovsPortGroup (HVM Standard Network)", "azureNetwork (", "Use config_azure instead",
	}
	for _, want := range wants {
		if !strings.Contains(detail, want) {
			t.Errorf("detail lacks %q: %s", want, detail)
		}
	}

	// A type with no block of its own is pointed at the free-form map.
	amazon := typeSpecFor("amazonSubnet")
	if d := checkBlockMatchesType(&ovsBlock, amazon); d == nil || !strings.Contains(d.Detail(), "free-form config map") {
		t.Errorf("a type without a block must be pointed at config: %v", d)
	}
}

func TestUnitImmutableAttributesExistInTheSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attrs := NetworkResourceSchema(ctx).Attributes

	// The attributes the catalogue may require replacement on must be schema
	// attributes, or ModifyPlan would add a RequiresReplace path Terraform
	// rejects. Fields the schema does not carry yet are fine: knownAndDifferent
	// never sees them differ.
	for _, code := range networkconfig.Codes() {
		typ, _ := networkconfig.Lookup(code)

		for _, a := range networkconfig.ImmutableAttributes(typ) {
			if _, ok := attrs[a]; !ok {
				switch a {
				case "scope_id", "vxlan_id", "parent_network_id", "availability_zone":
					// Added with the spec change that carries them.
				default:
					t.Errorf("%s: immutable attribute %q is not in the schema", code, a)
				}
			}
		}
	}
}

func TestUnitTypeSpecFallsBackToTheTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// No client (the appliance cannot be reached at plan time): the generated
	// table answers by code, and an id alone cannot be looked up.
	if spec := typeSpec(ctx, nil, 0, "ovsPortGroup"); spec == nil || spec.Code != "ovsPortGroup" {
		t.Errorf("table fallback by code: got %v", spec)
	}

	if spec := typeSpec(ctx, nil, 63, ""); spec != nil {
		t.Errorf("an id alone cannot be resolved without the appliance: got %v", spec)
	}

	if spec := typeSpec(ctx, nil, 0, "juniperSdn"); spec != nil {
		t.Errorf("a type in neither place is generic: got %v", spec)
	}

	// A resolved type whose live shape could not be read uses the table entry
	// for its code.
	if spec := liveOrTableSpec(networkType{id: 63, code: "ovsPortGroup"}); spec == nil || !spec.HasNetworkServer {
		t.Errorf("liveOrTableSpec without a live spec: got %v", spec)
	}

	live := networkconfig.TypeSpec{Code: "ovsPortGroup", Name: "live"}
	if spec := liveOrTableSpec(networkType{id: 63, code: "ovsPortGroup", spec: &live}); spec != &live {
		t.Error("liveOrTableSpec must prefer the live shape")
	}
}

func TestUnitCheckDeclared(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	block, _ := networkconfig.BlockFor("config_hvm_data")
	older := networkconfig.TypeSpec{
		Code: "hvmPortGroupDataNetwork", Name: "HVM Data Network", Live: true,
		Fields: []networkconfig.FieldSpec{
			{Name: "hvmHostNetworkInterface", Context: networkconfig.ContextConfig},
			{Name: "hvmMtu", Context: networkconfig.ContextConfig, Editable: true},
		},
	}

	inputs := configInputs{
		typ:         &older,
		block:       &block,
		blockValues: map[string]any{"cluster_network_layout": "HOST_INTERFACE", "host_network_interface": "bond0", "mtu": 1500},
	}

	// No client: the message simply has no version.
	diags := checkDeclared(ctx, nil, inputs)
	if !diags.HasError() || diags.ErrorsCount() != 1 {
		t.Fatalf("one undeclared block attribute must be one error, got %v", diags)
	}

	detail := diags.Errors()[0].Detail()
	for _, want := range []string{
		"hvmPortGroupDataNetwork (HVM Data Network)", "hvmClusterNetworkLayout",
		"config_hvm_data.cluster_network_layout", "ignore it", "newer Morpheus release",
	} {
		if !strings.Contains(detail, want) {
			t.Errorf("detail lacks %q: %s", want, detail)
		}
	}

	// From the generated table (not live): nothing is checked — it is another
	// release's view, and the appliance may well declare the field.
	table := older
	table.Live = false
	inputs.typ = &table

	if diags := checkDeclared(ctx, nil, inputs); len(diags) != 0 {
		t.Errorf("a table entry must not produce diagnostics: %v", diags)
	}

	// Free-form config with a key the type does not declare: a warning, and
	// the declared key passes.
	inputs.typ = &older
	inputs.block, inputs.blockValues = nil, nil
	inputs.generic = map[string]any{"hvmMtu": 9000, "hvmJumbo": true}

	diags = checkDeclared(ctx, nil, inputs)
	if diags.HasError() || diags.WarningsCount() != 1 || !strings.Contains(diags.Warnings()[0].Detail(), "hvmJumbo") {
		t.Errorf("an undeclared config key must be one warning naming it, got %v", diags)
	}

	// No type at all: nothing to compare against.
	inputs.typ = nil
	if diags := checkDeclared(ctx, nil, inputs); len(diags) != 0 {
		t.Errorf("no type, no diagnostics: %v", diags)
	}
}
