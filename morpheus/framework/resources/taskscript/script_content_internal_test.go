// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskscript

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitScriptContentSemanticEquality verifies that a plan that differs from
// state only in leading or trailing whitespace is suppressed (the state value
// is kept), while a genuine change is left as planned.
func TestUnitScriptContentSemanticEquality(t *testing.T) {
	t.Parallel()

	m := ScriptContentSemanticEquality()

	cases := []struct {
		name  string
		state types.String
		plan  types.String
		want  string // expected resulting plan value; "" means unchanged-plan check below
		keep  bool   // true: plan must equal state
	}{
		{"trailing newline suppressed", types.StringValue("echo hi"), types.StringValue("echo hi\n"), "echo hi", true},
		{"leading space suppressed", types.StringValue("echo hi"), types.StringValue("  echo hi"), "echo hi", true},
		{"genuine change kept", types.StringValue("echo hi"), types.StringValue("echo bye"), "echo bye", false},
		{"null state ignored", types.StringNull(), types.StringValue("echo hi\n"), "echo hi\n", false},
		{"unknown plan ignored", types.StringValue("echo hi"), types.StringUnknown(), "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{StateValue: tc.state, PlanValue: tc.plan}
			resp := planmodifier.StringResponse{PlanValue: tc.plan}

			m.PlanModifyString(context.Background(), req, &resp)

			if tc.keep {
				if resp.PlanValue.ValueString() != tc.state.ValueString() {
					t.Errorf("expected plan suppressed to state %q, got %q",
						tc.state.ValueString(), resp.PlanValue.ValueString())
				}

				return
			}

			if resp.PlanValue.ValueString() != tc.want {
				t.Errorf("expected plan %q, got %q", tc.want, resp.PlanValue.ValueString())
			}
		})
	}
}
