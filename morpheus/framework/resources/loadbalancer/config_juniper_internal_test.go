// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The SDN load balancer block is sent under the plugin's option type keys
// and read back from them; unset optional fields are omitted, not nulled.
func TestUnitJuniperSdnLoadBalancerConfigRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const uid = "0be6b79b-757c-4d92-affc-8084e328cb4c"

	cfg, d := NewConfigJuniperSdnValue(ConfigJuniperSdnValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"resource_pool_id":    types.Int64Value(153047),
		"virtual_network_uid": types.StringValue(uid),
		"virtual_ip":          types.StringValue("10.0.0.10"),
		"source_port":         types.Int64Value(443),
		"destination_port":    types.Int64Null(),
		"protocol":            types.StringValue("TCP"),
	})
	if d.HasError() {
		t.Fatal(d)
	}

	b, _ := json.Marshal(juniperSdnConfig(cfg))
	want := `{"juniperResourcePool":153047,"protocol":"TCP",` +
		`"sourcePort":443,"virtualIp":"10.0.0.10","virtualNetwork":"` + uid + `"}`
	if string(b) != want {
		t.Errorf("create body = %s, want %s", b, want)
	}

	// Morpheus hands option type values back as it stored them: numbers or
	// numeric strings.
	back, err := parseJuniperSdnConfig(ctx, map[string]interface{}{
		"juniperResourcePool": "153047", "virtualNetwork": uid, "virtualIp": "10.0.0.10",
		"sourcePort": float64(443), "protocol": "TCP",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !back.ResourcePoolId.Equal(types.Int64Value(153047)) || !back.VirtualNetworkUid.Equal(types.StringValue(uid)) ||
		back.VirtualIp.ValueString() != "10.0.0.10" || !back.SourcePort.Equal(types.Int64Value(443)) ||
		!back.DestinationPort.IsNull() || back.Protocol.ValueString() != "TCP" {
		t.Errorf("read back: %v", back)
	}
}
