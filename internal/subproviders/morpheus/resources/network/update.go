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

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()
	name := plan.Name.ValueString()

	network := sdk.NewUpdateNetworkRequestNetwork()

	// Set all updateable fields from plan
	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
		network.SetDisplayName(plan.DisplayName.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		network.SetDescription(plan.Description.ValueString())
	}

	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		network.SetCidr(plan.Cidr.ValueString())
	}

	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		network.SetGateway(plan.Gateway.ValueString())
	}

	if !plan.DnsPrimary.IsNull() && !plan.DnsPrimary.IsUnknown() {
		network.SetDnsPrimary(plan.DnsPrimary.ValueString())
	}

	if !plan.DnsSecondary.IsNull() && !plan.DnsSecondary.IsUnknown() {
		network.SetDnsSecondary(plan.DnsSecondary.ValueString())
	}

	if !plan.VlanId.IsNull() && !plan.VlanId.IsUnknown() {
		network.SetVlanId(plan.VlanId.ValueInt64())
	}

	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		network.SetPool(plan.PoolId.ValueInt64())
	}

	if !plan.ZonePoolId.IsNull() && !plan.ZonePoolId.IsUnknown() {
		zonePool := sdk.NewCreateNetworksRequestNetworkZonePool()
		zonePool.SetId(plan.ZonePoolId.ValueInt64())
		network.SetZonePool(*zonePool)
	}

	if !plan.AllowStaticOverride.IsNull() && !plan.AllowStaticOverride.IsUnknown() {
		network.SetAllowStaticOverride(plan.AllowStaticOverride.ValueBool())
	}

	if !plan.AssignPublicIp.IsNull() && !plan.AssignPublicIp.IsUnknown() {
		network.SetAssignPublicIp(plan.AssignPublicIp.ValueBool())
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		network.SetActive(plan.Active.ValueBool())
	}

	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		network.SetDhcpServer(plan.DhcpServer.ValueBool())
	}

	if !plan.SearchDomains.IsNull() && !plan.SearchDomains.IsUnknown() {
		network.SetSearchDomains(plan.SearchDomains.ValueString())
	}

	if !plan.ApplianceUrlProxyBypass.IsNull() && !plan.ApplianceUrlProxyBypass.IsUnknown() {
		network.SetApplianceUrlProxyBypass(plan.ApplianceUrlProxyBypass.ValueBool())
	}

	if !plan.NoProxy.IsNull() && !plan.NoProxy.IsUnknown() {
		network.SetNoProxy(plan.NoProxy.ValueString())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		network.SetVisibility(plan.Visibility.ValueString())
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []types.String
		diags := plan.Labels.ElementsAs(ctx, &labels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var labelStrings []string
		for _, label := range labels {
			if !label.IsNull() {
				labelStrings = append(labelStrings, label.ValueString())
			}
		}
		network.SetLabels(labelStrings)
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"update network resource",
				"network "+name+": failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configMap.(map[string]any)
		if ok {
			network.SetConfig(configDataMap)
		} else {
			resp.Diagnostics.AddError(
				"update network resource",
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
		for _, tenantID := range tenantIDs {
			if !tenantID.IsNull() {
				tenant := sdk.GetAlerts200ResponseAllOfChecksInnerAccount{}
				tenant.SetId(tenantID.ValueInt64())
				tenants = append(tenants, tenant)
			}
		}
		network.SetTenants(tenants)
	}

	if !plan.NetworkDomainId.IsNull() && !plan.NetworkDomainId.IsUnknown() {
		networkDomain := sdk.NewListNetworks200ResponseAllOfNetworksInnerNetworkDomain()
		networkDomain.SetId(plan.NetworkDomainId.ValueInt64())
		network.SetNetworkDomain(*networkDomain)
	}

	if !plan.NetworkProxyId.IsNull() && !plan.NetworkProxyId.IsUnknown() {
		networkProxy := sdk.NewListNetworks200ResponseAllOfNetworksInnerNetworkProxy()
		networkProxy.SetId(plan.NetworkProxyId.ValueInt64())
		network.SetNetworkProxy(*networkProxy)
	}

	// Handle resource permissions
	if !plan.ResourcePermissions.IsNull() &&
		!plan.ResourcePermissions.IsUnknown() {
		resourcePermissions := sdk.
			NewUpdateNetworkRequestNetworkResourcePermissions()

		allValue := plan.ResourcePermissions.All.ValueBool()
		resourcePermissions.SetAll(allValue)

		if !plan.ResourcePermissions.GroupIds.IsNull() &&
			!plan.ResourcePermissions.GroupIds.IsUnknown() {
			var groupIDs []types.Int64
			diags := plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &groupIDs, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			var sites []sdk.UpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsSitesInner
			for _, groupID := range groupIDs {
				if !groupID.IsNull() {
					site := sdk.NewUpdateClusterDatastoreRequestDatastorePermissionsResourcePermissionsSitesInner()
					site.SetId(groupID.ValueInt64())
					sites = append(sites, *site)
				}
			}
			if len(sites) > 0 {
				resourcePermissions.SetSites(sites)
			}
		}

		network.SetResourcePermissions(*resourcePermissions)
	}

	updateNetworkReq := sdk.NewUpdateNetworkRequest()
	updateNetworkReq.SetNetwork(*network)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update network resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	_, hresp, err := client.NetworksAPI.UpdateNetwork(ctx, id).
		UpdateNetworkRequest(*updateNetworkReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network resource",
			fmt.Sprintf("network %d UPDATE failed: %s",
				id, errors.ErrMsg(err, hresp)),
		)

		return
	}

	networkState, diags := getNetworkAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update network resource",
			fmt.Sprintf("network %d: failed to read from api", id),
		)

		return
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		networkState.Config = plan.Config
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &networkState)...)
}
