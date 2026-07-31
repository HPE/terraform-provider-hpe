// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestPreserveUnreturnedFields covers the fields the Morpheus API accepts on
// create/update but does not echo back: an empty receive_data comes back as
// null, and data_length is only round-tripped for ICMP monitors. Without the
// preservation the apply fails with "Provider produced inconsistent result
// after apply".
func TestPreserveUnreturnedFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		state    LoadBalancerMonitorModel
		prior    LoadBalancerMonitorModel
		wantRecv types.String
		wantLen  types.Int64
	}{
		"empty string configured, API returned null": {
			state: LoadBalancerMonitorModel{
				ReceiveData: types.StringNull(),
				DataLength:  types.Int64Null(),
			},
			prior: LoadBalancerMonitorModel{
				ReceiveData: types.StringValue(""),
				DataLength:  types.Int64Value(0),
			},
			wantRecv: types.StringValue(""),
			wantLen:  types.Int64Value(0),
		},
		"non-empty configured, API returned null": {
			state: LoadBalancerMonitorModel{
				ReceiveData: types.StringNull(),
				DataLength:  types.Int64Null(),
			},
			prior: LoadBalancerMonitorModel{
				ReceiveData: types.StringValue("HTTP/1.1 200"),
				DataLength:  types.Int64Value(56),
			},
			wantRecv: types.StringValue("HTTP/1.1 200"),
			wantLen:  types.Int64Value(56),
		},
		"API returned a value: API wins over prior": {
			state: LoadBalancerMonitorModel{
				ReceiveData: types.StringValue("from-api"),
				DataLength:  types.Int64Value(99),
			},
			prior: LoadBalancerMonitorModel{
				ReceiveData: types.StringValue("from-plan"),
				DataLength:  types.Int64Value(1),
			},
			wantRecv: types.StringValue("from-api"),
			wantLen:  types.Int64Value(99),
		},
		"unset by practitioner: unknown prior must not leak into state": {
			state: LoadBalancerMonitorModel{
				ReceiveData: types.StringNull(),
				DataLength:  types.Int64Null(),
			},
			prior: LoadBalancerMonitorModel{
				ReceiveData: types.StringUnknown(),
				DataLength:  types.Int64Unknown(),
			},
			wantRecv: types.StringNull(),
			wantLen:  types.Int64Null(),
		},
		"null prior (import): stays null": {
			state: LoadBalancerMonitorModel{
				ReceiveData: types.StringNull(),
				DataLength:  types.Int64Null(),
			},
			prior: LoadBalancerMonitorModel{
				ReceiveData: types.StringNull(),
				DataLength:  types.Int64Null(),
			},
			wantRecv: types.StringNull(),
			wantLen:  types.Int64Null(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			state := tc.state
			preserveUnreturnedFields(&state, tc.prior)

			if !state.ReceiveData.Equal(tc.wantRecv) {
				t.Errorf("receive_data = %v, want %v", state.ReceiveData, tc.wantRecv)
			}

			if !state.DataLength.Equal(tc.wantLen) {
				t.Errorf("data_length = %v, want %v", state.DataLength, tc.wantLen)
			}

			if state.ReceiveData.IsUnknown() || state.DataLength.IsUnknown() {
				t.Error("unknown value leaked into state")
			}
		})
	}
}
