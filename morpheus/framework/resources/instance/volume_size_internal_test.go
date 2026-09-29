// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

const gib = int64(1024 * 1024 * 1024)

// planWithVolumeSize builds a plan holding a single volume whose size is the
// given value, so the read path can be driven with a configured size, an
// omitted one, or an unknown one.
func planWithVolumeSize(t *testing.T, size basetypes.Int64Value) InstanceModel {
	t.Helper()

	ctx := context.Background()

	vol, d := NewVolumesValue(
		VolumesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"actual_size":              types.Int64Null(),
			"controller_mount_point":   types.StringNull(),
			"datastore_auto_selection": types.StringNull(),
			"datastore_id":             types.Int64Null(),
			"id":                       types.Int64Null(),
			"name":                     types.StringNull(),
			"root_volume":              types.BoolNull(),
			"size":                     size,
			"size_id":                  types.Int64Null(),
			"storage_profile":          types.StringNull(),
			"storage_type_id":          types.Int64Null(),
		},
	)
	if d.HasError() {
		t.Fatalf("building the fixture failed: %v", d)
	}

	list, d := types.ListValue(VolumesValue{}.Type(ctx), []attr.Value{vol})
	if d.HasError() {
		t.Fatalf("building the fixture failed: %v", d)
	}

	return InstanceModel{Volumes: list}
}

// apiVolume is one volume as the API reports it, with MaxStorage in bytes.
func apiVolume(bytes int64) []sdk.InstanceContainerServerVolume1 {
	return []sdk.InstanceContainerServerVolume1{{MaxStorage: &bytes}}
}

// size records what was asked for and actual_size what the platform made, so a
// configured size must survive the read even when the API reports something
// larger.
func TestVolumeSizeKeepsTheRequestWhenConfigured(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Asked for 10GB, the platform made 11GB.
	vols := setDatastoreAutoSelectionAndSize(
		apiVolume(11*gib), planWithVolumeSize(t, types.Int64Value(10)), false)

	list, diags := convertAPIVolumesToStateVolumes(ctx, vols)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	got := firstVolume(t, ctx, list)

	if got.Size.ValueInt64() != 10 {
		t.Errorf("size = %d, want 10 — the request was not preserved",
			got.Size.ValueInt64())
	}

	if got.ActualSize.ValueInt64() != 11 {
		t.Errorf("actual_size = %d, want 11 — the provisioned size was lost",
			got.ActualSize.ValueInt64())
	}
}

// size is Optional and Computed, so omitting it means the platform decides. The
// read must then report what was provisioned rather than recording null.
//
// This guards a regression introduced when size became Computed: the plan value
// was written over the API's unconditionally, so an omitted size wrote nil
// while still flagging that Terraform had set it, and size came back null.
func TestVolumeSizeFallsBackToTheAPIWhenOmitted(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for _, tc := range []struct {
		name string
		size basetypes.Int64Value
	}{
		{name: "omitted", size: types.Int64Null()},
		{name: "unknown", size: types.Int64Unknown()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			vols := setDatastoreAutoSelectionAndSize(
				apiVolume(11*gib), planWithVolumeSize(t, tc.size), false)

			list, diags := convertAPIVolumesToStateVolumes(ctx, vols)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			got := firstVolume(t, ctx, list)

			if got.Size.IsNull() {
				t.Error("size is null; a computed attribute must report the " +
					"size the platform provisioned when none was requested")
			}

			if got.Size.ValueInt64() != 11 {
				t.Errorf("size = %d, want 11 (the provisioned size in GB)",
					got.Size.ValueInt64())
			}

			if got.ActualSize.ValueInt64() != 11 {
				t.Errorf("actual_size = %d, want 11", got.ActualSize.ValueInt64())
			}
		})
	}
}

// MaxStorage is bytes on the API and gigabytes in state, and the flag saying
// which unit it currently holds is what keeps the two apart.
func TestVolumeSizeConvertsBytesToGigabytes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	list, diags := convertAPIVolumesToStateVolumes(ctx, apiVolume(40*gib))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	got := firstVolume(t, ctx, list)

	if got.Size.ValueInt64() != 40 {
		t.Errorf("size = %d, want 40 — bytes were not converted",
			got.Size.ValueInt64())
	}
}

func firstVolume(
	t *testing.T,
	ctx context.Context,
	list basetypes.ListValue,
) VolumesValue {
	t.Helper()

	elems := list.Elements()
	if len(elems) != 1 {
		t.Fatalf("got %d volumes, want 1", len(elems))
	}

	v, ok := elems[0].(VolumesValue)
	if !ok {
		t.Fatalf("element is %T, want VolumesValue", elems[0])
	}

	return v
}
