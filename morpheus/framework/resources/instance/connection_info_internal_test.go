// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestGetConnectionInfo guards MORPH-17158. connectionInfo is an optional field on
// the instance response (JSON omitempty), so it is absent for instances with no
// connection addresses (stopped / failed / not-yet-provisioned, or types that do
// not populate it). An absent (nil), empty, or IP-less connectionInfo must map to a
// null connection_info list WITHOUT a diagnostic, so refresh and import of such an
// instance succeed instead of failing with a manufactured "GET connectionInfo
// failed" error. A populated connectionInfo maps to the list of IP addresses.
func TestGetConnectionInfo(t *testing.T) {
	t.Parallel()

	ci := func(ips ...*string) []sdk.AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner {
		out := make([]sdk.AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner, 0, len(ips))
		for _, ip := range ips {
			out = append(out, sdk.AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner{Ip: ip})
		}

		return out
	}
	strp := func(s string) *string { return &s }
	list := func(ips ...string) types.List {
		vals := make([]attr.Value, 0, len(ips))
		for _, ip := range ips {
			vals = append(vals, types.StringValue(ip))
		}

		return types.ListValueMust(types.StringType, vals)
	}

	tests := map[string]struct {
		connInfo []sdk.AddInstance200ResponseAllOfOneOfInstanceConnectionInfoInner
		want     types.List
	}{
		"absent (nil) connectionInfo is null, not an error": {
			connInfo: nil,
			want:     types.ListNull(types.StringType),
		},
		"empty connectionInfo is null": {
			connInfo: ci(),
			want:     types.ListNull(types.StringType),
		},
		"entries without an IP are null": {
			connInfo: ci(nil, nil),
			want:     types.ListNull(types.StringType),
		},
		"populated connectionInfo maps to the IP list": {
			connInfo: ci(strp("10.0.0.1"), strp("10.0.0.2")),
			want:     list("10.0.0.1", "10.0.0.2"),
		},
		"only entries carrying an IP are included": {
			connInfo: ci(strp("10.0.0.1"), nil, strp("10.0.0.3")),
			want:     list("10.0.0.1", "10.0.0.3"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, diags := getConnectionInfo(sdk.GetInstance200ResponseInstance{
				ConnectionInfo: tc.connInfo,
			})

			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags.Errors())
			}

			if !got.Equal(tc.want) {
				t.Errorf("connection_info = %s, want %s", got, tc.want)
			}
		})
	}
}
