// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

func getNetworkAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (NetworkModel, diag.Diagnostics) {
	var state NetworkModel
	var diags diag.Diagnostics

	network, hresp, err := client.NetworksAPI.GetNetwork(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate network resource",
			fmt.Sprintf("network %d GET failed: ", id)+
				errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	net := network.GetNetwork()

	state.Id = convert.Int64ToType(net.Id)
	state.Name = convert.StrToType(net.Name)
	if net.DisplayName.IsSet() {
		state.DisplayName = convert.StrToType(net.DisplayName.Get())
	}
	if net.Description.IsSet() {
		state.Description = convert.StrToType(net.Description.Get())
	}
	state.Active = convert.BoolToType(net.Active)
	state.AllowStaticOverride = convert.BoolToType(
		net.AllowStaticOverride,
	)
	state.ApplianceUrlProxyBypass = convert.BoolToType(
		net.ApplianceUrlProxyBypass,
	)
	state.AssignPublicIp = convert.BoolToType(net.AssignPublicIp)
	if net.Cidr.IsSet() {
		state.Cidr = convert.StrToType(net.Cidr.Get())
	}
	if net.CidrIPv6.IsSet() {
		state.CidrIpv6 = convert.StrToType(net.CidrIPv6.Get())
	}
	state.DhcpServer = convert.BoolToType(net.DhcpServer)
	state.DhcpServerIpv6 = convert.BoolToType(net.DhcpServerIPv6)
	if net.DnsPrimary.IsSet() {
		state.DnsPrimary = convert.StrToType(net.DnsPrimary.Get())
	}
	if net.DnsPrimaryIPv6.IsSet() {
		state.DnsPrimaryIpv6 = convert.StrToType(
			net.DnsPrimaryIPv6.Get(),
		)
	}
	if net.DnsSecondary.IsSet() {
		state.DnsSecondary = convert.StrToType(net.DnsSecondary.Get())
	}
	if net.DnsSecondaryIPv6.IsSet() {
		state.DnsSecondaryIpv6 = convert.StrToType(
			net.DnsSecondaryIPv6.Get(),
		)
	}
	if net.Gateway.IsSet() {
		state.Gateway = convert.StrToType(net.Gateway.Get())
	}
	if net.GatewayIPv6.IsSet() {
		state.GatewayIpv6 = convert.StrToType(net.GatewayIPv6.Get())
	}
	state.Ipv4enabled = convert.BoolToType(net.Ipv4Enabled)
	state.Ipv6enabled = convert.BoolToType(net.Ipv6Enabled)
	if net.NetmaskIPv6.IsSet() {
		state.NetmaskIpv6 = convert.StrToType(net.NetmaskIPv6.Get())
	}
	if net.NoProxy.IsSet() {
		state.NoProxy = convert.StrToType(net.NoProxy.Get())
	}
	if net.SearchDomains.IsSet() {
		state.SearchDomains = convert.StrToType(
			net.SearchDomains.Get(),
		)
	}

	if net.Pool != nil && net.Pool.Id != nil {
		state.PoolId = convert.Int64ToType(net.Pool.Id)
	} else {
		state.PoolId = types.Int64Null()
	}

	if net.PoolIPv6 != nil && net.PoolIPv6.Id != nil {
		state.PoolIpv6Id = convert.Int64ToType(net.PoolIPv6.Id)
	} else {
		state.PoolIpv6Id = types.Int64Null()
	}

	if net.ZonePool != nil && net.ZonePool.Id != nil {
		state.ZonePoolId = convert.Int64ToType(net.ZonePool.Id)
	} else {
		state.ZonePoolId = types.Int64Null()
	}

	if net.VlanId.IsSet() {
		state.VlanId = convert.Int64ToType(net.VlanId.Get())
	}

	state.Labels = types.SetNull(types.StringType)
	if net.Labels != nil {
		var labelValues []attr.Value
		for _, label := range net.Labels {
			labelValues = append(labelValues, types.StringValue(label))
		}

		if len(labelValues) > 0 {
			labelsSet, d := types.SetValue(types.StringType, labelValues)
			diags.Append(d...)
			if diags.HasError() {
				return state, diags
			}
			state.Labels = labelsSet
		}
	}

	state.Config = types.DynamicNull()

	if net.NetworkDomain != nil {
		state.NetworkDomainId = convert.Int64ToType(
			net.NetworkDomain.Id,
		)
	} else {
		state.NetworkDomainId = types.Int64Null()
	}

	if net.NetworkProxy != nil {
		state.NetworkProxyId = convert.Int64ToType(
			net.NetworkProxy.Id,
		)
	} else {
		state.NetworkProxyId = types.Int64Null()
	}

	if net.Zone != nil {
		state.CloudId = convert.Int64ToType(net.Zone.Id)
	} else {
		state.CloudId = types.Int64Null()
	}

	group, ok := net.GetGroupOk()
	if ok && group.Id != nil {
		state.GroupId = convert.Int64ToType(group.Id)
	} else {
		state.GroupId = types.Int64Null()
	}

	if net.Type != nil {
		state.TypeId = convert.Int64ToType(net.Type.Id)
	} else {
		state.TypeId = types.Int64Null()
	}

	state.TenantIds = types.SetNull(types.Int64Type)
	if len(net.Tenants) > 0 {
		var tenantValues []attr.Value
		for _, tenant := range net.Tenants {
			if tenant.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*tenant.Id))
			}
		}
		if len(tenantValues) > 0 {
			tenantSet, d := types.SetValue(
				types.Int64Type, tenantValues,
			)
			diags.Append(d...)
			if diags.HasError() {
				return state, diags
			}
			state.TenantIds = tenantSet
		}
	}

	state.Visibility = convert.StrToType(net.Visibility)

	resourcePermission, ok := net.GetResourcePermissionOk()
	if ok {
		resourcePermissions, d := convertResourcePermissions(ctx, resourcePermission)
		diags.Append(d...)
		if diags.HasError() {
			return state, diags
		}
		state.ResourcePermissions = resourcePermissions
	} else {
		state.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return state, diags
}

func convertResourcePermissions(
	ctx context.Context,
	resourcePermission *sdk.ListNetworks200ResponseAllOfNetworksInnerResourcePermission,
) (ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupValues []attr.Value
	sites, ok := resourcePermission.GetSitesOk()
	if ok {
		for _, site := range sites {
			if site.Id != nil {
				groupValues = append(
					groupValues, types.Int64Value(*site.Id),
				)
			}
		}
	}

	var groupIDsSet attr.Value
	if len(groupValues) > 0 {
		groupIDsSet, _ = types.SetValue(types.Int64Type, groupValues)
	} else {
		groupIDsSet = types.SetNull(types.Int64Type)
	}

	resourcePermissions, d := NewResourcePermissionsValue(
		ResourcePermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"all": types.BoolValue(
				resourcePermission.All != nil &&
					*resourcePermission.All,
			),
			"group_ids": groupIDsSet,
		},
	)
	diags.Append(d...)

	return resourcePermissions, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	state, diags := getNetworkAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
