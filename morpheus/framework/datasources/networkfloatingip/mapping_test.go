// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip_test

import (
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfloatingip"
)

func strptr(s string) *string { return &s }
func i64ptr(i int64) *int64   { return &i }

// TestUnitFloatingIpToStateFull verifies every flattened attribute is mapped
// from a fully-populated normalised floating IP.
func TestUnitFloatingIpToStateFull(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updated := time.Date(2026, 6, 7, 8, 9, 10, 0, time.UTC)

	fip := networkfloatingip.FloatingIp{
		Id:                  i64ptr(42),
		ExternalId:          strptr("ext-1"),
		IpStatus:            strptr("assigned"),
		IpAddress:           strptr("10.0.0.5"),
		PtrId:               strptr("ptr-9"),
		CloudId:             i64ptr(7),
		CloudName:           strptr("prod-cloud"),
		ServerId:            i64ptr(11),
		ServerName:          strptr("web-01"),
		IpRangeStartAddress: strptr("10.0.0.1"),
		IpRangeEndAddress:   strptr("10.0.0.254"),
		NetworkDomainId:     i64ptr(3),
		NetworkDomainName:   strptr("example.local"),
		DateCreated:         &created,
		LastUpdated:         &updated,
	}

	state := networkfloatingip.FloatingIpToState(fip)

	if state.Id.ValueInt64() != 42 {
		t.Errorf("Id = %d, want 42", state.Id.ValueInt64())
	}

	if state.ExternalId.ValueString() != "ext-1" {
		t.Errorf("ExternalId = %q, want ext-1", state.ExternalId.ValueString())
	}

	if state.IpStatus.ValueString() != "assigned" {
		t.Errorf("IpStatus = %q, want assigned", state.IpStatus.ValueString())
	}

	if state.IpAddress.ValueString() != "10.0.0.5" {
		t.Errorf("IpAddress = %q, want 10.0.0.5", state.IpAddress.ValueString())
	}

	if state.PtrId.ValueString() != "ptr-9" {
		t.Errorf("PtrId = %q, want ptr-9", state.PtrId.ValueString())
	}

	if state.CloudId.ValueInt64() != 7 {
		t.Errorf("CloudId = %d, want 7", state.CloudId.ValueInt64())
	}

	if state.CloudName.ValueString() != "prod-cloud" {
		t.Errorf("CloudName = %q, want prod-cloud", state.CloudName.ValueString())
	}

	if state.ServerId.ValueInt64() != 11 {
		t.Errorf("ServerId = %d, want 11", state.ServerId.ValueInt64())
	}

	if state.ServerName.ValueString() != "web-01" {
		t.Errorf("ServerName = %q, want web-01", state.ServerName.ValueString())
	}

	if state.IpRangeStartAddress.ValueString() != "10.0.0.1" {
		t.Errorf("IpRangeStartAddress = %q, want 10.0.0.1", state.IpRangeStartAddress.ValueString())
	}

	if state.IpRangeEndAddress.ValueString() != "10.0.0.254" {
		t.Errorf("IpRangeEndAddress = %q, want 10.0.0.254", state.IpRangeEndAddress.ValueString())
	}

	if state.NetworkDomainId.ValueInt64() != 3 {
		t.Errorf("NetworkDomainId = %d, want 3", state.NetworkDomainId.ValueInt64())
	}

	if state.NetworkDomainName.ValueString() != "example.local" {
		t.Errorf("NetworkDomainName = %q, want example.local", state.NetworkDomainName.ValueString())
	}

	if state.DateCreated.IsNull() || state.DateCreated.ValueString() == "" {
		t.Errorf("DateCreated should be set, got %q", state.DateCreated.ValueString())
	}

	if state.LastUpdated.IsNull() || state.LastUpdated.ValueString() == "" {
		t.Errorf("LastUpdated should be set, got %q", state.LastUpdated.ValueString())
	}
}

// TestUnitFloatingIpToStateEmpty verifies that absent nested objects map to null
// without panicking.
func TestUnitFloatingIpToStateEmpty(t *testing.T) {
	t.Parallel()

	state := networkfloatingip.FloatingIpToState(networkfloatingip.FloatingIp{})

	if !state.CloudId.IsNull() {
		t.Error("CloudId should be null when unset")
	}

	if !state.ServerName.IsNull() {
		t.Error("ServerName should be null when unset")
	}

	if !state.NetworkDomainId.IsNull() {
		t.Error("NetworkDomainId should be null when unset")
	}

	if !state.DateCreated.IsNull() {
		t.Error("DateCreated should be null when unset")
	}
}

// TestUnitFloatingIpFromGet verifies that nested SDK objects are flattened and
// nil nested objects are tolerated.
func TestUnitFloatingIpFromGet(t *testing.T) {
	t.Parallel()

	fip := &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIp{
		Id:        i64ptr(5),
		IpStatus:  strptr("free"),
		IpAddress: strptr("192.168.1.10"),
		Cloud: &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIpCloud{
			Id:   i64ptr(2),
			Name: strptr("dev-cloud"),
		},
		IpRange: &sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIpIpRange{
			StartAddress: strptr("192.168.1.1"),
			EndAddress:   strptr("192.168.1.100"),
		},
	}

	got := networkfloatingip.FloatingIpFromGet(fip)

	if got.Id == nil || *got.Id != 5 {
		t.Error("Id not mapped")
	}

	if got.CloudName == nil || *got.CloudName != "dev-cloud" {
		t.Error("CloudName not mapped")
	}

	if got.IpRangeEndAddress == nil || *got.IpRangeEndAddress != "192.168.1.100" {
		t.Error("IpRangeEndAddress not mapped")
	}

	// Server and NetworkDomain were nil; they must remain unset.
	if got.ServerId != nil || got.NetworkDomainId != nil {
		t.Error("absent nested objects should remain nil")
	}
}
