// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

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
		} else {
			resp.Diagnostics.AddError(
				"create network resource",
				"network "+name+": config must be a valid object/map",
			)

			return
		}
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []types.Int64
		diags := plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var tenants []sdk.GetAlerts200ResponseAllOfChecksInnerAccount
		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
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
			var groupIDs []types.Int64
			diags := plan.ResourcePermissions.GroupIds.ElementsAs(ctx,
				&groupIDs, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			var sites []int64
			for _, idVal := range groupIDs {
				if !idVal.IsNull() {
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
