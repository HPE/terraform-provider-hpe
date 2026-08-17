// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancedisktype

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strPtr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }

func int32Ptr(i int32) *int32 { return &i }

// st builds a disk type with the given id and name, the two fields the lookup
// keys on.
func st(id int32, name string) diskType {
	return diskType{Id: int32Ptr(id), Name: strPtr(name)}
}

func plan(sts ...diskType) sdk.ListInstanceServicePlans200ResponsePlansInner {
	return sdk.ListInstanceServicePlans200ResponsePlansInner{StorageTypes: sts}
}

// TestCollectDiskTypesDedupesByID pins the parity-critical behaviour: Morpheus
// repeats the same storage types in every service plan, so the flattened list
// must be deduped by id or every lookup would report multiple matches.
func TestCollectDiskTypesDedupesByID(t *testing.T) {
	plans := []sdk.ListInstanceServicePlans200ResponsePlansInner{
		plan(st(1, "Standard"), st(2, "Thin")),
		plan(st(1, "Standard"), st(2, "Thin")),
		plan(st(1, "Standard")),
	}

	got := collectDiskTypes(plans)
	if len(got) != 2 {
		t.Fatalf("want 2 deduped disk types, got %d", len(got))
	}
}

// TestCollectDiskTypesSkipsNilID verifies disk types without an id are dropped,
// since the id is the value the data source exists to return.
func TestCollectDiskTypesSkipsNilID(t *testing.T) {
	plans := []sdk.ListInstanceServicePlans200ResponsePlansInner{
		plan(diskType{Name: strPtr("no id")}, st(5, "Standard")),
	}

	got := collectDiskTypes(plans)
	if len(got) != 1 || got[0].Id == nil || *got[0].Id != 5 {
		t.Fatalf("want only the disk type with an id, got %+v", got)
	}
}

func TestMatchDiskType(t *testing.T) {
	diskTypes := []diskType{
		st(1, "Standard"),
		st(2, "Thin"),
		st(3, "Thick"),
	}

	tests := []struct {
		name       string
		lookupName string
		input      []diskType
		wantID     int32
		wantErr    string
	}{
		{
			name:       "exact match",
			lookupName: "Thin",
			input:      diskTypes,
			wantID:     2,
		},
		{
			name:       "case-insensitive and whitespace-trimmed match (hpegl parity)",
			lookupName: "  standard  ",
			input:      diskTypes,
			wantID:     1,
		},
		{
			name:       "no match errors",
			lookupName: "NVMe",
			input:      diskTypes,
			wantErr:    ErrorNoInstanceDiskTypeFound,
		},
		{
			name:       "empty list errors as not found",
			lookupName: "Standard",
			input:      nil,
			wantErr:    ErrorNoInstanceDiskTypeFound,
		},
		{
			name:       "multiple matches errors",
			lookupName: "dup",
			input:      []diskType{st(4, "dup"), st(5, "DUP")},
			wantErr:    ErrorMultipleInstanceDiskTypes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchDiskType(tt.input, tt.lookupName)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Id == nil || *got.Id != tt.wantID {
				t.Errorf("id: want %d, got %v", tt.wantID, got.Id)
			}
		})
	}
}

// TestDiskTypeAsState covers the three mapping shapes the model mixes: plain
// pointers, NullableString/NullableBool (which must be read through Get), and
// int32 widening. It also pins that absent values stay null rather than being
// flattened to a zero value.
func TestDiskTypeAsState(t *testing.T) {
	in := diskType{
		Id:           int32Ptr(1),
		Name:         strPtr("Standard"),
		Code:         strPtr("standard"),
		DisplayOrder: int32Ptr(3),
		VolumeType:   strPtr("disk"),
		Enabled:      boolPtr(true),
		HasISO:       boolPtr(false),
	}
	in.Description.Set(strPtr("Standard"))
	in.StorageType.Set(strPtr("block"))
	in.HasActiveReplica.Set(boolPtr(true))
	// MinIOPS, MaxStorage, ExternalId and VolumeOptionSource left unset.

	got := diskTypeAsState(in, InstanceDiskTypeModel{})

	if got.Id.ValueInt64() != 1 {
		t.Errorf("id: want 1, got %v", got.Id)
	}

	if got.DisplayOrder.ValueInt64() != 3 {
		t.Errorf("display_order: want 3 (int32 widened), got %v", got.DisplayOrder)
	}

	if got.Description.ValueString() != "Standard" {
		t.Errorf("description: want %q, got %v", "Standard", got.Description)
	}

	if got.StorageType.ValueString() != "block" {
		t.Errorf("storage_type: want %q, got %v", "block", got.StorageType)
	}

	if !got.HasActiveReplica.ValueBool() {
		t.Errorf("has_active_replica: want true, got %v", got.HasActiveReplica)
	}

	if got.HasIso.ValueBool() {
		t.Errorf("has_iso: want false, got %v", got.HasIso)
	}

	for name, v := range map[string]interface{ IsNull() bool }{
		"min_iops":             got.MinIops,
		"max_storage":          got.MaxStorage,
		"external_id":          got.ExternalId,
		"volume_option_source": got.VolumeOptionSource,
		"deletable":            got.Deletable,
	} {
		if !v.IsNull() {
			t.Errorf("%s: want null when the API omits it, got %v", name, v)
		}
	}
}
