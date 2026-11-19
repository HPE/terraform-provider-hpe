package image

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func localObject(t *testing.T, values map[string]tftypes.Value) (tfsdk.Config, tfsdk.Plan, tfsdk.State) {
	t.Helper()
	ctx := context.Background()
	schema := ImageResourceSchema(ctx)
	object := schema.Type().TerraformType(ctx).(tftypes.Object)
	attributes := make(map[string]tftypes.Value, len(object.AttributeTypes))
	for name, kind := range object.AttributeTypes {
		attributes[name] = tftypes.NewValue(kind, nil)
	}
	for name, value := range values {
		attributes[name] = value
	}
	raw := tftypes.NewValue(object, attributes)

	return tfsdk.Config{Schema: schema, Raw: raw}, tfsdk.Plan{Schema: schema, Raw: raw}, tfsdk.State{Schema: schema, Raw: raw}
}

func TestUnitLocalFilePlan(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	filename := writeTempFile(t, "disk.iso", []byte("first"))
	for _, scenario := range []string{"computed", "uppercase", "mismatch", "missing", "removed", "unknown"} {
		t.Run(scenario, func(t *testing.T) {
			values := map[string]tftypes.Value{"file_path": tftypes.NewValue(tftypes.String, filename)}
			switch scenario {
			case "uppercase":
				values["file_sha256"] = tftypes.NewValue(tftypes.String, strings.ToUpper(sumOf([]byte("first"))))
			case "mismatch":
				values["file_sha256"] = tftypes.NewValue(tftypes.String, strings.Repeat("0", 64))
			case "missing":
				values["file_path"] = tftypes.NewValue(tftypes.String, filepath.Join(t.TempDir(), "future.iso"))
			case "removed":
				values["file_path"] = tftypes.NewValue(tftypes.String, nil)
			case "unknown":
				values["file_path"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
			}
			config, plan, state := localObject(t, values)
			response := resource.ModifyPlanResponse{Plan: plan}
			planLocalFile(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan, State: state}, &response)
			if scenario == "mismatch" {
				if !response.Diagnostics.HasError() {
					t.Fatal("bad checksum was accepted")
				}

				return
			}
			if response.Diagnostics.HasError() || len(response.RequiresReplace) != 0 {
				t.Fatalf("local file plan failed or replaced image: %v", response.Diagnostics)
			}
			var actual ImageModel
			if diagnostics := response.Plan.Get(ctx, &actual); diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			switch scenario {
			case "computed":
				if actual.FileSha256.ValueString() != sumOf([]byte("first")) || actual.FileSize.ValueInt64() != 5 {
					t.Fatal("local file was not hashed")
				}
			case "uppercase":
				if actual.FileSha256.ValueString() != strings.ToUpper(sumOf([]byte("first"))) {
					t.Fatal("configured checksum casing changed")
				}
			case "missing":
				if !actual.FileSha256.IsUnknown() || !actual.FileSize.IsUnknown() {
					t.Fatal("apply-time file was not deferred")
				}
			case "removed":
				if !actual.FileSha256.IsNull() || !actual.FileName.IsNull() || !actual.FileSize.IsNull() {
					t.Fatal("file metadata not cleared when relinquishing management")
				}
			}
		})
	}
	config, plan, state := localObject(t, map[string]tftypes.Value{"file_path": tftypes.NewValue(tftypes.String, filename)})
	if err := os.WriteFile(filename, []byte("other"), 0o600); err != nil {
		t.Fatal(err)
	}
	response := resource.ModifyPlanResponse{Plan: plan}
	planLocalFile(ctx, resource.ModifyPlanRequest{Config: config, Plan: plan, State: state}, &response)
	var actual ImageModel
	if diagnostics := response.Plan.Get(ctx, &actual); diagnostics.HasError() || response.Diagnostics.HasError() {
		t.Fatal(diagnostics, response.Diagnostics)
	}
	if actual.FileSha256 != types.StringValue(sumOf([]byte("other"))) {
		t.Fatal("same-size local content change was not detected")
	}
}
