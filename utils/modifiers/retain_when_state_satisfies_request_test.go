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

// nonNullRaw is a stand-in for a populated plan or state object. The modifier
// only inspects whether these are null, to tell create and destroy apart from
// an update.
func nonNullRaw() tftypes.Value {
	return tftypes.NewValue(
		tftypes.Object{AttributeTypes: map[string]tftypes.Type{}},
		map[string]tftypes.Value{},
	)
}

func nullRaw() tftypes.Value {
	return tftypes.NewValue(
		tftypes.Object{AttributeTypes: map[string]tftypes.Type{}},
		nil,
	)
}

func TestRetainWhenStateSatisfiesRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		stateRaw  tftypes.Value
		planRaw   tftypes.Value
		configVal types.Int64
		stateVal  types.Int64
		planVal   types.Int64
		want      attr.Value
	}{
		{
			// The motivating case: the platform grew a 10GB request to 11GB,
			// so state exceeds config. Keeping state avoids planning a shrink
			// on every run. This is also the post-import case.
			name:      "server grew the volume, state retained",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(10),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Value(10),
			want:      types.Int64Value(11),
		},
		{
			name:      "config equals state, state retained",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(10),
			stateVal:  types.Int64Value(10),
			planVal:   types.Int64Value(10),
			want:      types.Int64Value(10),
		},
		{
			// A genuine request to grow must survive untouched so the resize
			// path acts on it.
			name:      "config exceeds state, plan left alone",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(20),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Value(20),
			want:      types.Int64Value(20),
		},
		{
			// Documented no-op: a reduction cannot be told apart from platform
			// growth, and Morpheus cannot shrink a disk in place.
			name:      "config below state is a no-op",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(5),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Value(5),
			want:      types.Int64Value(11),
		},
		{
			name:      "config silent, state carried forward",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Null(),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Unknown(),
			want:      types.Int64Value(11),
		},
		{
			// Create: no prior state to retain, so the request stands.
			name:      "create leaves the plan untouched",
			stateRaw:  nullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(10),
			stateVal:  types.Int64Null(),
			planVal:   types.Int64Value(10),
			want:      types.Int64Value(10),
		},
		{
			name:      "destroy leaves the plan untouched",
			stateRaw:  nonNullRaw(),
			planRaw:   nullRaw(),
			configVal: types.Int64Value(10),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Value(10),
			want:      types.Int64Value(10),
		},
		{
			name:      "unknown state leaves the plan untouched",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Value(10),
			stateVal:  types.Int64Unknown(),
			planVal:   types.Int64Value(10),
			want:      types.Int64Value(10),
		},
		{
			name:      "unknown config leaves the plan untouched",
			stateRaw:  nonNullRaw(),
			planRaw:   nonNullRaw(),
			configVal: types.Int64Unknown(),
			stateVal:  types.Int64Value(11),
			planVal:   types.Int64Unknown(),
			want:      types.Int64Unknown(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.Int64Request{
				State:       tfsdk.State{Raw: tt.stateRaw},
				Plan:        tfsdk.Plan{Raw: tt.planRaw},
				ConfigValue: tt.configVal,
				StateValue:  tt.stateVal,
				PlanValue:   tt.planVal,
			}
			resp := &planmodifier.Int64Response{PlanValue: tt.planVal}

			modifiers.RetainWhenStateSatisfiesRequest().
				PlanModifyInt64(context.Background(), req, resp)

			if !resp.PlanValue.Equal(tt.want) {
				t.Errorf("PlanValue = %s, want %s", resp.PlanValue, tt.want)
			}
		})
	}
}
