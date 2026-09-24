// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestUnitDigestIsCanonical(t *testing.T) {
	t.Parallel()

	a := map[string]any{"hvmVlanTrunks": "100,200", "hvmTrunkOnly": false, "nested": map[string]any{"b": 2, "a": []any{1, "x"}}}
	b := map[string]any{"nested": map[string]any{"a": []any{1, "x"}, "b": 2}, "hvmTrunkOnly": false, "hvmVlanTrunks": "100,200"}

	da, err := Digest(a)
	if err != nil {
		t.Fatal(err)
	}

	db, err := Digest(b)
	if err != nil {
		t.Fatal(err)
	}

	if da != db {
		t.Errorf("same content, different digests: %s vs %s", da, db)
	}

	if len(da) != 64 {
		t.Errorf("sha256 hex must be 64 chars, got %d", len(da))
	}

	c := map[string]any{"hvmVlanTrunks": "100,201", "hvmTrunkOnly": false, "nested": a["nested"]}

	dc, _ := Digest(c)
	if dc == da {
		t.Error("a changed value must change the digest")
	}

	empty, _ := Digest(map[string]any{})

	nilMap, _ := Digest(nil)
	if empty != nilMap {
		t.Error("nil and empty maps must digest equally: an unset block and an empty one are the same request")
	}
}

// fakePrivate is the framework's private state, as far as the digests use it.
type fakePrivate struct {
	data map[string][]byte
}

func (f *fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f.data[key], nil
}

func (f *fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if f.data == nil {
		f.data = map[string][]byte{}
	}

	f.data[key] = value

	return nil
}

func TestUnitDigestsRoundTripThroughPrivateState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := &fakePrivate{}

	if _, ok, _ := Load(ctx, private); ok {
		t.Fatal("empty private state must report no digests")
	}

	want, err := Compute(map[string]any{"hvmVlanTrunks": "100"}, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}

	if diags := Store(ctx, private, want); diags.HasError() {
		t.Fatal(diags)
	}

	got, ok, diags := Load(ctx, private)
	if diags.HasError() || !ok {
		t.Fatalf("expected stored digests, ok=%v diags=%v", ok, diags)
	}

	if got != want {
		t.Errorf("round trip changed the digests: %+v vs %+v", got, want)
	}

	// A digest written by another version is treated as missing, never compared.
	private.data[PrivateKey] = []byte(`{"all":"abc","immutable":"def","v":99}`)
	if _, ok, _ := Load(ctx, private); ok {
		t.Error("a different digest version must be treated as missing")
	}

	// Garbage is treated as missing too.
	private.data[PrivateKey] = []byte(`not json`)
	if _, ok, _ := Load(ctx, private); ok {
		t.Error("undecodable private state must be treated as missing")
	}

	// Nil private state: Load says missing; Store reports an error rather than
	// silently dropping the record.
	if _, ok, _ := Load(ctx, nil); ok {
		t.Error("nil private state must report no digests")
	}

	if diags := Store(ctx, nil, want); !diags.HasError() {
		t.Error("storing into nil private state must fail loudly")
	}
}

func TestUnitRouteOvsPortGroup(t *testing.T) {
	t.Parallel()

	b, ok := BlockFor("config_ovs_port_group")
	if !ok || !b.Serves("ovsPortGroup") {
		t.Fatalf("config_ovs_port_group must be declared for ovsPortGroup, got %+v ok=%v", b, ok)
	}

	typ, ok := Lookup("ovsPortGroup")
	if !ok {
		t.Fatal("the catalogue must know ovsPortGroup")
	}

	generic := map[string]any{"connectedGateway": "gw"}
	all, immutable := Route(&typ, &b, generic, map[string]any{"vlan_trunks": "100,200", "trunk_only": true})

	if all["hvmVlanTrunks"] != "100,200" || all["hvmTrunkOnly"] != true || all["connectedGateway"] != "gw" {
		t.Errorf("routing: got %v", all)
	}

	if len(immutable) != 0 {
		t.Errorf("both OVS config fields are editable; immutable subset must be empty, got %v", immutable)
	}

	// Attributes the block does not set are not sent, so the server keeps its
	// current value rather than receiving a null.
	all, _ = Route(&typ, &b, nil, map[string]any{"trunk_only": false})
	if _, present := all["hvmVlanTrunks"]; present {
		t.Errorf("an unset attribute must not appear in the request, got %v", all)
	}

	// No type and no block: the generic map passes through untouched.
	all, immutable = Route(nil, nil, map[string]any{"x": 1}, nil)
	if all["x"] != 1 || len(immutable) != 0 {
		t.Errorf("generic pass-through: got %v / %v", all, immutable)
	}
}

func TestUnitRouteImmutableSubset(t *testing.T) {
	t.Parallel()

	// azureNetwork declares its three config fields as not editable, so a
	// change to any of them must replace the network; the digest over the
	// immutable subset is what detects that.
	typ, ok := Lookup("azureNetwork")
	if !ok {
		t.Fatal("the catalogue must know azureNetwork")
	}

	for _, key := range []string{"resourceGroupId", "subnetName", "subnetCidr"} {
		if f, ok := typ.Field(key); !ok || f.Context != ContextConfig || f.Editable {
			t.Fatalf("azureNetwork.%s: expected a non-editable config field, got %+v ok=%v", key, f, ok)
		}
	}

	// Through the generic map, since azureNetwork has no typed block yet: the
	// type still decides what is immutable.
	all, immutable := Route(&typ, nil, map[string]any{"subnetCidr": "10.0.0.0/24", "custom": "x"}, nil)

	if len(all) != 2 || immutable["subnetCidr"] != "10.0.0.0/24" || len(immutable) != 1 {
		t.Errorf("immutable subset wrong: all=%v immutable=%v", all, immutable)
	}

	d1, _ := Compute(all, immutable)
	d2, _ := Compute(map[string]any{"subnetCidr": "10.0.0.0/24", "custom": "y"}, immutable)
	d3, _ := Compute(map[string]any{"subnetCidr": "10.0.1.0/24", "custom": "x"}, map[string]any{"subnetCidr": "10.0.1.0/24"})

	if d1.All == d2.All || d1.Immutable != d2.Immutable {
		t.Error("an editable change must alter All and leave Immutable alone")
	}

	if d1.Immutable == d3.Immutable {
		t.Error("an immutable change must alter Immutable")
	}

	// Type unknown at plan time: nothing is immutable, so no replacement is
	// planned on a guess.
	_, immutable = Route(nil, nil, map[string]any{"subnetCidr": "10.0.0.0/24"}, nil)
	if len(immutable) != 0 {
		t.Errorf("without a type nothing may be immutable, got %v", immutable)
	}
}

func TestUnitBlocksMatchTheCatalogue(t *testing.T) {
	t.Parallel()

	if problems := Check(); len(problems) != 0 {
		t.Errorf("the typed blocks disagree with the generated catalogue:\n  %s", strings.Join(problems, "\n  "))
	}

	for _, b := range Blocks {
		if !strings.HasPrefix(b.Attr, "config_") {
			t.Errorf("block %q must be a config_* attribute", b.Attr)
		}

		if got, ok := BlockForType(b.Types[0]); !ok || got.Attr != b.Attr {
			t.Errorf("BlockForType(%q) = %q, want %q", b.Types[0], got.Attr, b.Attr)
		}
	}
}

func TestUnitCheckReportsInconsistencies(t *testing.T) {
	t.Parallel()

	ovs, _ := Lookup("ovsPortGroup")

	bad := Block{Attr: "config_bad", Types: []string{"ovsPortGroup"}, Fields: []BlockField{
		{Attr: "vlan_trunks", Key: "hvmVlanTrunks"},
		{Attr: "switch", Key: "switchId"},   // network-context: belongs at the top level
		{Attr: "nothing", Key: "noSuchKey"}, // not an option type
		// hvmTrunkOnly neither mapped nor omitted
	}, Omit: []string{"alsoMissing"}}

	problems := checkBlockAgainstType(bad, ovs)

	want := []string{
		`config_bad.switch: "switchId" is a network-context field of ovsPortGroup and belongs at the top level`,
		`config_bad.nothing: "noSuchKey" is not an option type of ovsPortGroup`,
		`config_bad: omitted "alsoMissing" is not an option type of ovsPortGroup`,
		`config_bad: config field "hvmTrunkOnly" of ovsPortGroup has no attribute and is not omitted`,
	}

	if len(problems) != len(want) {
		t.Fatalf("got %d problems, want %d:\n  %s", len(problems), len(want), strings.Join(problems, "\n  "))
	}

	for i := range want {
		if problems[i] != want[i] {
			t.Errorf("problem %d:\n  got  %s\n  want %s", i, problems[i], want[i])
		}
	}
}

func TestUnitCatalogueShape(t *testing.T) {
	t.Parallel()

	if KnownTypesRelease == "" || len(knownTypes) < 60 {
		t.Fatalf("catalogue looks empty: release %q, %d types", KnownTypesRelease, len(knownTypes))
	}

	for code, typ := range knownTypes {
		if typ.Code != code {
			t.Errorf("catalogue key %q does not match Code %q", code, typ.Code)
		}

		for _, f := range typ.Fields {
			if f.Context != ContextNetwork && f.Context != ContextConfig {
				t.Errorf("%s.%s: unexpected context %q", code, f.Name, f.Context)
			}
		}
	}

	// Facts the resource relies on, pinned so a re-harvest that changes them
	// is noticed here rather than in an acceptance test.
	ovs, _ := Lookup("ovsPortGroup")
	if !ovs.HasNetworkServer || ovs.Flags.Name || !ovs.Flags.VlanID {
		t.Errorf("ovsPortGroup facts changed: %+v", ovs)
	}

	if sw, _ := ovs.Field("switchId"); sw.Editable || sw.Context != ContextNetwork {
		t.Errorf("ovsPortGroup.switchId must be a non-editable network field, got %+v", sw)
	}
}

func TestUnitImmutableAttributes(t *testing.T) {
	t.Parallel()

	cases := map[string][]string{
		// nameEditable off; switchId is an option type off the bind list and not
		// editable; vlanId follows the type's vlanIdEditable flag, which is on;
		// zonePool.id is bound on every update whatever its option type says.
		"ovsPortGroup": {"name", "switch_id"},
		// availabilityZone is bound on create only; zonePool.id always.
		"amazonSubnet": {"availability_zone"},
	}

	for code, want := range cases {
		typ, ok := Lookup(code)
		if !ok {
			t.Fatalf("the catalogue must know %s", code)
		}

		// The flag-bound core fields join the list when their flag is off, for
		// every type; the fixed expectations above cover the rest.
		if !typ.Flags.Cidr {
			want = append(want, "cidr", "cidr_ipv6", "ipv4enabled", "ipv6enabled")
		}

		if !typ.Flags.Gateway {
			want = append(want, "gateway", "gateway_ipv6")
		}

		if !typ.Flags.DhcpServer {
			want = append(want, "dhcp_server", "dhcp_server_ipv6")
		}

		if !typ.Flags.DNS {
			want = append(want, "dns_primary", "dns_secondary", "dns_primary_ipv6", "dns_secondary_ipv6")
		}

		if !typ.Flags.Name && code != "ovsPortGroup" {
			want = append(want, "name")
		}

		if !typ.Flags.VlanID {
			want = append(want, "vlan_id")
		}

		sort.Strings(want)

		if got := ImmutableAttributes(typ); strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("%s: immutable attributes\n  got  %v\n  want %v (flags %+v)", code, got, want, typ.Flags)
		}
	}

	// A synthetic type shows the three classes side by side.
	synthetic := TypeSpec{
		Code:  "synthetic",
		Flags: TypeFlags{Name: true, Cidr: true, Gateway: true, VlanID: false, DhcpServer: true, DNS: true},
		Fields: []FieldSpec{
			{Name: "zonePool.id", Context: ContextNetwork, Editable: false}, // always bound: editable
			{Name: "vlanId", Context: ContextNetwork, Editable: true},       // flag off wins: immutable
			{Name: "scopeId", Context: ContextNetwork, Editable: false},     // off the list, not editable
			{Name: "vxlanId", Context: ContextNetwork, Editable: true},      // off the list, editable
			{Name: "availabilityZone", Context: ContextNetwork, Editable: true},
			{Name: "streamType", Context: ContextNetwork, Editable: false}, // no attribute: not reported
			{Name: "subnetCidr", Context: ContextConfig, Editable: false},  // config context: not top-level
		},
	}

	if got := ImmutableAttributes(synthetic); strings.Join(got, ",") != "availability_zone,scope_id,vlan_id" {
		t.Errorf("synthetic: got %v", got)
	}
}

func TestUnitEveryNetworkFieldIsMappedOrKnownUnsupported(t *testing.T) {
	t.Parallel()

	for code, typ := range knownTypes {
		for _, f := range typ.NetworkFields() {
			_, mapped := attributeFor[f.Name]
			_, unsupported := unsupportedNetworkFields[f.Name]

			switch {
			case mapped && unsupported:
				t.Errorf("%s.%s: both mapped to an attribute and listed as unsupported", code, f.Name)
			case !mapped && !unsupported:
				t.Errorf("%s.%s: new network-context field; map it in attributeFor or record why it is unsupported", code, f.Name)
			}
		}
	}

	for field := range attributeFor {
		if alwaysBound[field] && createOnly[field] {
			t.Errorf("%s: cannot be both always bound and create-only", field)
		}
	}

	if attr, ok := NetworkAttribute("switchId"); !ok || attr != "switch_id" {
		t.Errorf("NetworkAttribute(switchId) = %q, %v", attr, ok)
	}
}

func TestUnitUndeclaredFieldsAndKeys(t *testing.T) {
	t.Parallel()

	// An appliance whose data-network type predates the virtual-switch layout:
	// it declares the interface and the MTU, not the layout or the segment.
	older := TypeSpec{Code: "hvmPortGroupDataNetwork", Live: true, Fields: []FieldSpec{
		{Name: "zonePool.id", Context: ContextNetwork},
		{Name: "vlanId", Context: ContextNetwork},
		{Name: "hvmHostNetworkInterface", Context: ContextConfig, Editable: false},
		{Name: "hvmMtu", Context: ContextConfig, Editable: true},
	}}

	block, ok := BlockFor("config_hvm_data")
	if !ok {
		t.Fatal("config_hvm_data must be declared")
	}

	// Only the attributes that are set count.
	got := block.Undeclared(older, map[string]any{
		"cluster_network_layout": "HOST_INTERFACE",
		"host_network_interface": "bond0",
		"mtu":                    1500,
	})

	if len(got) != 1 || got[0].Attr != "cluster_network_layout" || got[0].Key != "hvmClusterNetworkLayout" {
		t.Errorf("undeclared: got %+v", got)
	}

	if got := block.Undeclared(older, map[string]any{"host_network_interface": "bond0"}); len(got) != 0 {
		t.Errorf("declared fields must pass: %+v", got)
	}

	// A network-context field is not a config key either.
	if got := block.Undeclared(TypeSpec{Fields: []FieldSpec{{Name: "hvmMtu", Context: ContextNetwork}}},
		map[string]any{"mtu": 1500}); len(got) != 1 {
		t.Errorf("a key declared in the network context is not a config field: %+v", got)
	}

	// The current catalogue declares everything the block maps.
	current, _ := Lookup("hvmPortGroupDataNetwork")
	if got := block.Undeclared(current, map[string]any{
		"cluster_network_layout": "VSWITCH", "virtual_switch_segment_id": 7, "host_network_interface": "bond0", "mtu": 9000,
	}); len(got) != 0 {
		t.Errorf("the harvested type must declare every block field: %+v", got)
	}

	keys := UndeclaredKeys(older, map[string]any{"hvmMtu": 1500, "zzz": 1, "aaa": 2, "vlanId": 3})
	if strings.Join(keys, ",") != "aaa,vlanId,zzz" {
		t.Errorf("undeclared keys: got %v", keys)
	}

	if keys := UndeclaredKeys(older, nil); len(keys) != 0 {
		t.Errorf("no config, no keys: %v", keys)
	}
}
