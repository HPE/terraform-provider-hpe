// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hosts

import (
	"context"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func serverWithType(code string, _ bool) *sdk.ListHosts200ResponseAllOfServersInner {
	return &sdk.ListHosts200ResponseAllOfServersInner{
		Id:   sdk.PtrInt64(1),
		Name: sdk.PtrString("h"),
		ComputeServerType: &sdk.ListHosts200ResponseAllOfServersInnerComputeServerType{
			Code: sdk.PtrString(code),
		},
	}
}

func TestUnitMatchesClientFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		srv          *sdk.ListHosts200ResponseAllOfServersInner
		wantType     string
		filterByType bool
		want         bool
	}{
		{
			name: "no filters passes",
			srv:  serverWithType("vmwareKvm", false),
			want: true,
		},
		{
			name:         "type matches",
			srv:          serverWithType("omega.baremetal.stub-server", true),
			wantType:     "omega.baremetal.stub-server",
			filterByType: true,
			want:         true,
		},
		{
			name:         "type mismatch",
			srv:          serverWithType("vmwareKvm", false),
			wantType:     "omega.baremetal.stub-server",
			filterByType: true,
			want:         false,
		},
		{
			name:         "missing compute server type with type filter excludes",
			srv:          &sdk.ListHosts200ResponseAllOfServersInner{Id: sdk.PtrInt64(1)},
			wantType:     "vmwareKvm",
			filterByType: true,
			want:         false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := matchesClientFilters(tc.srv, tc.wantType, tc.filterByType)
			if got != tc.want {
				t.Errorf("matchesClientFilters = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUnitHostToValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	srv := &sdk.ListHosts200ResponseAllOfServersInner{
		Id:         sdk.PtrInt64(10),
		Name:       sdk.PtrString("metal-10"),
		Status:     sdk.PtrString("provisioned"),
		PowerState: sdk.PtrString("on"),
		ExternalIp: *sdk.NewNullableString(sdk.PtrString("203.0.113.10")),
		InternalIp: *sdk.NewNullableString(sdk.PtrString("10.0.0.10")),
		ZoneId:     sdk.PtrInt64(17),
		ComputeServerType: &sdk.ListHosts200ResponseAllOfServersInnerComputeServerType{
			Code: sdk.PtrString("omega.baremetal.stub-server"),
		},
		ParentServer: *sdk.NewNullableListHosts200ResponseAllOfServersInnerParentServer(
			&sdk.ListHosts200ResponseAllOfServersInnerParentServer{Id: sdk.PtrInt64(3)},
		),
	}

	v, diags := hostToValue(ctx, srv, map[string]bool{"omega.baremetal.stub-server": true})
	if diags.HasError() {
		t.Fatalf("hostToValue: %v", diags.Errors())
	}

	if got := v.Id.ValueInt64(); got != 10 {
		t.Errorf("id = %d", got)
	}
	if got := v.ComputeServerType.ValueString(); got != "omega.baremetal.stub-server" {
		t.Errorf("compute_server_type = %q", got)
	}
	if !v.BareMetal.ValueBool() {
		t.Error("bare_metal = false, want true")
	}
	if got := v.CloudId.ValueInt64(); got != 17 {
		t.Errorf("cloud_id = %d", got)
	}
	if got := v.ParentHostId.ValueInt64(); got != 3 {
		t.Errorf("parent_host_id = %d", got)
	}
}
