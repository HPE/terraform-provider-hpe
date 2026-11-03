// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitHVMConfigReadsUserDataForImport: user_data forces replacement, so
// import must recover it from the API's untyped config keys; a missing or
// empty value reads as null, never "".
func TestUnitHVMConfigReadsUserDataForImport(t *testing.T) {
	t.Parallel()

	pool := "pool-1"
	poolID := &sdk.GetInstance200ResponseInstanceConfigResourcePoolId{String: &pool}
	apiConfig := &apiConfigType{
		ResourcePoolId:       poolID,
		AdditionalProperties: map[string]interface{}{"userData": "#cloud-config\npackages: [jq]"},
	}
	got, diags := getInstanceHVMConfig(context.Background(), 1, apiConfig)
	if diags.HasError() {
		t.Fatal(diags)
	}
	if got.UserData.ValueString() != "#cloud-config\npackages: [jq]" {
		t.Errorf("user_data: %v", got.UserData)
	}

	for name, props := range map[string]map[string]interface{}{
		"absent": {},
		"empty":  {"userData": ""},
		"nonstr": {"userData": 42},
	} {
		got, diags := getInstanceHVMConfig(context.Background(), 1, &apiConfigType{ResourcePoolId: poolID, AdditionalProperties: props})
		if diags.HasError() {
			t.Fatal(name, diags)
		}
		if !got.UserData.IsNull() {
			t.Errorf("%s: user_data must be null, got %v", name, got.UserData)
		}
	}
}

// TestUnitHVMConfigSendsUserData: the typed HVM configuration carries the
// value as config.userData, the key the KVM provisioning service reads.
func TestUnitHVMConfigSendsUserData(t *testing.T) {
	t.Parallel()

	ud := "#cloud-config\nbootcmd: [true]"
	cfg := sdk.HVMInstanceConfiguration{UserData: *sdk.NewNullableString(&ud)}
	b, err := cfg.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"userData":"#cloud-config\nbootcmd: [true]"}` {
		t.Errorf("marshalled: %s", b)
	}
}
