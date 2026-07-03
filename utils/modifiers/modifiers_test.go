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
		config  types.String
		want    string
		wantSet bool // whether the modifier overrides PlanValue
	}{
		"crlf normalized to lf": {
			config:  types.StringValue("line1\r\nline2\r\nline3"),
			want:    "line1\nline2\nline3",
			wantSet: true,
		},
		"lone cr normalized to lf": {
			config:  types.StringValue("line1\rline2"),
			want:    "line1\nline2",
			wantSet: true,
		},
		"mixed cr and crlf": {
			config:  types.StringValue("a\r\nb\rc\n"),
			want:    "a\nb\nc\n",
			wantSet: true,
		},
		"lf only left unchanged": {
			config:  types.StringValue("already\nnormalized\n"),
			wantSet: false,
		},
		"no line endings left unchanged": {
			config:  types.StringValue("single line"),
			wantSet: false,
		},
		"null config left unchanged": {
			config:  types.StringNull(),
			wantSet: false,
		},
		"unknown config left unchanged": {
			config:  types.StringUnknown(),
			wantSet: false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// The framework initialises PlanValue to the configured value for a
			// configured attribute; mirror that here.
			resp := &planmodifier.StringResponse{PlanValue: tc.config}
			req := planmodifier.StringRequest{ConfigValue: tc.config, PlanValue: tc.config}

			NormalizeLineEndingsModifier{}.PlanModifyString(context.Background(), req, resp)

			if tc.wantSet {
				if got := resp.PlanValue.ValueString(); got != tc.want {
					t.Errorf("PlanValue = %q, want %q", got, tc.want)
				}

				return
			}

			// When nothing needs normalizing the modifier must leave PlanValue
			// as the original configured value.
			if !resp.PlanValue.Equal(tc.config) {
				t.Errorf("PlanValue changed unexpectedly: got %#v, want %#v", resp.PlanValue, tc.config)
			}
		})
	}
}
