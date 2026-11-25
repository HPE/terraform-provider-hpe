// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// nullAttrs builds a map of all-null values for an object type's attributes, so a test only has
// to override the handful it cares about.
func nullAttrs(obj tftypes.Object) map[string]tftypes.Value {
	attrs := make(map[string]tftypes.Value, len(obj.AttributeTypes))
	for name, typ := range obj.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}

	return attrs
}

func TestUnitUpgradeAlertCorrelationPolicyStateV0toV1(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := alertCorrelationPolicySchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	attrs := nullAttrs(priorType)
	attrs["id"] = tftypes.NewValue(tftypes.String, "policy-1")
	attrs["name"] = tftypes.NewValue(tftypes.String, "test-policy")
	attrs["enabled_mode"] = tftypes.NewValue(tftypes.String, "ON")

	algorithmCorrelationType := priorType.AttributeTypes["algorithm_correlation"].(tftypes.Object)
	algorithmCorrelationAttrs := nullAttrs(algorithmCorrelationType)
	algorithmCorrelationAttrs["alerts_time_window"] = tftypes.NewValue(tftypes.String, "45")
	attrs["algorithm_correlation"] = tftypes.NewValue(algorithmCorrelationType, algorithmCorrelationAttrs)

	req := resource.UpgradeStateRequest{
		State: &tfsdk.State{Schema: prior, Raw: tftypes.NewValue(priorType, attrs)},
	}
	resp := resource.UpgradeStateResponse{
		State: tfsdk.State{Schema: alertCorrelationPolicySchema(1, int64TimeWindowAttribute())},
	}

	upgradeAlertCorrelationPolicyStateV0toV1(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var upgraded AlertCorrelationPolicyModel
	if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
		t.Fatal(diags)
	}

	if upgraded.Id.ValueString() != "policy-1" {
		t.Errorf("id: got %q, want policy-1", upgraded.Id.ValueString())
	}

	if upgraded.Name.ValueString() != "test-policy" {
		t.Errorf("name: got %q, want test-policy", upgraded.Name.ValueString())
	}

	if upgraded.AlgorithmCorrelation == nil {
		t.Fatal("algorithm_correlation must not be nil")
	}

	if upgraded.AlgorithmCorrelation.AlertsTimeWindow.ValueInt64() != 45 {
		t.Errorf(
			"alerts_time_window: got %d, want 45",
			upgraded.AlgorithmCorrelation.AlertsTimeWindow.ValueInt64(),
		)
	}
}

// A legacy alerts_time_window that is not a number is refused rather than silently dropped.
func TestUnitUpgradeAlertCorrelationPolicyStateRejectsBadAlertsTimeWindow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := alertCorrelationPolicySchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	attrs := nullAttrs(priorType)
	attrs["id"] = tftypes.NewValue(tftypes.String, "policy-1")
	attrs["name"] = tftypes.NewValue(tftypes.String, "test-policy")

	algorithmCorrelationType := priorType.AttributeTypes["algorithm_correlation"].(tftypes.Object)
	algorithmCorrelationAttrs := nullAttrs(algorithmCorrelationType)
	algorithmCorrelationAttrs["alerts_time_window"] = tftypes.NewValue(tftypes.String, "not-a-number")
	attrs["algorithm_correlation"] = tftypes.NewValue(algorithmCorrelationType, algorithmCorrelationAttrs)

	req := resource.UpgradeStateRequest{
		State: &tfsdk.State{Schema: prior, Raw: tftypes.NewValue(priorType, attrs)},
	}
	resp := resource.UpgradeStateResponse{
		State: tfsdk.State{Schema: alertCorrelationPolicySchema(1, int64TimeWindowAttribute())},
	}

	upgradeAlertCorrelationPolicyStateV0toV1(ctx, req, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for a non-numeric legacy alerts_time_window")
	}
}
