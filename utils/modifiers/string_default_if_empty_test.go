package modifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringUnknownIfEmpty(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		configured types.String
		state      types.String
		expected   types.String
	}{
		"null": {
			configured: types.StringNull(),
			expected:   types.StringUnknown(),
		},
		"empty": {
			configured: types.StringValue(""),
			expected:   types.StringValue(""),
		},
		"explicit": {
			configured: types.StringValue("custom"),
			expected:   types.StringValue("custom"),
		},
		"unknown": {
			configured: types.StringUnknown(),
			expected:   types.StringUnknown(),
		},
		"null with known prior state": {
			configured: types.StringNull(),
			state:      types.StringValue("api-resolved"),
			expected:   types.StringValue("api-resolved"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := planmodifier.StringRequest{
				ConfigValue: test.configured,
				PlanValue:   test.configured,
				StateValue:  test.state,
			}
			response := planmodifier.StringResponse{PlanValue: test.configured}

			StringUnknownIfEmpty().PlanModifyString(context.Background(), request, &response)

			if !response.PlanValue.Equal(test.expected) {
				t.Fatalf("expected %s, got %s", test.expected, response.PlanValue)
			}
		})
	}
}
