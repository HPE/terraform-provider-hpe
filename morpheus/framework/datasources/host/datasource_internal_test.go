// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func TestUnitMapServer(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updated := time.Date(2026, 1, 2, 4, 5, 6, 0, time.UTC)

	srv := &sdk.GetHost200ResponseServer{
		Id:             sdk.PtrInt64(42),
		Uuid:           sdk.PtrString("uuid-42"),
		ExternalId:     *sdk.NewNullableString(sdk.PtrString("ext-42")),
		Name:           sdk.PtrString("metal-01"),
		Description:    *sdk.NewNullableString(sdk.PtrString("a bare metal host")),
		Status:         sdk.PtrString("provisioned"),
		PowerState:     sdk.PtrString("on"),
		InternalIp:     *sdk.NewNullableString(sdk.PtrString("10.0.0.5")),
		ExternalIp:     *sdk.NewNullableString(sdk.PtrString("203.0.113.5")),
		Hostname:       sdk.PtrString("metal-01.local"),
		Managed:        sdk.PtrBool(false),
		Enabled:        sdk.PtrBool(true),
		Visibility:     sdk.PtrString("private"),
		AgentInstalled: sdk.PtrBool(true),
		AgentVersion:   *sdk.NewNullableString(sdk.PtrString("7.0.0")),
		DateCreated:    &created,
		LastUpdated:    &updated,
		ZoneId:         sdk.PtrInt64(17),
		Labels:         []string{"prod", "metal"},
		ComputeServerType: &sdk.AddBaremetalHost200ResponseServerComputeServerType{
			Id:   sdk.PtrInt64(99),
			Code: sdk.PtrString("omega.baremetal.stub-server"),
			Name: sdk.PtrString("Omega Stub Server"),
		},
		Plan: &sdk.AddBaremetalHost200ResponseServerPlan{
			Id: *sdk.NewNullableInt64(sdk.PtrInt64(5)),
		},
		ParentServer: &sdk.AddBaremetalHost200ResponseServerParentServer{
			Id: sdk.PtrInt64(7),
		},
		Interfaces: []sdk.AddBaremetalHost200ResponseServerInterfacesInner{
			{IpAddress: sdk.PtrString("10.0.0.5")},
			{IpAddress: sdk.PtrString("")}, // skipped
			{IpAddress: sdk.PtrString("10.0.0.6")},
		},
		Volumes: []sdk.AddBaremetalHost200ResponseServerVolumesInner{
			{Name: sdk.PtrString("root")},
			{Name: sdk.PtrString("data")},
		},
	}

	var data HostModel
	mapServer(srv, &data)

	if got := data.Id.ValueInt64(); got != 42 {
		t.Errorf("id = %d, want 42", got)
	}
	if got := data.Name.ValueString(); got != "metal-01" {
		t.Errorf("name = %q", got)
	}
	if got := data.ComputeServerType.ValueString(); got != "omega.baremetal.stub-server" {
		t.Errorf("compute_server_type = %q", got)
	}
	if got := data.ComputeServerTypeId.ValueInt64(); got != 99 {
		t.Errorf("compute_server_type_id = %d", got)
	}
	if got := data.ComputeServerTypeName.ValueString(); got != "Omega Stub Server" {
		t.Errorf("compute_server_type_name = %q", got)
	}
	// bare_metal is not in the server's rendering; getHostByID derives it
	// from the server types list after mapServer, so it is null here.
	if !data.BareMetal.IsNull() {
		t.Errorf("bare_metal = %v, want null from mapServer alone", data.BareMetal)
	}
	if data.Managed.ValueBool() {
		t.Error("managed = true, want false")
	}
	if got := data.CloudId.ValueInt64(); got != 17 {
		t.Errorf("cloud_id = %d", got)
	}
	if got := data.PlanId.ValueInt64(); got != 5 {
		t.Errorf("plan_id = %d", got)
	}
	if got := data.ParentHostId.ValueInt64(); got != 7 {
		t.Errorf("parent_host_id = %d", got)
	}
	if got := data.InternalIp.ValueString(); got != "10.0.0.5" {
		t.Errorf("internal_ip = %q", got)
	}
	if got := data.DateCreated.ValueString(); got != "2026-01-02T03:04:05Z" {
		t.Errorf("date_created = %q", got)
	}

	labels := data.Labels.Elements()
	if len(labels) != 2 {
		t.Errorf("labels len = %d, want 2", len(labels))
	}

	ifaces := data.Interfaces.Elements()
	if len(ifaces) != 2 {
		t.Errorf("interfaces len = %d, want 2 (empty skipped)", len(ifaces))
	}

	vols := data.Volumes.Elements()
	if len(vols) != 2 {
		t.Errorf("volumes len = %d, want 2", len(vols))
	}
}

func TestUnitMapServerNullNested(t *testing.T) {
	t.Parallel()

	// A minimal server with no compute server type, plan, parent, interfaces or
	// volumes maps its nested fields to null without panicking.
	srv := &sdk.GetHost200ResponseServer{
		Id:   sdk.PtrInt64(1),
		Name: sdk.PtrString("minimal"),
	}

	var data HostModel
	mapServer(srv, &data)

	if !data.ComputeServerType.IsNull() {
		t.Error("compute_server_type should be null")
	}
	if !data.BareMetal.IsNull() {
		t.Error("bare_metal should be null")
	}
	if !data.PlanId.IsNull() {
		t.Error("plan_id should be null")
	}
	if !data.ParentHostId.IsNull() {
		t.Error("parent_host_id should be null")
	}
	if !data.Interfaces.IsNull() {
		t.Error("interfaces should be null")
	}
	if !data.Volumes.IsNull() {
		t.Error("volumes should be null")
	}
}
