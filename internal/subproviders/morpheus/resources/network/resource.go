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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkResourceSchema(ctx)
}

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

	if net.Labels != nil {
		var labelValues []attr.Value
		for _, label := range net.Labels {
			labelValues = append(labelValues, types.StringValue(label))
		}

		labelsSet, d := types.SetValue(types.StringType, labelValues)
		diags.Append(d...)
		if diags.HasError() {
			return state, diags
		}
		state.Labels = labelsSet
	} else {
		state.Labels = types.SetNull(types.StringType)
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

	if group, ok := net.GetGroupOk(); ok && group.Id != nil {
		state.GroupId = convert.Int64ToType(group.Id)
	} else {
		state.GroupId = types.Int64Null()
	}

	if net.Type != nil {
		state.TypeId = convert.Int64ToType(net.Type.Id)
	} else {
		state.TypeId = types.Int64Null()
	}

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
		} else {
			state.TenantIds = types.SetNull(types.Int64Type)
		}
	} else {
		state.TenantIds = types.SetNull(types.Int64Type)
	}

	state.Visibility = convert.StrToType(net.Visibility)

	if resourcePermission, ok := net.GetResourcePermissionOk(); ok {
		var groupValues []attr.Value
		if sites, sitesOk := resourcePermission.GetSitesOk(); sitesOk {
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
		if diags.HasError() {
			return state, diags
		}
		state.ResourcePermissions = resourcePermissions
	} else {
		state.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	name := plan.Name.ValueString()

	createNetwork := sdk.NewCreateNetworksRequestNetworkWithDefaults()
	createNetwork.SetName(name)
	createNetwork.SetSite(*sdk.NewCreateNetworksRequestNetworkSite(
		plan.GroupId.ValueInt64(),
	))
	createNetwork.SetZone(*sdk.NewCreateNetworksRequestNetworkZone(
		plan.CloudId.ValueInt64(),
	))

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createNetwork.SetDescription(plan.Description.ValueString())
	}

	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
		createNetwork.SetDisplayName(plan.DisplayName.ValueString())
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		createNetwork.SetActive(plan.Active.ValueBool())
	}

	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		createNetwork.SetCidr(plan.Cidr.ValueString())
	}

	if !plan.CidrIpv6.IsNull() && !plan.CidrIpv6.IsUnknown() {
		createNetwork.SetCidrIPv6(plan.CidrIpv6.ValueString())
	}

	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		createNetwork.SetGateway(plan.Gateway.ValueString())
	}

	if !plan.GatewayIpv6.IsNull() && !plan.GatewayIpv6.IsUnknown() {
		createNetwork.SetGatewayIPv6(plan.GatewayIpv6.ValueString())
	}

	if !plan.DnsPrimary.IsNull() && !plan.DnsPrimary.IsUnknown() {
		createNetwork.SetDnsPrimary(plan.DnsPrimary.ValueString())
	}

	if !plan.DnsSecondary.IsNull() && !plan.DnsSecondary.IsUnknown() {
		createNetwork.SetDnsSecondary(plan.DnsSecondary.ValueString())
	}

	if !plan.DnsPrimaryIpv6.IsNull() && !plan.DnsPrimaryIpv6.IsUnknown() {
		createNetwork.SetDnsPrimaryIPv6(plan.DnsPrimaryIpv6.ValueString())
	}

	if !plan.DnsSecondaryIpv6.IsNull() &&
		!plan.DnsSecondaryIpv6.IsUnknown() {
		createNetwork.SetDnsSecondaryIPv6(plan.DnsSecondaryIpv6.ValueString())
	}

	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		createNetwork.SetDhcpServer(plan.DhcpServer.ValueBool())
	}

	if !plan.DhcpServerIpv6.IsNull() &&
		!plan.DhcpServerIpv6.IsUnknown() {
		createNetwork.SetDhcpServerIPv6(plan.DhcpServerIpv6.ValueBool())
	}

	if !plan.AllowStaticOverride.IsNull() &&
		!plan.AllowStaticOverride.IsUnknown() {
		createNetwork.SetAllowStaticOverride(
			plan.AllowStaticOverride.ValueBool(),
		)
	}

	if !plan.AssignPublicIp.IsNull() &&
		!plan.AssignPublicIp.IsUnknown() {
		createNetwork.SetAssignPublicIp(plan.AssignPublicIp.ValueBool())
	}

	if !plan.ApplianceUrlProxyBypass.IsNull() &&
		!plan.ApplianceUrlProxyBypass.IsUnknown() {
		createNetwork.SetApplianceUrlProxyBypass(
			plan.ApplianceUrlProxyBypass.ValueBool(),
		)
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		createNetwork.SetVisibility(plan.Visibility.ValueString())
	}

	if !plan.VlanId.IsNull() && !plan.VlanId.IsUnknown() {
		createNetwork.SetVlanId(plan.VlanId.ValueInt64())
	}

	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		createNetwork.SetPool(plan.PoolId.ValueInt64())
	}

	if !plan.PoolIpv6Id.IsNull() && !plan.PoolIpv6Id.IsUnknown() {
		createNetwork.SetPoolIPv6(plan.PoolIpv6Id.ValueInt64())
	}

	if !plan.ZonePoolId.IsNull() && !plan.ZonePoolId.IsUnknown() {
		zonePool := sdk.NewCreateNetworksRequestNetworkZonePool()
		zonePool.SetId(plan.ZonePoolId.ValueInt64())
		createNetwork.SetZonePool(*zonePool)
	}

	if !plan.Ipv4enabled.IsNull() && !plan.Ipv4enabled.IsUnknown() {
		createNetwork.SetIpv4Enabled(plan.Ipv4enabled.ValueBool())
	}

	if !plan.Ipv6enabled.IsNull() && !plan.Ipv6enabled.IsUnknown() {
		createNetwork.SetIpv6Enabled(plan.Ipv6enabled.ValueBool())
	}

	if !plan.NetmaskIpv6.IsNull() && !plan.NetmaskIpv6.IsUnknown() {
		createNetwork.SetNetmaskIPv6(plan.NetmaskIpv6.ValueString())
	}

	if !plan.NoProxy.IsNull() && !plan.NoProxy.IsUnknown() {
		createNetwork.SetNoProxy(plan.NoProxy.ValueString())
	}

	if !plan.SearchDomains.IsNull() && !plan.SearchDomains.IsUnknown() {
		createNetwork.SetSearchDomains(plan.SearchDomains.ValueString())
	}

	if !plan.TypeId.IsNull() && !plan.TypeId.IsUnknown() {
		networkType := sdk.NewCreateNetworksRequestNetworkType(
			plan.TypeId.ValueInt64(),
		)
		createNetwork.SetType(*networkType)
	}

	if !plan.NetworkDomainId.IsNull() &&
		!plan.NetworkDomainId.IsUnknown() {
		networkDomain := sdk.
			NewListNetworks200ResponseAllOfNetworksInnerNetworkDomain()
		networkDomain.SetId(plan.NetworkDomainId.ValueInt64())
		createNetwork.SetNetworkDomain(*networkDomain)
	}

	if !plan.NetworkProxyId.IsNull() &&
		!plan.NetworkProxyId.IsUnknown() {
		networkProxy := sdk.
			NewListNetworks200ResponseAllOfNetworksInnerNetworkProxy()
		networkProxy.SetId(plan.NetworkProxyId.ValueInt64())
		createNetwork.SetNetworkProxy(*networkProxy)
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": failed to parse labels: "+
					err.Error(),
			)

			return
		}
		createNetwork.SetLabels(labels)
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": failed to convert config: "+
					err.Error(),
			)

			return
		}

		if configDataMap, ok := configMap.(map[string]any); ok {
			networkConfig := sdk.CreateNetworksRequestNetworkConfig{}
			networkConfig.MapmapOfStringAny = &configDataMap
			createNetwork.SetConfig(networkConfig)
			tflog.Debug(ctx, fmt.Sprintf(
				"Config set with %d properties", len(configDataMap),
			))
		} else {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": config must be a valid object/map",
			)

			return
		}
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenants []sdk.GetAlerts200ResponseAllOfChecksInnerAccount
		for _, elem := range plan.TenantIds.Elements() {
			if idVal, ok := elem.(types.Int64); ok && !idVal.IsNull() {
				tenant := sdk.
					GetAlerts200ResponseAllOfChecksInnerAccount{}
				tenant.SetId(idVal.ValueInt64())
				tenants = append(tenants, tenant)
			}
		}
		if len(tenants) > 0 {
			createNetwork.SetTenants(tenants)
		}
	}

	if !plan.ResourcePermissions.IsNull() &&
		!plan.ResourcePermissions.IsUnknown() {
		resourcePermission := sdk.
			NewCreateNetworksRequestNetworkResourcePermission()

		allValue := plan.ResourcePermissions.All.ValueBool()
		resourcePermission.SetAll(allValue)

		if !plan.ResourcePermissions.GroupIds.IsNull() &&
			!plan.ResourcePermissions.GroupIds.IsUnknown() {
			var sites []int64
			for _, elem := range plan.ResourcePermissions.GroupIds.Elements() {
				if idVal, ok := elem.(types.Int64); ok && !idVal.IsNull() {
					sites = append(sites, idVal.ValueInt64())
				}
			}
			if len(sites) > 0 {
				resourcePermission.SetSites(sites)
			}
		}

		createNetwork.SetResourcePermission(*resourcePermission)
	}

	createNetworkReq := sdk.NewCreateNetworksRequest()
	createNetworkReq.SetNetwork(*createNetwork)

	tflog.Debug(ctx, fmt.Sprintf("Creating network '%s'", name))

	network, hresp, err := client.NetworksAPI.CreateNetworks(ctx).
		CreateNetworksRequest(*createNetworkReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create network resource",
			fmt.Sprintf("network %s POST failed: %s",
				name, errors.ErrMsg(err, hresp)),
		)

		return
	}

	if network.GetNetwork().Id == nil {
		resp.Diagnostics.AddError(
			"create network resource",
			"network "+name+": id is nil",
		)

		return
	}

	id := *network.GetNetwork().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getNetworkAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create network resource",
			fmt.Sprintf("network %d: failed to read from api", id),
		)

		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	// Create a custom timeout for the delete operation
	// Mainly needed for GCP network delete, which is
	// synchronous, and can take some time
	deleteCtx, cancel := context.WithTimeout(ctx, constants.NetworkDeleteTimeout)
	defer cancel()

	tflog.Debug(ctx, fmt.Sprintf("Deleting network %d", id))
	_, hresp, err := client.NetworksAPI.DeleteNetwork(deleteCtx, id).
		Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete network resource",
			fmt.Sprintf("network %d DELETE failed: %s",
				id, errors.ErrMsg(err, hresp)),
		)
	}
}

func (r *Resource) ImportState(
	_ context.Context,
	_ resource.ImportStateRequest,
	_ *resource.ImportStateResponse,
) {
}
