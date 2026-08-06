// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makeConfig(t *testing.T, m instanceNodeModel) tfsdk.Config {
	t.Helper()

	ctx := context.Background()
	r := &Resource{}
	var schemaResp resource.SchemaResponse

	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	if m.Timeouts.IsNull() || m.Timeouts.IsUnknown() {
		m.Timeouts = timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"delete": types.StringType,
			}),
		}
	}

	// Config has no Set method; use Plan.Set then copy the raw value.
	p := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := p.Set(ctx, &m); diags.HasError() {
		t.Fatalf("build config: %v", diags)
	}

	return tfsdk.Config{Schema: schemaResp.Schema, Raw: p.Raw}
}

func TestUnitValidator_PreProvisionedTrueWithoutServerID(t *testing.T) {
	t.Parallel()

	cfg := makeConfig(t, instanceNodeModel{
		InstanceID:       types.Int64Value(1),
		ResourcePoolID:   types.Int64Null(),
		PreProvisioned:   types.BoolValue(true),
		SelectedServerID: types.Int64Null(),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	})

	v := preProvisionedRequiresServerID{}
	resp := &resource.ValidateConfigResponse{}

	v.ValidateResource(context.Background(), resource.ValidateConfigRequest{
		Config: cfg,
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("expected error when pre_provisioned=true without selected_server_id")
	}
}

func TestUnitValidator_PreProvisionedFalseAlone(t *testing.T) {
	t.Parallel()

	cfg := makeConfig(t, instanceNodeModel{
		InstanceID:       types.Int64Value(1),
		ResourcePoolID:   types.Int64Null(),
		PreProvisioned:   types.BoolValue(false),
		SelectedServerID: types.Int64Null(),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	})

	v := preProvisionedRequiresServerID{}
	resp := &resource.ValidateConfigResponse{}

	v.ValidateResource(context.Background(), resource.ValidateConfigRequest{
		Config: cfg,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error when pre_provisioned=false, got: %s",
			resp.Diagnostics.Errors())
	}
}

func TestUnitValidator_PreProvisionedAbsent(t *testing.T) {
	t.Parallel()

	cfg := makeConfig(t, instanceNodeModel{
		InstanceID:       types.Int64Value(1),
		ResourcePoolID:   types.Int64Null(),
		PreProvisioned:   types.BoolNull(),
		SelectedServerID: types.Int64Null(),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	})

	v := preProvisionedRequiresServerID{}
	resp := &resource.ValidateConfigResponse{}

	v.ValidateResource(context.Background(), resource.ValidateConfigRequest{
		Config: cfg,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error when pre_provisioned is absent, got: %s",
			resp.Diagnostics.Errors())
	}
}

func TestUnitValidator_PreProvisionedTrueWithServerID(t *testing.T) {
	t.Parallel()

	cfg := makeConfig(t, instanceNodeModel{
		InstanceID:       types.Int64Value(1),
		ResourcePoolID:   types.Int64Null(),
		PreProvisioned:   types.BoolValue(true),
		SelectedServerID: types.Int64Value(42),
		SshHost:          types.StringNull(),
		SshUsername:      types.StringNull(),
		SshPassword:      types.StringNull(),
		SshKeyPairID:     types.Int64Null(),
		WaitForIPAddress: types.BoolValue(false),
		ContainerID:      types.Int64Unknown(),
		ServerID:         types.Int64Unknown(),
		IPAddress:        types.StringUnknown(),
	})

	v := preProvisionedRequiresServerID{}
	resp := &resource.ValidateConfigResponse{}

	v.ValidateResource(context.Background(), resource.ValidateConfigRequest{
		Config: cfg,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error when pre_provisioned=true with selected_server_id, got: %s",
			resp.Diagnostics.Errors())
	}
}
