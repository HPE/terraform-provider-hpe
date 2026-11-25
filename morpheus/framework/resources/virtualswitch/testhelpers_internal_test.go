// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func tfstr(v string) types.String { return types.StringValue(v) }
func tfint(v int64) types.Int64   { return types.Int64Value(v) }

// buildSegmentsMapForTest builds a segments map value from a list of
// (key, SegmentsValue) for tests that need prior state.
func buildSegmentsMapForTest(t *testing.T, entries map[string]SegmentsValue) types.Map {
	t.Helper()

	ctx := context.Background()
	segType := SegmentsValue{}.Type(ctx)

	elems := make(map[string]attr.Value, len(entries))
	for k, v := range entries {
		elems[k] = v
	}

	m, d := types.MapValue(segType, elems)
	if d.HasError() {
		t.Fatalf("map: %v", d)
	}

	return m
}

// newSeg builds a minimal known SegmentsValue for tests.
func newSeg(t *testing.T, id int64, trafficType string, vlan *int64) SegmentsValue {
	t.Helper()

	ctx := context.Background()

	vlanVal := types.Int64Null()
	if vlan != nil {
		vlanVal = types.Int64Value(*vlan)
	}

	idVal := types.Int64Null()
	if id != 0 {
		idVal = types.Int64Value(id)
	}

	return SegmentsValue{
		Gateway:           types.StringNull(),
		HostIps:           types.SetNull(HostIpsValue{}.Type(ctx)),
		Id:                idVal,
		InterfaceName:     types.StringNull(),
		IpAddress:         types.StringNull(),
		Mtu:               types.Int64Null(),
		Netmask:           types.StringNull(),
		OvsBridge:         types.StringNull(),
		OvsLibvirtNetwork: types.StringNull(),
		OvsUuid:           types.StringNull(),
		Status:            types.StringNull(),
		StatusMessage:     types.StringNull(),
		TrafficType:       types.StringValue(trafficType),
		VlanId:            vlanVal,
		state:             attr.ValueStateKnown,
	}
}
