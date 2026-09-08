// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestRetryDelaySecondsToState guards the Optional+Computed drift fix. When the API omits
// retry_delay_seconds (nil), a non-zero planned value must be preserved rather than
// nulled; a nil API value with a zero or absent plan maps to null. The nil-apiVal rows
// also guard the pointer dereference that previously panicked.
func TestRetryDelaySecondsToState(t *testing.T) {
	t.Parallel()

	forty := int64(45)
	thirty := int64(30)

	tests := map[string]struct {
		planned types.Int64
		apiVal  *int64
		want    types.Int64
	}{
		"api omits, non-zero plan preserved": {
			planned: types.Int64Value(30),
			apiVal:  nil,
			want:    types.Int64Value(30),
		},
		"api omits, zero plan becomes null": {
			planned: types.Int64Value(0),
			apiVal:  nil,
			want:    types.Int64Null(),
		},
		"api omits, absent plan becomes null": {
			planned: types.Int64Null(),
			apiVal:  nil,
			want:    types.Int64Null(),
		},
		"api differs from non-zero plan, plan wins": {
			planned: types.Int64Value(30),
			apiVal:  &forty,
			want:    types.Int64Value(30),
		},
		"api matches non-zero plan, api used": {
			planned: types.Int64Value(30),
			apiVal:  &thirty,
			want:    types.Int64Value(30),
		},
		"absent plan, api value used": {
			planned: types.Int64Null(),
			apiVal:  &forty,
			want:    types.Int64Value(45),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := retryDelaySecondsToState(tc.planned, tc.apiVal); !got.Equal(tc.want) {
				t.Errorf("retryDelaySecondsToState(%v, %v) = %v, want %v",
					tc.planned, tc.apiVal, got, tc.want)
			}
		})
	}
}

// TestTaskTypeCodeToState guards the nil taskType dereference: the API may omit taskType,
// which previously panicked at task.TaskType.Code. An absent taskType (or a taskType with
// no code) maps to a null string; a present code maps through.
func TestTaskTypeCodeToState(t *testing.T) {
	t.Parallel()

	if got := taskTypeCodeToState(nil); !got.IsNull() {
		t.Errorf("nil taskType = %v, want null", got)
	}

	if got := taskTypeCodeToState(&sdk.GetTasks200ResponseAllOfTaskTaskType{}); !got.IsNull() {
		t.Errorf("taskType with nil code = %v, want null", got)
	}

	code := "ansiblePlaybook"
	got := taskTypeCodeToState(&sdk.GetTasks200ResponseAllOfTaskTaskType{Code: &code})
	if got.IsNull() || got.ValueString() != code {
		t.Errorf("taskType code = %v, want %q", got, code)
	}
}
