// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapRuleToStateProtocolPortRange covers MORPH-16362: protocol and
// port_range are not reliably returned by the API, so the configured value must
// be preserved rather than nulled out after apply.
func TestMapRuleToStateProtocolPortRange(t *testing.T) {
	apiVal := "tcp"

	tests := []struct {
		name         string
		apiProtocol  *string
		planProtocol types.String
		wantProtocol types.String
		wantAbsent   bool
	}{
		{"api returns value", &apiVal, types.StringValue("udp"), types.StringValue("tcp"), false},
		{"api null, plan known preserves", nil, types.StringValue("tcp"), types.StringValue("tcp"), true},
		{"api null, plan unknown nulls", nil, types.StringUnknown(), types.StringNull(), true},
		{"api null, plan null nulls", nil, types.StringNull(), types.StringNull(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &sdk.GetNetworkRouterFirewallRule200ResponseRule{}
			if tt.apiProtocol != nil {
				rule.Protocol.Set(tt.apiProtocol)
			}

			plan := NetworkRouterFirewallRuleModel{Protocol: tt.planProtocol}
			state, absent := mapRuleToState(rule, plan)

			if !state.Protocol.Equal(tt.wantProtocol) {
				t.Errorf("Protocol = %v, want %v", state.Protocol, tt.wantProtocol)
			}
			if absent.Protocol != tt.wantAbsent {
				t.Errorf("absent.Protocol = %v, want %v", absent.Protocol, tt.wantAbsent)
			}
		})
	}
}
