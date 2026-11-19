// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func int64Ptr(v int64) *int64 { return &v }
func strPtr(s string) *string { return &s }

// The short import form leaves the two inputs nil; the long form parses all
// three integers.
func TestUnitParseImportID(t *testing.T) {
	t.Parallel()

	short, err := parseImportID("42")
	if err != nil {
		t.Fatalf("short form: %s", err)
	}
	if short.id != 42 || short.networkServerID != nil || short.floatingIpPoolID != nil {
		t.Errorf("short form = %+v, want id=42 with nil inputs", short)
	}

	long, err := parseImportID("42:7:9")
	if err != nil {
		t.Fatalf("long form: %s", err)
	}
	if long.id != 42 || long.networkServerID == nil || *long.networkServerID != 7 ||
		long.floatingIpPoolID == nil || *long.floatingIpPoolID != 9 {
		t.Errorf("long form = %+v, want id=42 nsid=7 poolid=9", long)
	}

	for _, bad := range []string{"", "x", "42:7", "42:7:9:1", "42:x:9"} {
		if _, err := parseImportID(bad); err == nil {
			t.Errorf("parseImportID(%q) succeeded, want error", bad)
		}
	}
}

// successMsg reads the human-readable msg from a success=false envelope.
func TestUnitSuccessMsg(t *testing.T) {
	t.Parallel()

	if got := successMsg(map[string]interface{}{"msg": "no free IPs"}); got != "no free IPs" {
		t.Errorf("successMsg = %q, want %q", got, "no free IPs")
	}
	if got := successMsg(nil); got != "no message returned" {
		t.Errorf("successMsg(nil) = %q, want fallback", got)
	}
}

// floatingIpAsState carries the two API-omitted inputs forward from prior state
// and maps the nested sub-structs, treating nil nested objects as null.
func TestUnitFloatingIpAsState(t *testing.T) {
	t.Parallel()

	prior := NetworkFloatingIpModel{
		NetworkServerId:  types.Int64Value(7),
		FloatingIpPoolId: types.Int64Value(9),
	}

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	fip := &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIp{
		Id:          int64Ptr(42),
		IpAddress:   strPtr("10.0.0.5"),
		IpStatus:    strPtr("assigned"),
		DateCreated: &created,
		Cloud: &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIpCloud{
			Id:   int64Ptr(3),
			Name: strPtr("hvm"),
		},
		Server: &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIpServer{
			Id:   int64Ptr(11),
			Name: strPtr("web-01"),
		},
		IpRange: &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIpIpRange{
			StartAddress: strPtr("10.0.0.1"),
			EndAddress:   strPtr("10.0.0.254"),
		},
	}

	state := floatingIpAsState(fip, &prior)

	if state.Id.ValueInt64() != 42 {
		t.Errorf("id = %d, want 42", state.Id.ValueInt64())
	}
	if state.NetworkServerId.ValueInt64() != 7 || state.FloatingIpPoolId.ValueInt64() != 9 {
		t.Errorf("inputs not carried forward: nsid=%d poolid=%d",
			state.NetworkServerId.ValueInt64(), state.FloatingIpPoolId.ValueInt64())
	}
	if state.CloudName.ValueString() != "hvm" || state.ServerName.ValueString() != "web-01" {
		t.Errorf("nested names wrong: cloud=%q server=%q",
			state.CloudName.ValueString(), state.ServerName.ValueString())
	}
	if state.IpRangeStartAddress.ValueString() != "10.0.0.1" {
		t.Errorf("range start = %q, want 10.0.0.1", state.IpRangeStartAddress.ValueString())
	}
	if !state.NetworkDomainId.IsNull() || !state.NetworkDomainName.IsNull() {
		t.Errorf("nil network domain should be null, got id=%v name=%v",
			state.NetworkDomainId, state.NetworkDomainName)
	}
	if state.DateCreated.ValueString() == "" {
		t.Error("date_created should be set")
	}
}
