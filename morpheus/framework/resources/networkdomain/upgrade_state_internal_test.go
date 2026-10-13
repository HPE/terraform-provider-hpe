// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

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

// The version-0 schema must describe what the previous (SDKv2) release actually
// wrote, and the upgrader must accept it. testdata/state_v0.json is the
// attributes object of a real hpe_morpheus_network_domain instance captured
// from a live appliance (schema_version 0), which is exactly what Terraform
// hands the provider as the raw state to upgrade.
func TestUnitUpgradeStateV0toV1(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	raw, err := os.ReadFile("testdata/state_v0.json")
	if err != nil {
		t.Fatal(err)
	}

	// The file is the whole stored instance { "attributes": {...},
	// "schema_version": 0 }; the raw state the framework decodes is just the
	// attributes object.
	attrs := extractAttributes(t, raw)

	prior := networkDomainSchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	rawState := tfprotov6.RawState{JSON: attrs}

	value, err := rawState.UnmarshalWithOpts(priorType, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		t.Fatalf("decode of the old state against the version-0 schema failed: %v", err)
	}

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: value}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: NetworkDomainResourceSchema(ctx)}}

	upgradeNetworkDomainStateV0toV1(ctx, req, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var upgraded NetworkDomainModel
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatal(diags)
	}

	if upgraded.Id.ValueInt64() != 38785 {
		t.Errorf("id: got %d, want 38785", upgraded.Id.ValueInt64())
	}

	if upgraded.DomainPasswordWoVersion.ValueInt64() != 2 {
		t.Errorf("domain_password_wo_version: got %d, want 2", upgraded.DomainPasswordWoVersion.ValueInt64())
	}

	if upgraded.TenantId.ValueInt64() != 1 {
		t.Errorf("tenant_id: got %d, want 1", upgraded.TenantId.ValueInt64())
	}

	if upgraded.Name.ValueString() != "morph18860-v0cap.example.test" {
		t.Errorf("name: got %q", upgraded.Name.ValueString())
	}

	if !upgraded.DomainController.ValueBool() || !upgraded.Active.ValueBool() {
		t.Errorf("active/domain_controller must carry across: active=%v controller=%v",
			upgraded.Active, upgraded.DomainController)
	}

	// Write-only is never in state; the revision is seeded null so the first
	// plan after the upgrade is one ~ write_only_revision update.
	// The deprecated attributes are still accepted and are carried as they
	// were (the captured state has auto_join_domain = true, domain_password null).
	if !upgraded.AutoJoinDomain.ValueBool() {
		t.Error("auto_join_domain must be carried from the legacy state")
	}
	if !upgraded.DomainPassword.IsNull() {
		t.Error("domain_password must stay null when the legacy state had none")
	}

	if !upgraded.DomainPasswordWo.IsNull() {
		t.Error("domain_password_wo must be null after upgrade")
	}

	if !upgraded.WriteOnlyRevision.IsNull() {
		t.Error("write_only_revision must be null after upgrade")
	}
}

// A legacy state with a non-numeric id is refused rather than carried forward.
func TestUnitUpgradeStateRejectsBadId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := networkDomainSchemaV0()
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
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: NetworkDomainResourceSchema(ctx)}}

	upgradeNetworkDomainStateV0toV1(ctx, req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for a non-numeric id")
	}
}

// extractAttributes returns the "attributes" object of a stored SDKv2 state
// instance as raw JSON, which is what Terraform hands the provider to upgrade.
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
