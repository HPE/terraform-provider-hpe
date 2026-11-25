// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

// baseNsxPlan is a minimal generic (non-SDN) plan.
func baseNsxPlan() NetworkServerModel {
	return NetworkServerModel{
		Name:             types.StringValue("test"),
		TypeCode:         types.StringValue("nsx-t"),
		ServiceUrl:       types.StringValue("https://nsxt.example.com"),
		ServiceUsername:  types.StringValue("admin"),
		Config:           types.DynamicNull(),
		ConfigJuniperSdn: NewConfigJuniperSdnValueNull(),
		CloudId:          types.Int64Value(7),
		CredentialId:     types.Int64Null(),
		ServiceHost:      types.StringNull(),
		ServiceMode:      types.StringNull(),
		ServicePort:      types.Int64Null(),
		ServiceTokenWo:   types.StringNull(),
		NetworkFilter:    types.StringNull(),
		TenantMatch:      types.StringNull(),
		Visibility:       types.StringNull(),
		TenantIds:        types.SetNull(types.Int64Type),
	}
}

func marshalCreate(t *testing.T, plan *NetworkServerModel, pw types.String) map[string]any {
	t.Helper()

	return marshalCreateWithToken(t, plan, pw, types.StringNull())
}

func marshalCreateWithToken(t *testing.T, plan *NetworkServerModel, pw, token types.String) map[string]any {
	t.Helper()

	req, diags := buildCreateRequest(context.Background(), plan, pw, token)
	if diags.HasError() {
		t.Fatalf("buildCreateRequest: %s", diags.Errors())
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %s", err)
	}

	return out
}

// The generic variant sends the common bind fields, and the extra fields the
// struct lacks (serviceHost, servicePort, serviceToken, ...) plus the config
// map through AdditionalProperties, flattened into the body by the SDK.
func TestUnitBuildCreateGenericBody(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.ServiceHost = types.StringValue("nsxt.example.com")
	plan.ServicePort = types.Int64Value(443)
	plan.ServiceMode = types.StringValue("standard")
	plan.NetworkFilter = types.StringValue("10.0.0.0/8")
	plan.TenantMatch = types.StringValue(".*")
	plan.CloudId = types.Int64Value(7)
	plan.TenantIds = types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(3)})

	config := types.DynamicValue(types.ObjectValueMust(
		map[string]attr.Type{"tier1Gateway": types.StringType},
		map[string]attr.Value{"tier1Gateway": types.StringValue("t1")},
	))
	plan.Config = config

	body := marshalCreateWithToken(t, &plan, types.StringValue("secret"), types.StringValue("tok"))

	checks := map[string]any{
		"type":            "nsx-t",
		"name":            "test",
		"serviceUrl":      "https://nsxt.example.com",
		"serviceUsername": "admin",
		"servicePassword": "secret",
		"serviceHost":     "nsxt.example.com",
		"servicePort":     float64(443),
		"serviceMode":     "standard",
		"networkFilter":   "10.0.0.0/8",
		"tenantMatch":     ".*",
		"serviceToken":    "tok",
		"zoneId":          float64(7),
	}
	for k, want := range checks {
		if got := body[k]; got != want {
			t.Errorf("body[%q] = %v, want %v", k, got, want)
		}
	}

	cfg, ok := body["config"].(map[string]any)
	if !ok || cfg["tier1Gateway"] != "t1" {
		t.Errorf("config = %v, want map with tier1Gateway=t1", body["config"])
	}

	tenants, ok := body["tenants"].([]any)
	if !ok || len(tenants) != 1 {
		t.Fatalf("tenants = %v, want one entry", body["tenants"])
	}
	if tenants[0].(map[string]any)["id"] != float64(3) {
		t.Errorf("tenant id = %v, want 3", tenants[0])
	}
}

// credential_id set sends {credential:{id}} and omits it when unset.
func TestUnitBuildCreateCredential(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.CredentialId = types.Int64Value(42)

	body := marshalCreate(t, &plan, types.StringNull())

	cred, ok := body["credential"].(map[string]any)
	if !ok || cred["id"] != float64(42) {
		t.Errorf("credential = %v, want {id:42}", body["credential"])
	}

	plan.CredentialId = types.Int64Null()
	body = marshalCreate(t, &plan, types.StringNull())
	if _, ok := body["credential"]; ok {
		t.Errorf("credential should be omitted when credential_id is unset, got %v", body["credential"])
	}
}

// The SDN type uses the HPEMorpheusSDNNetworkServer variant with the typed
// config block.
func TestUnitBuildCreateSdnBody(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.TypeCode = types.StringValue(typeCodeJuniperSdn)
	plan.ServiceUrl = types.StringNull()
	plan.ServiceUsername = types.StringNull()
	plan.CloudId = types.Int64Value(9)
	plan.ConfigJuniperSdn = NewConfigJuniperSdnValueMust(
		ConfigJuniperSdnValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"controller_address": types.StringValue("10.0.0.10"),
			"show_introspect":    types.BoolValue(true),
			"enable_sub_tenants": types.BoolValue(false),
		},
	)

	body := marshalCreate(t, &plan, types.StringNull())

	if body["type"] != typeCodeJuniperSdn {
		t.Errorf("type = %v, want %s", body["type"], typeCodeJuniperSdn)
	}
	if body["zoneId"] != float64(9) {
		t.Errorf("zoneId = %v, want 9", body["zoneId"])
	}

	cfg, ok := body["config"].(map[string]any)
	if !ok {
		t.Fatalf("config = %v, want object", body["config"])
	}
	if cfg["controllerAddress"] != "10.0.0.10" {
		t.Errorf("controllerAddress = %v, want 10.0.0.10", cfg["controllerAddress"])
	}
	if cfg["showIntrospect"] != true {
		t.Errorf("showIntrospect = %v, want true", cfg["showIntrospect"])
	}
}

// A type_code that conflicts with the config_juniper_sdn block is rejected.
func TestUnitBuildCreateSdnTypeMismatch(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.TypeCode = types.StringValue("nsx-t")
	plan.ConfigJuniperSdn = NewConfigJuniperSdnValueMust(
		ConfigJuniperSdnValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"controller_address": types.StringValue("10.0.0.10"),
			"show_introspect":    types.BoolValue(false),
			"enable_sub_tenants": types.BoolValue(false),
		},
	)

	_, diags := buildCreateRequest(context.Background(), &plan, types.StringNull(), types.StringNull())
	if !diags.HasError() {
		t.Fatal("expected a diagnostic for the type/config mismatch")
	}
}

// On update, removing credential_id from a plan that previously had it set
// sends {credential:{type:local}} so the server drops the stored link.
func TestUnitBuildUpdateCredentialLocalSwitch(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.CredentialId = types.Int64Null()
	plan.ServiceUsername = types.StringValue("admin")

	prior := baseNsxPlan()
	prior.CredentialId = types.Int64Value(42)

	req, diags := buildUpdateRequest(
		context.Background(), &plan, &prior, types.StringNull(), types.StringNull(), clearable.Pending{}, false,
	)
	if diags.HasError() {
		t.Fatalf("buildUpdateRequest: %s", diags.Errors())
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}

	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal: %s", err)
	}

	cred, ok := body["credential"].(map[string]any)
	if !ok || cred["type"] != "local" {
		t.Errorf("credential = %v, want {type:local}", body["credential"])
	}
}

// servicePassword is sent on update only when the write-only digest changed.
func TestUnitBuildUpdatePasswordOnlyOnResend(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	prior := baseNsxPlan()

	body := marshalUpdate(t, &plan, &prior, types.StringValue("secret"), false)
	if _, ok := body["servicePassword"]; ok {
		t.Errorf("servicePassword should be omitted when not resending, got %v", body["servicePassword"])
	}

	body = marshalUpdate(t, &plan, &prior, types.StringValue("secret"), true)
	if body["servicePassword"] != "secret" {
		t.Errorf("servicePassword = %v, want secret", body["servicePassword"])
	}
}

func marshalUpdate(
	t *testing.T,
	plan, prior *NetworkServerModel,
	pw types.String,
	resend bool,
) map[string]any {
	t.Helper()

	req, diags := buildUpdateRequest(context.Background(), plan, prior, pw, types.StringNull(), clearable.Pending{}, resend)
	if diags.HasError() {
		t.Fatalf("buildUpdateRequest: %s", diags.Errors())
	}

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %s", err)
	}

	return out
}

// configBool parses the truthy forms Morpheus may store an option-type
// checkbox as.
func TestUnitConfigBoolTolerantParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value any
		want  bool
	}{
		"real bool true":  {true, true},
		"real bool false": {false, false},
		"string on":       {"on", true},
		"string off":      {"off", false},
		"string true":     {"true", true},
		"string false":    {"false", false},
		"missing":         {nil, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := map[string]interface{}{}
			if tc.value != nil {
				m["showIntrospect"] = tc.value
			}

			if got := configBool(m, "showIntrospect"); got != tc.want {
				t.Errorf("configBool(%v) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

// A generic (non-SDN) server must name its cloud: the API requires zoneId and
// the model would otherwise send 0.
func TestUnitBuildCreateRequiresCloudForGenericTypes(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.CloudId = types.Int64Null()

	_, diags := buildGenericCreate(context.Background(), &plan, types.StringNull(), types.StringNull())
	if !diags.HasError() {
		t.Fatal("expected an error for a generic server without cloud_id")
	}
	if !strings.Contains(diags.Errors()[0].Summary(), "cloud_id is required") {
		t.Errorf("unexpected diagnostic: %v", diags.Errors()[0])
	}

	plan.CloudId = types.Int64Value(7)
	body := marshalCreate(t, &plan, types.StringNull())
	if body["zoneId"] != float64(7) {
		t.Errorf("zoneId = %v, want 7", body["zoneId"])
	}
}

// TestUnitBuildSdnUpdateSendsCloud pins that the SDN update carries cloud_id
// as zoneId — set, and null when cleared — since the typed variant has no
// field for it and the server binds it for every type.
func TestUnitBuildSdnUpdateSendsCloud(t *testing.T) {
	t.Parallel()

	plan := baseNsxPlan()
	plan.TypeCode = types.StringValue(typeCodeJuniperSdn)
	plan.CloudId = types.Int64Value(9)
	plan.ConfigJuniperSdn = NewConfigJuniperSdnValueMust(
		ConfigJuniperSdnValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"controller_address": types.StringValue("10.0.0.10"),
			"show_introspect":    types.BoolValue(true),
			"enable_sub_tenants": types.BoolValue(false),
		},
	)

	upd := buildSdnUpdate(&plan, clearable.Pending{})

	raw, err := json.Marshal(upd)
	if err != nil {
		t.Fatal(err)
	}

	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}

	if got, ok := body["zoneId"].(float64); !ok || got != 9 {
		t.Errorf("zoneId = %v, want 9", body["zoneId"])
	}

	if _, ok := body["tenants"]; ok && body["tenants"] == nil {
		t.Error("tenants must not be sent as null")
	}

	plan.CloudId = types.Int64Null()

	upd = buildSdnUpdate(&plan, clearable.Pending{"cloud_id": clearable.Attribute{Name: "cloud_id", Mode: clearable.Clear}})

	raw, err = json.Marshal(upd)
	if err != nil {
		t.Fatal(err)
	}

	body = map[string]any{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}

	if v, ok := body["zoneId"]; !ok || v != nil {
		t.Errorf("cleared cloud_id: zoneId = %v (present %v), want an explicit null", v, ok)
	}
}
