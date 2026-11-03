// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// testdata/state_v0.json is the attributes object of a real
// hpe_morpheus_task_powershell_script instance captured from a live appliance
// (schema_version 0, SDKv2), which is what Terraform hands the provider to
// upgrade.
func TestUnitUpgradeStateV0toV1(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	raw, err := os.ReadFile("testdata/state_v0.json")
	if err != nil {
		t.Fatal(err)
	}

	attrs := extractAttributes(t, raw)

	prior := taskPowershellScriptSchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	rawState := tfprotov6.RawState{JSON: attrs}

	value, err := rawState.UnmarshalWithOpts(priorType, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		t.Fatalf("decode of the old state against the version-0 schema failed: %v", err)
	}

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: value}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: TaskPowershellScriptResourceSchema(ctx)}}

	upgradeTaskPowershellScriptStateV0toV1(ctx, req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var upgraded TaskPowershellScriptModel
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatal(diags)
	}

	if upgraded.Id.ValueInt64() != 19778 {
		t.Errorf("id: got %d, want 19778", upgraded.Id.ValueInt64())
	}

	if upgraded.RemoteTargetPasswordWoVersion.ValueInt64() != 2 {
		t.Errorf("remote_target_password_wo_version: got %d, want 2", upgraded.RemoteTargetPasswordWoVersion.ValueInt64())
	}

	if upgraded.Name.ValueString() != "morph18860-v0cap-ps" {
		t.Errorf("name: got %q", upgraded.Name.ValueString())
	}

	if !upgraded.ElevatedShell.ValueBool() || !upgraded.Retryable.ValueBool() {
		t.Errorf("elevated_shell/retryable must carry across: elevated_shell=%v retryable=%v", upgraded.ElevatedShell, upgraded.Retryable)
	}

	if upgraded.ExecuteTarget.ValueString() != "remote" {
		t.Errorf("execute_target: got %q, want remote", upgraded.ExecuteTarget.ValueString())
	}

	if !upgraded.RemoteTargetPasswordWo.IsNull() {
		t.Error("remote_target_password_wo must be null after upgrade")
	}

	// The legacy state held the API's password hash in remote_target_password;
	// it is not a configured value and must not be carried.
	if !upgraded.RemoteTargetPassword.IsNull() {
		t.Error("remote_target_password must be null after upgrade (legacy state held the API hash)")
	}

	if !upgraded.WriteOnlyRevision.IsNull() {
		t.Error("write_only_revision must be null after upgrade")
	}
}

// A legacy state with a non-numeric id is refused rather than carried forward.
func TestUnitUpgradeStateRejectsBadId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := taskPowershellScriptSchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	attrs := map[string]tftypes.Value{}
	for name, typ := range priorType.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}
	attrs["id"] = tftypes.NewValue(tftypes.String, "not-a-number")
	attrs["name"] = tftypes.NewValue(tftypes.String, "n")

	req := resource.UpgradeStateRequest{
		State: &tfsdk.State{Schema: prior, Raw: tftypes.NewValue(priorType, attrs)},
	}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: TaskPowershellScriptResourceSchema(ctx)}}

	upgradeTaskPowershellScriptStateV0toV1(ctx, req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for a non-numeric id")
	}
}

func extractAttributes(t *testing.T, raw []byte) []byte {
	t.Helper()

	var stored struct {
		Attributes json.RawMessage `json:"attributes"`
	}
	if err := json.Unmarshal(raw, &stored); err != nil {
		t.Fatal(err)
	}

	return stored.Attributes
}
