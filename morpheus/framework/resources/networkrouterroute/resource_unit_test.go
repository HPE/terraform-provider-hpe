// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapRouteToStateNullOptionals guards the read-path nil-safety fix: an explicit JSON
// null (or an omitted field) for description and network_mtu must map to a Terraform
// null without panicking. The SDK marks an explicit null as IsSet with a nil Get(), so
// the prior *Get() dereference panicked on import or refresh of a route omitting these.
func TestMapRouteToStateNullOptionals(t *testing.T) {
	t.Parallel()

	plan := NetworkRouterRouteModel{RouterId: types.Int64Value(7)}

	id := int64(3)
	tests := map[string]sdk.GetNetworkRouterRoute200ResponseNetworkRoute{
		"explicit null optionals": {
			Id:          &id,
			Description: *sdk.NewNullableString(nil),
			NetworkMtu:  *sdk.NewNullableFloat32(nil),
			Priority:    *sdk.NewNullableString(nil),
		},
		"absent optionals (zero value)": {Id: &id},
	}

	for name, route := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			state := mapRouteToState(&route, plan) // must not panic

			if !state.Description.IsNull() {
				t.Errorf("description = %v, want null", state.Description)
			}

			if !state.NetworkMtu.IsNull() {
				t.Errorf("network_mtu = %v, want null", state.NetworkMtu)
			}

			// RouterId is carried from the plan, not echoed by the API response.
			if !state.RouterId.Equal(plan.RouterId) {
				t.Errorf("router_id = %v, want %v", state.RouterId, plan.RouterId)
			}
		})
	}
}

// TestMapRouteToStatePopulatedOptionals confirms present values still map through.
func TestMapRouteToStatePopulatedOptionals(t *testing.T) {
	t.Parallel()

	id := int64(3)
	desc := "north-south"
	mtu := float32(1500)
	route := sdk.GetNetworkRouterRoute200ResponseNetworkRoute{
		Id:          &id,
		Description: *sdk.NewNullableString(&desc),
		NetworkMtu:  *sdk.NewNullableFloat32(&mtu),
	}

	state := mapRouteToState(&route, NetworkRouterRouteModel{RouterId: types.Int64Value(7)})

	if state.Description.IsNull() || state.Description.ValueString() != desc {
		t.Errorf("description = %v, want %q", state.Description, desc)
	}

	if state.NetworkMtu.IsNull() || state.NetworkMtu.ValueFloat64() != float64(mtu) {
		t.Errorf("network_mtu = %v, want %v", state.NetworkMtu, mtu)
	}
}
