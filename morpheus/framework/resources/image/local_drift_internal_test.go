package image

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUnitLocalFileRemoteDriftPlansInPlaceRepair(t *testing.T) {
	t.Parallel()
	for _, missing := range []bool{false, true} {
		contents := []byte("complete image")
		store := newFakeImageStore(map[string][]byte{"disk.iso": []byte("partial")})
		if missing {
			delete(store.files, "disk.iso")
		}
		client := newFakeClient(t, store)
		values := map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.Number, 42),
			"file_path":   tftypes.NewValue(tftypes.String, writeTempFile(t, "disk.iso", contents)),
			"file_name":   tftypes.NewValue(tftypes.String, "disk.iso"),
			"file_size":   tftypes.NewValue(tftypes.Number, len(contents)),
			"file_sha256": tftypes.NewValue(tftypes.String, sumOf(contents)),
		}
		config, plan, state := localObject(t, values)
		ctx := context.Background()
		var previous ImageModel
		if diagnostics := state.Get(ctx, &previous); diagnostics.HasError() {
			t.Fatal(diagnostics)
		}
		refreshed, _, found, diagnostics := getImageAsState(ctx, 42, client, previous)
		if diagnostics.HasError() || !found || refreshed.Id.ValueInt64() != 42 {
			t.Fatalf("file drift lost image identity: %v", diagnostics)
		}
		if refreshed.FileSize.ValueInt64() == int64(len(contents)) {
			t.Fatal("file listing drift was not recorded")
		}
		if diagnostics := state.Set(ctx, &refreshed); diagnostics.HasError() {
			t.Fatal(diagnostics)
		}
		response := resource.ModifyPlanResponse{Plan: plan}
		planLocalFile(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan, State: state}, &response)
		if response.Diagnostics.HasError() || len(response.RequiresReplace) != 0 {
			t.Fatal("file drift must not replace the image")
		}
		var desired ImageModel
		if diagnostics := response.Plan.Get(ctx, &desired); diagnostics.HasError() {
			t.Fatal(diagnostics)
		}
		if desired.FileSize == refreshed.FileSize {
			t.Fatal("file drift did not produce a repair diff")
		}
		if err := applyLocalFile(ctx, client, &desired, refreshed); err != nil {
			t.Fatal(err)
		}
		if string(store.files["disk.iso"]) != string(contents) || desired.Id.ValueInt64() != 42 {
			t.Fatal("repair did not restore content in place")
		}
	}
}
