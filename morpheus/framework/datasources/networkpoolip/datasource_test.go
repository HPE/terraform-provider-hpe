// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import (
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strPtr(s string) *string { return &s }
func intPtr(i int64) *int64   { return &i }
func boolPtr(b bool) *bool    { return &b }

// TestUnitNetworkPoolIPAsStateNil verifies that every nil pointer in the
// common field set becomes a Terraform null, exercising the nil-handling of
// the convert helpers used by the mapper.
func TestUnitNetworkPoolIPAsStateNil(t *testing.T) {
	t.Parallel()

	state := networkPoolIPAsState(networkPoolIPFields{})

	if !state.Id.IsNull() {
		t.Errorf("expected Id null, got %v", state.Id)
	}
	if !state.IpAddress.IsNull() {
		t.Errorf("expected IpAddress null, got %v", state.IpAddress)
	}
	if !state.GatewayAddress.IsNull() {
		t.Errorf("expected GatewayAddress null, got %v", state.GatewayAddress)
	}
	if !state.Active.IsNull() {
		t.Errorf("expected Active null, got %v", state.Active)
	}
	if !state.StartDate.IsNull() {
		t.Errorf("expected StartDate null, got %v", state.StartDate)
	}
	if !state.EndDate.IsNull() {
		t.Errorf("expected EndDate null, got %v", state.EndDate)
	}
	if !state.RefId.IsNull() {
		t.Errorf("expected RefId null, got %v", state.RefId)
	}
}

// TestUnitNetworkPoolIPAsStateValues verifies that populated fields, including
// RFC3339 date formatting, are mapped onto Terraform state correctly.
func TestUnitNetworkPoolIPAsStateValues(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	state := networkPoolIPAsState(networkPoolIPFields{
		Id:             intPtr(7),
		NetworkPoolId:  intPtr(99),
		IpAddress:      strPtr("10.0.0.5"),
		GatewayAddress: strPtr("10.0.0.1"),
		Hostname:       strPtr("host-a"),
		Active:         boolPtr(true),
		StaticIp:       boolPtr(false),
		RefId:          intPtr(42),
		DateCreated:    &created,
	})

	if state.Id.ValueInt64() != 7 {
		t.Errorf("expected Id 7, got %d", state.Id.ValueInt64())
	}
	if state.IpAddress.ValueString() != "10.0.0.5" {
		t.Errorf("expected IpAddress 10.0.0.5, got %s", state.IpAddress.ValueString())
	}
	if !state.Active.ValueBool() {
		t.Errorf("expected Active true")
	}
	if state.StaticIp.ValueBool() {
		t.Errorf("expected StaticIp false")
	}
	if state.RefId.ValueInt64() != 42 {
		t.Errorf("expected RefId 42, got %d", state.RefId.ValueInt64())
	}
	if got := state.DateCreated.ValueString(); got != "2026-01-02T03:04:05Z" {
		t.Errorf("expected RFC3339 date, got %s", got)
	}
}

// TestUnitFromSingleNullable verifies that unset nullable SDK fields on the
// single-get model flatten to nil pointers (and thus Terraform nulls).
func TestUnitFromSingleNullable(t *testing.T) {
	t.Parallel()

	ip := &sdk.GetNetworkPoolIp200ResponseNetworkPoolIp{
		Id:             intPtr(3),
		IpAddress:      strPtr("192.0.2.10"),
		GatewayAddress: *sdk.NewNullableString(strPtr("192.0.2.1")),
		DnsServer:      *sdk.NewNullableString(nil),
		RefId:          *sdk.NewNullableInt64(intPtr(5)),
	}

	f := fromSingle(ip)

	if f.Id == nil || *f.Id != 3 {
		t.Errorf("expected Id 3")
	}
	if f.GatewayAddress == nil || *f.GatewayAddress != "192.0.2.1" {
		t.Errorf("expected GatewayAddress set")
	}
	if f.DNSServer != nil {
		t.Errorf("expected DnsServer nil for unset nullable")
	}
	if f.RefId == nil || *f.RefId != 5 {
		t.Errorf("expected RefId 5")
	}

	state := networkPoolIPAsState(f)
	if !state.DnsServer.IsNull() {
		t.Errorf("expected DnsServer null in state")
	}
}

// TestUnitFromListItem verifies the list-item model flattens identically to
// the single-get model for the shared fields.
func TestUnitFromListItem(t *testing.T) {
	t.Parallel()

	item := &sdk.GetNetworkPoolIps200ResponseAllOfNetworkPoolIpsInner{
		Id:        intPtr(11),
		Hostname:  strPtr("host-b"),
		IpAddress: strPtr("192.0.2.20"),
	}

	f := fromListItem(item)

	if f.Id == nil || *f.Id != 11 {
		t.Errorf("expected Id 11")
	}
	if f.Hostname == nil || *f.Hostname != "host-b" {
		t.Errorf("expected Hostname host-b")
	}
	if f.IpAddress == nil || *f.IpAddress != "192.0.2.20" {
		t.Errorf("expected IpAddress 192.0.2.20")
	}
}
