// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// buildCreateRequest assembles the AddBaremetalHostRequest from the plan and the
// write-only config values read from configuration. The typed config_* blocks and
// the free-form config map are mutually exclusive (enforced in the schema).
//
// configBmaasPwd, configOmegaPwd and sshPwd are the write-only password values
// read from the configuration (never from plan/state, where they are null).
func buildCreateRequest(
	ctx context.Context,
	plan HostModel,
	sshPwd *string,
	configBmaasPwd *string,
	configOmegaPwd *string,
) (sdk.AddBaremetalHostRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	server := sdk.AddBaremetalHostRequestServer{
		Name: plan.Name.ValueStringPointer(),
		Cloud: &sdk.AddBaremetalHostRequestServerCloud{
			Id: plan.CloudId.ValueInt64Pointer(),
		},
		Group: &sdk.AddBaremetalHostRequestServerGroup{
			Id: plan.GroupId.ValueInt64Pointer(),
		},
		ComputeServerType: &sdk.AddBaremetalHostRequestServerComputeServerType{
			AddBaremetalHostRequestServerComputeServerTypeAnyOf: &sdk.AddBaremetalHostRequestServerComputeServerTypeAnyOf{
				Code: plan.ComputeServerType.ValueStringPointer(),
			},
		},
	}

	if !plan.PlanId.IsNull() && !plan.PlanId.IsUnknown() {
		server.Plan = &sdk.AddBaremetalHostRequestServerPlan{
			Id: plan.PlanId.ValueInt64Pointer(),
		}
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		server.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		server.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.SshUsername.IsNull() && !plan.SshUsername.IsUnknown() {
		server.SshUsername = plan.SshUsername.ValueStringPointer()
	}

	if sshPwd != nil {
		server.SshPassword = sshPwd
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		diags.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		server.Labels = labels
	}

	// config.
	cfg, cfgDiags := buildConfig(ctx, plan, configBmaasPwd, configOmegaPwd)
	diags.Append(cfgDiags...)
	if diags.HasError() {
		return sdk.AddBaremetalHostRequest{}, diags
	}

	server.Config = cfg

	req := sdk.AddBaremetalHostRequest{
		Server: &server,
	}

	// volumes.
	vols, vDiags := buildVolumes(ctx, plan)
	diags.Append(vDiags...)
	req.Volumes = vols

	// network interfaces.
	nics, nDiags := buildNetworkInterfaces(ctx, plan)
	diags.Append(nDiags...)
	req.NetworkInterfaces = nics

	return req, diags
}

// buildConfig routes the plan's config to the matching SDK config union variant.
func buildConfig(
	ctx context.Context,
	plan HostModel,
	configBmaasPwd *string,
	configOmegaPwd *string,
) (*sdk.AddBaremetalHostRequestServerConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The add-host handler binds server.config wholesale, so the keys every
	// type shares ride in whichever variant is configured: userData, and
	// provisionKey, which is how the handler assigns the SSH key pair.
	extras := commonConfig(plan)

	switch {
	case !plan.ConfigBmaas.IsNull() && !plan.ConfigBmaas.IsUnknown():
		b := plan.ConfigBmaas
		obj := sdk.BmaasHostConfigObject{
			IloIpAddress:         b.IloIpAddress.ValueString(),
			ResourcePoolId:       b.ResourcePoolId.ValueString(),
			AdditionalProperties: extras,
		}
		if !b.IloUsername.IsNull() {
			obj.IloUsername.Set(b.IloUsername.ValueStringPointer())
		}
		if configBmaasPwd != nil {
			obj.IloPassword.Set(configBmaasPwd)
		}
		if !b.IloMacAddress.IsNull() {
			obj.IloMacAddress.Set(b.IloMacAddress.ValueStringPointer())
		}
		if !b.PreProvisioned.IsNull() {
			obj.PreProvisioned.Set(b.PreProvisioned.ValueBoolPointer())
		}

		return &sdk.AddBaremetalHostRequestServerConfig{
			BmaasHostConfigObject: &obj,
		}, diags

	case !plan.ConfigOmegaBmaas.IsNull() && !plan.ConfigOmegaBmaas.IsUnknown():
		o := plan.ConfigOmegaBmaas
		obj := sdk.OmegaBmaasHostConfigObject{
			ResourcePoolId:       o.ResourcePoolId.ValueString(),
			AdditionalProperties: extras,
		}
		if !o.ConsoleHost.IsNull() {
			obj.ConsoleHost.Set(o.ConsoleHost.ValueStringPointer())
		}
		if !o.ConsoleUsername.IsNull() {
			obj.ConsoleUsername.Set(o.ConsoleUsername.ValueStringPointer())
		}
		if configOmegaPwd != nil {
			obj.ConsolePassword.Set(configOmegaPwd)
		}
		if !o.NumNics.IsNull() {
			obj.NumNics.Set(o.NumNics.ValueInt64Pointer())
		}
		if !o.PreProvisioned.IsNull() {
			obj.PreProvisioned.Set(o.PreProvisioned.ValueBoolPointer())
		}

		return &sdk.AddBaremetalHostRequestServerConfig{
			OmegaBmaasHostConfigObject: &obj,
		}, diags

	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		generic, gDiags := genericConfig(ctx, plan)
		diags.Append(gDiags...)
		if diags.HasError() {
			return nil, diags
		}

		for k, v := range extras {
			generic.AdditionalProperties[k] = v
		}

		return &sdk.AddBaremetalHostRequestServerConfig{
			HostConfigGeneric: generic,
		}, diags

	default:
		// No config block set; still carry the shared keys if present.
		if len(extras) == 0 {
			return nil, diags
		}

		return &sdk.AddBaremetalHostRequestServerConfig{
			HostConfigGeneric: &sdk.HostConfigGeneric{AdditionalProperties: extras},
		}, diags
	}
}

// commonConfig returns the config keys shared by every host type that the
// plan sets: userData and provisionKey (the SSH key pair id). Nil when none.
func commonConfig(plan HostModel) map[string]interface{} {
	extras := map[string]interface{}{}

	if !plan.UserData.IsNull() && !plan.UserData.IsUnknown() {
		extras["userData"] = plan.UserData.ValueString()
	}

	if !plan.SshKeyPairId.IsNull() && !plan.SshKeyPairId.IsUnknown() {
		extras["provisionKey"] = plan.SshKeyPairId.ValueInt64()
	}

	if len(extras) == 0 {
		return nil
	}

	return extras
}

// buildVolumes converts the plan's volumes list to SDK volume objects.
func buildVolumes(
	ctx context.Context,
	plan HostModel,
) ([]sdk.AddBaremetalHostRequestVolumesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if plan.Volumes.IsNull() || plan.Volumes.IsUnknown() {
		return nil, diags
	}

	var vals []VolumesValue
	diags.Append(plan.Volumes.ElementsAs(ctx, &vals, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]sdk.AddBaremetalHostRequestVolumesInner, 0, len(vals))
	for _, v := range vals {
		vol := sdk.AddBaremetalHostRequestVolumesInner{}
		if !v.RootVolume.IsNull() {
			vol.RootVolume = v.RootVolume.ValueBoolPointer()
		}
		if !v.Name.IsNull() {
			vol.Name = v.Name.ValueStringPointer()
		}
		if !v.Size.IsNull() {
			vol.Size = v.Size.ValueInt64Pointer()
		}
		if !v.StorageTypeId.IsNull() {
			vol.StorageType = v.StorageTypeId.ValueInt64Pointer()
		}
		if !v.DatastoreId.IsNull() {
			vol.DatastoreId = &sdk.AddBaremetalHostRequestVolumesInnerDatastoreId{
				Int64: v.DatastoreId.ValueInt64Pointer(),
			}
		}
		out = append(out, vol)
	}

	return out, diags
}

// buildNetworkInterfaces converts the plan's network_interfaces list to SDK
// interface objects.
func buildNetworkInterfaces(
	ctx context.Context,
	plan HostModel,
) ([]sdk.AddBaremetalHostRequestNetworkInterfacesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if plan.NetworkInterfaces.IsNull() || plan.NetworkInterfaces.IsUnknown() {
		return nil, diags
	}

	var vals []NetworkInterfacesValue
	diags.Append(plan.NetworkInterfaces.ElementsAs(ctx, &vals, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]sdk.AddBaremetalHostRequestNetworkInterfacesInner, 0, len(vals))
	for _, v := range vals {
		nic := sdk.AddBaremetalHostRequestNetworkInterfacesInner{}
		if !v.NetworkId.IsNull() {
			nic.Network = &sdk.AddBaremetalHostRequestNetworkInterfacesInnerNetwork{
				Id: &sdk.AddBaremetalHostRequestNetworkInterfacesInnerNetworkId{
					Int64: v.NetworkId.ValueInt64Pointer(),
				},
			}
		}
		if !v.IpAddress.IsNull() {
			nic.IpAddress = v.IpAddress.ValueStringPointer()
		}
		if !v.IpMode.IsNull() {
			nic.IpMode = v.IpMode.ValueStringPointer()
		}
		if !v.NetworkInterfaceTypeId.IsNull() {
			nic.NetworkInterfaceTypeId = v.NetworkInterfaceTypeId.ValueInt64Pointer()
		}
		out = append(out, nic)
	}

	return out, diags
}
