// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ptrStr(s string) *string { return &s }

// knownConfigBmaas builds a known ConfigBmaasValue for tests.
func knownConfigBmaas(t *testing.T) ConfigBmaasValue {
	t.Helper()
	v, diags := NewConfigBmaasValue(
		ConfigBmaasValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"ilo_ip_address":   types.StringValue("10.1.2.3"),
			"ilo_mac_address":  types.StringValue("aa:bb:cc:dd:ee:ff"),
			"ilo_password_wo":  types.StringNull(),
			"ilo_username":     types.StringValue("admin"),
			"pre_provisioned":  types.BoolValue(true),
			"resource_pool_id": types.StringValue("pool-1"),
		},
	)
	if diags.HasError() {
		t.Fatalf("building config_bmaas: %v", diags)
	}

	return v
}

func knownConfigOmega(t *testing.T) ConfigOmegaBmaasValue {
	t.Helper()
	v, diags := NewConfigOmegaBmaasValue(
		ConfigOmegaBmaasValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"console_host":        types.StringValue("10.0.0.1"),
			"console_password_wo": types.StringNull(),
			"console_username":    types.StringValue("admin"),
			"num_nics":            types.Int64Value(1),
			"pre_provisioned":     types.BoolNull(),
			"resource_pool_id":    types.StringValue("pool-9"),
		},
	)
	if diags.HasError() {
		t.Fatalf("building config_omega_bmaas: %v", diags)
	}

	return v
}

// baseModel returns a HostModel with the common required fields set and all the
// config variants null.
func baseModel() HostModel {
	return HostModel{
		Name:              types.StringValue("host-1"),
		CloudId:           types.Int64Value(15174),
		GroupId:           types.Int64Value(2),
		ComputeServerType: types.StringValue("omega.baremetal.stub-server"),
		Config:            types.DynamicNull(),
		ConfigBmaas:       NewConfigBmaasValueNull(),
		ConfigOmegaBmaas:  NewConfigOmegaBmaasValueNull(),
		Labels:            types.SetNull(types.StringType),
		NetworkInterfaces: types.ListNull(NetworkInterfacesValue{}.Type(context.Background())),
		Volumes:           types.ListNull(VolumesValue{}.Type(context.Background())),
		Description:       types.StringNull(),
		Visibility:        types.StringNull(),
		SshUsername:       types.StringNull(),
		PlanId:            types.Int64Null(),
		UserData:          types.StringNull(),
	}
}

// marshalServer marshals the request's server.config to a key map for assertion.
func configKeys(t *testing.T, req any) map[string]any {
	t.Helper()
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return m
}

func TestUnitBuildConfigBmaasRequestShape(t *testing.T) {
	t.Parallel()

	m := baseModel()
	m.ComputeServerType = types.StringValue("hpe-baremetal")
	m.ConfigBmaas = knownConfigBmaas(t)

	cfg, diags := buildConfig(context.Background(), m, ptrStr("secret"), nil)
	if diags.HasError() {
		t.Fatalf("buildConfig: %v", diags)
	}
	if cfg == nil || cfg.BmaasHostConfigObject == nil {
		t.Fatalf("expected BmaasHostConfigObject variant")
	}

	keys := configKeys(t, cfg)
	for _, want := range []string{"iloIpAddress", "iloUsername", "iloPassword", "iloMacAddress", "resourcePoolId", "preProvisioned"} {
		if _, ok := keys[want]; !ok {
			t.Errorf("missing key %q in %v", want, keys)
		}
	}
	if keys["iloPassword"] != "secret" {
		t.Errorf("iloPassword = %v, want secret", keys["iloPassword"])
	}
	if keys["iloIpAddress"] != "10.1.2.3" {
		t.Errorf("iloIpAddress = %v", keys["iloIpAddress"])
	}
}

func TestUnitBuildConfigOmegaRequestShape(t *testing.T) {
	t.Parallel()

	m := baseModel()
	m.ConfigOmegaBmaas = knownConfigOmega(t)

	cfg, diags := buildConfig(context.Background(), m, nil, ptrStr("x"))
	if diags.HasError() {
		t.Fatalf("buildConfig: %v", diags)
	}
	if cfg == nil || cfg.OmegaBmaasHostConfigObject == nil {
		t.Fatalf("expected OmegaBmaasHostConfigObject variant")
	}

	keys := configKeys(t, cfg)
	for _, want := range []string{"consoleHost", "consoleUsername", "consolePassword", "numNics", "resourcePoolId"} {
		if _, ok := keys[want]; !ok {
			t.Errorf("missing key %q in %v", want, keys)
		}
	}
	if keys["consolePassword"] != "x" {
		t.Errorf("consolePassword = %v", keys["consolePassword"])
	}
	if keys["consoleHost"] != "10.0.0.1" {
		t.Errorf("consoleHost = %v", keys["consoleHost"])
	}
}

func TestUnitBuildConfigGenericPassthrough(t *testing.T) {
	t.Parallel()

	m := baseModel()
	m.Config = types.DynamicValue(types.ObjectValueMust(
		map[string]attr.Type{
			"customKey": types.StringType,
			"userData":  types.StringType,
		},
		map[string]attr.Value{
			"customKey": types.StringValue("v"),
			"userData":  types.StringValue("#cloud-config"),
		},
	))

	cfg, diags := buildConfig(context.Background(), m, nil, nil)
	if diags.HasError() {
		t.Fatalf("buildConfig: %v", diags)
	}
	if cfg == nil || cfg.HostConfigGeneric == nil {
		t.Fatalf("expected HostConfigGeneric variant")
	}

	keys := configKeys(t, cfg)
	if keys["customKey"] != "v" {
		t.Errorf("customKey = %v, want v", keys["customKey"])
	}
	if keys["userData"] != "#cloud-config" {
		t.Errorf("userData = %v", keys["userData"])
	}
}

func TestUnitClassifyStatus(t *testing.T) {
	t.Parallel()

	cases := map[string]statusClass{
		"provisioning": statusProvisioning,
		"":             statusProvisioning,
		"queued":       statusProvisioning,
		"provisioned":  statusReady,
		"running":      statusReady,
		"ok":           statusReady,
		"available":    statusReady, // a registered bare-metal node
		"reserved":     statusReady,
		"resizing":     statusProvisioning,
		"failed":       statusFailed,
		"error":        statusFailed,
		"teardown":     statusFailed,
		"weird":        statusUnknown,
	}
	for status, want := range cases {
		if got := classifyStatus(status); got != want {
			t.Errorf("classifyStatus(%q) = %d, want %d", status, got, want)
		}
	}
}

// The keys every host type shares -- userData and provisionKey (the SSH key
// pair) -- ride in whichever config variant is configured, and alone when
// none is.
func TestUnitBuildConfigCommonKeysInEveryVariant(t *testing.T) {
	t.Parallel()

	variants := map[string]func(*HostModel){
		"bmaas": func(m *HostModel) { m.ConfigBmaas = knownConfigBmaas(t) },
		"omega": func(m *HostModel) { m.ConfigOmegaBmaas = knownConfigOmega(t) },
		"generic": func(m *HostModel) {
			m.Config = types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"customKey": types.StringType},
				map[string]attr.Value{"customKey": types.StringValue("v")},
			))
		},
		"none": func(*HostModel) {},
	}

	for name, set := range variants {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := baseModel()
			m.UserData = types.StringValue("#cloud-config")
			m.SshKeyPairId = types.Int64Value(42)
			set(&m)

			cfg, diags := buildConfig(context.Background(), m, nil, nil)
			if diags.HasError() {
				t.Fatalf("buildConfig: %v", diags)
			}
			if cfg == nil {
				t.Fatal("expected a config")
			}

			keys := configKeys(t, cfg)
			if keys["userData"] != "#cloud-config" {
				t.Errorf("userData = %v", keys["userData"])
			}
			if keys["provisionKey"] != float64(42) {
				t.Errorf("provisionKey = %v, want 42", keys["provisionKey"])
			}
		})
	}
}

// attrToGo keeps structure: nested objects become maps, lists become slices
// and nulls become nil, rather than their textual rendering.
func TestUnitAttrToGoRecurses(t *testing.T) {
	t.Parallel()

	v := types.ObjectValueMust(
		map[string]attr.Type{
			"nested": types.ObjectType{AttrTypes: map[string]attr.Type{"k": types.StringType}},
			"list":   types.ListType{ElemType: types.Int64Type},
			"gone":   types.StringType,
			"flag":   types.BoolType,
		},
		map[string]attr.Value{
			"nested": types.ObjectValueMust(map[string]attr.Type{"k": types.StringType}, map[string]attr.Value{"k": types.StringValue("v")}),
			"list":   types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(1), types.Int64Value(2)}),
			"gone":   types.StringNull(),
			"flag":   types.BoolValue(true),
		},
	)

	got, diags := attrToGo(context.Background(), v)
	if diags.HasError() {
		t.Fatalf("attrToGo: %v", diags)
	}

	b, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	const want = `{"flag":true,"gone":null,"list":[1,2],"nested":{"k":"v"}}`
	if string(b) != want {
		t.Errorf("got %s, want %s", b, want)
	}
}
