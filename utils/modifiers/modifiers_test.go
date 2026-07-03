// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package modifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizeLineEndingsModifier(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		plan    types.String
		want    string
		wantSet bool
	}{
		"crlf normalized to lf": {
			plan:    types.StringValue("line1\r\nline2\r\nline3"),
			want:    "line1\nline2\nline3",
			wantSet: true,
		},
		"lone cr normalized to lf": {
			plan:    types.StringValue("line1\rline2"),
			want:    "line1\nline2",
			wantSet: true,
		},
		"mixed cr and crlf": {
			plan:    types.StringValue("a\r\nb\rc\n"),
			want:    "a\nb\nc\n",
			wantSet: true,
		},
		"lf only left unchanged": {
			plan:    types.StringValue("already\nnormalized\n"),
			wantSet: false,
		},
		"no line endings left unchanged": {
			plan:    types.StringValue("single line"),
			wantSet: false,
		},
		"null plan left unchanged": {
			plan:    types.StringNull(),
			wantSet: false,
		},
		"unknown plan left unchanged": {
			plan:    types.StringUnknown(),
			wantSet: false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// ConfigValue is deliberately null for every case. The modifier must
			// normalize the accumulated PlanValue (as populated by prior plan
			// modifiers), never the raw ConfigValue — so a case such as "crlf
			// normalized to lf" with a null ConfigValue and a CRLF PlanValue
			// exercises the chained-modifier scenario and would fail if the
			// modifier read ConfigValue.
			resp := &planmodifier.StringResponse{PlanValue: tc.plan}
			req := planmodifier.StringRequest{
				ConfigValue: types.StringNull(),
				PlanValue:   tc.plan,
			}

			NormalizeLineEndingsModifier{}.PlanModifyString(context.Background(), req, resp)

			if tc.wantSet {
				if got := resp.PlanValue.ValueString(); got != tc.want {
					t.Errorf("PlanValue = %q, want %q", got, tc.want)
				}

				return
			}

			// When nothing needs normalizing the modifier must leave PlanValue
			// as the accumulated planned value.
			if !resp.PlanValue.Equal(tc.plan) {
				t.Errorf("PlanValue changed unexpectedly: got %#v, want %#v", resp.PlanValue, tc.plan)
			}
		})
	}
}
