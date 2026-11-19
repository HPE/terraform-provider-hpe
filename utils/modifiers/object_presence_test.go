// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/modifiers"
)

func TestObjectRequiresReplaceOnPresenceChange(t *testing.T) {
	t.Parallel()

	attrs := map[string]attr.Type{"a": types.StringType}
	set := func(v string) types.Object {
		return types.ObjectValueMust(attrs, map[string]attr.Value{"a": types.StringValue(v)})
	}
	null := types.ObjectNull(attrs)

	// A resource that exists before and after the plan (an update); create
	// and destroy have a null state or plan respectively.
	existing := tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})}
	planned := tfsdk.Plan{Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})}
	none := tftypes.NewValue(tftypes.Object{}, nil)

	cases := []struct {
		name        string
		state, plan types.Object
		stateRaw    tfsdk.State
		planRaw     tfsdk.Plan
		replace     bool
	}{
		{"added", null, set("x"), existing, planned, true},
		{"removed", set("x"), null, existing, planned, true},
		{"edited in place", set("x"), set("y"), existing, planned, false},
		{"unchanged", set("x"), set("x"), existing, planned, false},
		{"both null", null, null, existing, planned, false},
		{"plan unknown", set("x"), types.ObjectUnknown(attrs), existing, planned, false},
		{"create", null, set("x"), tfsdk.State{Raw: none}, planned, false},
		{"destroy", set("x"), null, existing, tfsdk.Plan{Raw: none}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.ObjectRequest{
				State: c.stateRaw, Plan: c.planRaw,
				StateValue: c.state, PlanValue: c.plan, ConfigValue: c.plan,
			}
			resp := &planmodifier.ObjectResponse{PlanValue: c.plan}

			modifiers.ObjectRequiresReplaceOnPresenceChange().PlanModifyObject(context.Background(), req, resp)

			if resp.RequiresReplace != c.replace {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, c.replace)
			}
		})
	}
}
