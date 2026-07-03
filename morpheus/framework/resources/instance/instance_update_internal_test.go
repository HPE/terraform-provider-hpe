// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mkVol builds a VolumesValue for the resize-matching tests. A null id represents
// a volume the user has just added (no server id yet); a set id represents an
// existing volume carried in state.
func mkVol(id types.Int64, name string, size int64) VolumesValue {
	return VolumesValue{
		Id:                     id,
		Name:                   types.StringValue(name),
		RootVolume:             types.BoolNull(),
		Size:                   types.Int64Value(size),
		StorageTypeId:          types.Int64Null(),
		DatastoreId:            types.Int64Null(),
		DatastoreAutoSelection: types.StringNull(),
		StorageProfile:         types.StringNull(),
	}
}

func ptrInt64(v int64) *int64 { return &v }

// TestUnitBuildResizeVolumes pins the plan<->state volume matching used to build an
// instance resize request. The morpheus-ui resize changelist matches existing
// volumes by id, so the request must carry the existing (state) id for kept/resized
// volumes, id=-1 (null in the model, set by updateVolumeMapper) for added volumes,
// and omit dropped volumes. When the counts match, pairing is positional (unchanged
// behaviour: a rename with no size change is not a resize); when they differ, matching
// is by name so a middle insert/remove does not mis-assign ids.
func TestUnitBuildResizeVolumes(t *testing.T) {
	t.Parallel()

	type wantVol struct {
		id   *int64 // nil => new/null id
		name string
		size int64
	}

	tests := []struct {
		name         string
		plan         []VolumesValue
		state        []VolumesValue
		wantResizing bool
		want         []wantVol
	}{
		{
			name: "no change keeps state ids and does not resize",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "root", 10),
				mkVol(types.Int64Value(2), "data", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "root", 10),
				mkVol(types.Int64Null(), "data", 20),
			},
			wantResizing: false,
			want: []wantVol{
				{ptrInt64(1), "root", 10},
				{ptrInt64(2), "data", 20},
			},
		},
		{
			name: "size change resizes the matched volume",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "root", 10),
				mkVol(types.Int64Value(2), "data", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "root", 10),
				mkVol(types.Int64Null(), "data", 30),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "root", 10},
				{ptrInt64(2), "data", 30},
			},
		},
		{
			name: "add at end",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "root", 10),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "root", 10),
				mkVol(types.Int64Null(), "extra", 5),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "root", 10},
				{nil, "extra", 5},
			},
		},
		{
			// The case the old positional logic got wrong: inserting "mid" between
			// "root" and "data" must not steal data's id.
			name: "add in the middle keeps existing ids by name",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "root", 10),
				mkVol(types.Int64Value(2), "data", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "root", 10),
				mkVol(types.Int64Null(), "mid", 5),
				mkVol(types.Int64Null(), "data", 20),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "root", 10},
				{nil, "mid", 5},
				{ptrInt64(2), "data", 20},
			},
		},
		{
			// Removing "mid" must drop id 2 and keep data's id 3 (positional logic
			// would have mis-mapped data to id 2).
			name: "remove in the middle drops the right volume",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "root", 10),
				mkVol(types.Int64Value(2), "mid", 5),
				mkVol(types.Int64Value(3), "data", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "root", 10),
				mkVol(types.Int64Null(), "data", 20),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "root", 10},
				{ptrInt64(3), "data", 20},
			},
		},
		{
			// Count is unchanged, so positional pairing applies (unchanged
			// behaviour): a rename with no size change must NOT trigger a resize —
			// the volume keeps its state id and name.
			name: "rename with same size does not resize",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "data", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "newdata", 20),
			},
			wantResizing: false,
			want: []wantVol{
				{ptrInt64(1), "data", 20},
			},
		},
		{
			name: "unnamed volumes fall back to positional matching",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "", 10),
				mkVol(types.Int64Value(2), "", 20),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "", 10),
				mkVol(types.Int64Null(), "", 30),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "", 10},
				{ptrInt64(2), "", 30},
			},
		},
		{
			name: "unnamed add appends a new volume",
			state: []VolumesValue{
				mkVol(types.Int64Value(1), "", 10),
			},
			plan: []VolumesValue{
				mkVol(types.Int64Null(), "", 10),
				mkVol(types.Int64Null(), "", 5),
			},
			wantResizing: true,
			want: []wantVol{
				{ptrInt64(1), "", 10},
				{nil, "", 5},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, resizing := buildResizeVolumes(tc.plan, tc.state)

			if resizing != tc.wantResizing {
				t.Errorf("resizing = %v, want %v", resizing, tc.wantResizing)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("got %d volumes, want %d", len(got), len(tc.want))
			}

			for i, w := range tc.want {
				gotID := got[i].Id.ValueInt64Pointer()
				switch {
				case w.id == nil && gotID != nil:
					t.Errorf("volume %d (%s): id = %d, want null (new volume)", i, w.name, *gotID)
				case w.id != nil && gotID == nil:
					t.Errorf("volume %d (%s): id = null, want %d", i, w.name, *w.id)
				case w.id != nil && gotID != nil && *gotID != *w.id:
					t.Errorf("volume %d (%s): id = %d, want %d", i, w.name, *gotID, *w.id)
				}

				if gotName := got[i].Name.ValueString(); gotName != w.name {
					t.Errorf("volume %d: name = %q, want %q", i, gotName, w.name)
				}

				if gotSize := got[i].Size.ValueInt64(); gotSize != w.size {
					t.Errorf("volume %d (%s): size = %d, want %d", i, w.name, gotSize, w.size)
				}
			}
		})
	}
}
