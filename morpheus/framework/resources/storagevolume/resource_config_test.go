// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestBuildCreateConfigAlletraMPBMaaS exercises the typed Alletra MP BMaaS
// branch of buildCreateConfig, isolating the formatting logic (bool -> on/off,
// id -> "[id: N]", instance list expansion) without a full acceptance round-trip.
func TestBuildCreateConfigAlletraMPBMaaS(t *testing.T) {
	ctx := context.Background()

	model := &StorageVolumeModel{
		ConfigAlletrampBmaas: ConfigAlletrampBmaasValue{
			DatastoreId:     types.Int64Value(5),
			Shared:          types.BoolValue(true),
			ComputeServerId: types.Int64Value(10),
			InstanceIds: types.ListValueMust(types.Int64Type, []attr.Value{
				types.Int64Value(7),
				types.Int64Value(8),
			}),
			RemoteCopyTargetId:   types.StringValue("rct-1"),
			UseExistingVolumeSet: types.BoolValue(false),
			VolumeSetId:          types.StringValue("vs-7"),
			VolumeSetName:        types.StringValue("myset"),
			state:                attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	cfg := buildCreateConfig(ctx, model, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if cfg == nil || cfg.AlletraMPBMaaSVolumeConfiguration == nil {
		t.Fatalf("expected AlletraMPBMaaSVolumeConfiguration, got %+v", cfg)
	}

	v := cfg.AlletraMPBMaaSVolumeConfiguration
	if v.HpeStorageDatastore != 5 {
		t.Errorf("HpeStorageDatastore = %d, want 5", v.HpeStorageDatastore)
	}
	if got := derefString(v.HpeStorageVolumeShared); got != "on" {
		t.Errorf("HpeStorageVolumeShared = %q, want \"on\"", got)
	}
	if got := derefString(v.HpeStorageComputeServer); got != "[id: 10]" {
		t.Errorf("HpeStorageComputeServer = %q, want \"[id: 10]\"", got)
	}
	wantInstances := []string{"[id: 7]", "[id: 8]"}
	if len(v.HpeStorageInstances) != len(wantInstances) {
		t.Fatalf("HpeStorageInstances = %v, want %v", v.HpeStorageInstances, wantInstances)
	}
	for i, want := range wantInstances {
		if v.HpeStorageInstances[i] != want {
			t.Errorf("HpeStorageInstances[%d] = %q, want %q", i, v.HpeStorageInstances[i], want)
		}
	}
	if got := derefString(v.HpeStorageExistingVolumeSet); got != "off" {
		t.Errorf("HpeStorageExistingVolumeSet = %q, want \"off\"", got)
	}
	if got := derefString(v.HpeStorageRemotecopytargetId); got != "rct-1" {
		t.Errorf("HpeStorageRemotecopytargetId = %q, want \"rct-1\"", got)
	}
	if got := derefString(v.HpeStorageVolumesetId); got != "vs-7" {
		t.Errorf("HpeStorageVolumesetId = %q, want \"vs-7\"", got)
	}
	if got := derefString(v.HpeStorageVolumeSetName); got != "myset" {
		t.Errorf("HpeStorageVolumeSetName = %q, want \"myset\"", got)
	}
}

// TestBuildCreateConfigAlletraMPBMaaSMinimal confirms that unset optional fields
// are omitted (left nil) while the required datastore id is still sent.
func TestBuildCreateConfigAlletraMPBMaaSMinimal(t *testing.T) {
	ctx := context.Background()

	model := &StorageVolumeModel{
		ConfigAlletrampBmaas: ConfigAlletrampBmaasValue{
			DatastoreId:          types.Int64Value(42),
			Shared:               types.BoolNull(),
			ComputeServerId:      types.Int64Null(),
			InstanceIds:          types.ListNull(types.Int64Type),
			RemoteCopyTargetId:   types.StringNull(),
			UseExistingVolumeSet: types.BoolNull(),
			VolumeSetId:          types.StringNull(),
			VolumeSetName:        types.StringNull(),
			state:                attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	cfg := buildCreateConfig(ctx, model, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if cfg == nil || cfg.AlletraMPBMaaSVolumeConfiguration == nil {
		t.Fatalf("expected AlletraMPBMaaSVolumeConfiguration, got %+v", cfg)
	}

	v := cfg.AlletraMPBMaaSVolumeConfiguration
	if v.HpeStorageDatastore != 42 {
		t.Errorf("HpeStorageDatastore = %d, want 42", v.HpeStorageDatastore)
	}
	if v.HpeStorageVolumeShared != nil {
		t.Errorf("HpeStorageVolumeShared = %q, want nil", *v.HpeStorageVolumeShared)
	}
	if v.HpeStorageComputeServer != nil {
		t.Errorf("HpeStorageComputeServer = %q, want nil", *v.HpeStorageComputeServer)
	}
	if len(v.HpeStorageInstances) != 0 {
		t.Errorf("HpeStorageInstances = %v, want empty", v.HpeStorageInstances)
	}
}

// TestBuildCreateConfigEmpty confirms that when neither config block is set the
// request carries no config union (nil), so the API receives no config object.
func TestBuildCreateConfigEmpty(t *testing.T) {
	ctx := context.Background()

	model := &StorageVolumeModel{
		ConfigAlletrampBmaas: NewConfigAlletrampBmaasValueNull(),
		Config:               types.DynamicNull(),
	}

	var diags diag.Diagnostics
	cfg := buildCreateConfig(ctx, model, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if cfg != nil {
		t.Errorf("expected nil config, got %+v", cfg)
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
