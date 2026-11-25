// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitOverlayPlanConfigNullRead pins that a plan overlaid on a null read
// yields a Known config whose unset fields are typed nulls that can be
// written to state — not the zero value's untyped object and dynamic.
func TestUnitOverlayPlanConfigNullRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	plan := ConfigValue{
		HttpMethod:     types.StringValue("GET"),
		HttpUri:        types.StringValue("/health"),
		HttpPort:       types.Int64Unknown(),
		HttpReturnCode: types.StringUnknown(),
		Monitor:        types.ObjectUnknown(MonitorValue{}.AttributeTypes(ctx)),
		MonitorConfig:  types.DynamicUnknown(),
		state:          attr.ValueStateKnown,
	}

	out := overlayPlanConfig(ctx, plan, NewConfigValueNull())

	if out.IsNull() || out.IsUnknown() {
		t.Fatal("overlay on a null read: want a known config")
	}

	if out.HttpMethod.ValueString() != "GET" || out.HttpUri.ValueString() != "/health" {
		t.Errorf("planned values not carried: %+v", out)
	}

	for name, v := range map[string]attr.Value{
		"http_port": out.HttpPort, "http_return_code": out.HttpReturnCode,
		"monitor": out.Monitor, "monitor_config": out.MonitorConfig,
	} {
		if !v.IsNull() {
			t.Errorf("%s: want null, got %v", name, v)
		}
	}

	// The result must be convertible to a Terraform value — the whole point.
	if _, err := out.ToTerraformValue(ctx); err != nil {
		t.Errorf("ToTerraformValue: %v", err)
	}
}
