// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancedisktype

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strPtr(s string) *string { return &s }

func int32Ptr(i int32) *int32 { return &i }

// st builds a storage type (disk type) with the given id and name, the two
// fields the lookup keys on.
func st(id int32, name string) sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner {
	return sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner{
		Id:   int32Ptr(id),
		Name: strPtr(name),
	}
}

func plan(
	sts ...sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner,
) sdk.ListInstanceServicePlans200ResponsePlansInner {
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

// TestCollectDiskTypesSkipsNilID verifies storage types without an id are
// dropped, since the id is the value the data source exists to return.
func TestCollectDiskTypesSkipsNilID(t *testing.T) {
	plans := []sdk.ListInstanceServicePlans200ResponsePlansInner{
		plan(
			sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner{Name: strPtr("no id")},
			st(5, "Standard"),
		),
	}

	got := collectDiskTypes(plans)
	if len(got) != 1 || got[0].Id == nil || *got[0].Id != 5 {
		t.Fatalf("want only the disk type with an id, got %+v", got)
	}
}

func TestMatchDiskType(t *testing.T) {
	diskTypes := []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner{
		st(1, "Standard"),
		st(2, "Thin"),
		st(3, "Thick"),
	}

	tests := []struct {
		name       string
		lookupName string
		input      []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner
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
			input: []sdk.ListInstanceServicePlans200ResponsePlansInnerStorageTypesInner{
				st(4, "dup"),
				st(5, "DUP"),
			},
			wantErr: ErrorMultipleInstanceDiskTypes,
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

			if got.id == nil || *got.id != tt.wantID {
				t.Errorf("id: want %d, got %v", tt.wantID, got.id)
			}
		})
	}
}

func TestInt32PtrToType(t *testing.T) {
	if v := int32PtrToType(nil); !v.IsNull() {
		t.Errorf("nil int32 should map to a null Int64, got %v", v)
	}

	if v := int32PtrToType(int32Ptr(7)); v.IsNull() || v.ValueInt64() != 7 {
		t.Errorf("int32(7) should map to Int64(7), got %v", v)
	}
}
