// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitMonitorPasswordHash checks the appliance-hash extractor: a set hash
// is returned present; a present-but-null field is present with an empty hash
// (the appliance holds no secret); an unset (API-omitted) field is not present
// (no information); a wrong type is not present.
func TestUnitMonitorPasswordHash(t *testing.T) {
	t.Parallel()

	hash := "deadbeef"

	set := &sdk.GetLoadBalancerMonitor200ResponseLoadBalancerMonitor{
		MonitorPasswordHash: *sdk.NewNullableString(&hash),
	}
	if h, ok := monitorPasswordHash(set); !ok || h != hash {
		t.Fatalf("set hash: got (%q,%v), want (%q,true)", h, ok, hash)
	}

	null := &sdk.GetLoadBalancerMonitor200ResponseLoadBalancerMonitor{
		MonitorPasswordHash: *sdk.NewNullableString(nil),
	}
	if h, ok := monitorPasswordHash(null); !ok || h != "" {
		t.Fatalf("null hash: got (%q,%v), want (\"\",true)", h, ok)
	}

	unset := &sdk.GetLoadBalancerMonitor200ResponseLoadBalancerMonitor{}
	if h, ok := monitorPasswordHash(unset); ok || h != "" {
		t.Fatalf("unset hash: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := monitorPasswordHash("not the object"); ok || h != "" {
		t.Fatalf("wrong type: got (%q,%v), want (\"\",false)", h, ok)
	}

	if h, ok := monitorPasswordHash((*sdk.GetLoadBalancerMonitor200ResponseLoadBalancerMonitor)(nil)); ok || h != "" {
		t.Fatalf("nil object: got (%q,%v), want (\"\",false)", h, ok)
	}
}
