// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingips

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strptr(s string) *string { return &s }

// TestUnitFieldValue verifies the named-field extraction used by the client-side
// filter, including the "field absent" case which must report no value.
func TestUnitFieldValue(t *testing.T) {
	t.Parallel()

	item := &sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner{
		IpAddress: strptr("10.1.1.1"),
		IpStatus:  strptr("assigned"),
		Cloud: &sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInnerCloud{
			Name: strptr("prod-cloud"),
		},
		Server: &sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInnerServer{
			Name: strptr("web-01"),
		},
	}
	item.ExternalId.Set(strptr("ext-abc"))

	cases := []struct {
		field  string
		want   string
		wantOK bool
	}{
		{"ip_address", "10.1.1.1", true},
		{"ip_status", "assigned", true},
		{"external_id", "ext-abc", true},
		{"cloud_name", "prod-cloud", true},
		{"server_name", "web-01", true},
		{"unknown_field", "", false},
	}

	for _, c := range cases {
		got, ok := fieldValue(item, c.field)
		if got != c.want || ok != c.wantOK {
			t.Errorf("fieldValue(%q) = (%q, %v), want (%q, %v)",
				c.field, got, ok, c.want, c.wantOK)
		}
	}
}

// TestUnitFieldValueAbsentNested verifies that a nil nested object reports no
// value rather than panicking.
func TestUnitFieldValueAbsentNested(t *testing.T) {
	t.Parallel()

	item := &sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner{}

	if _, ok := fieldValue(item, "cloud_name"); ok {
		t.Error("cloud_name should report no value when cloud is nil")
	}

	if _, ok := fieldValue(item, "server_name"); ok {
		t.Error("server_name should report no value when server is nil")
	}

	if _, ok := fieldValue(item, "external_id"); ok {
		t.Error("external_id should report no value when unset")
	}
}
