// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

	network := &sdk.UpdateNetworkRequestNetwork{}

	// Set all updateable fields from plan
	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
		displayName := plan.DisplayName.ValueString()
		network.DisplayName = &displayName
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		network.Description.Set(&description)
	}

	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		cidr := plan.Cidr.ValueString()
		network.Cidr = &cidr
	}

	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		gateway := plan.Gateway.ValueString()
		network.Gateway = &gateway
	}

	if !plan.DnsPrimary.IsNull() && !plan.DnsPrimary.IsUnknown() {
		dnsPrimary := plan.DnsPrimary.ValueString()
		network.DnsPrimary = &dnsPrimary
	}

	if !plan.DnsSecondary.IsNull() && !plan.DnsSecondary.IsUnknown() {
		dnsSecondary := plan.DnsSecondary.ValueString()
		network.DnsSecondary = &dnsSecondary
	}

	if !plan.VlanId.IsNull() && !plan.VlanId.IsUnknown() {
		vlanID := plan.VlanId.ValueInt64()
		network.VlanId = &vlanID
	}

	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		poolID := plan.PoolId.ValueInt64()
		network.Pool.Set(&poolID)
	}

	if !plan.ZonePoolId.IsNull() && !plan.ZonePoolId.IsUnknown() {
		zonePoolID := plan.ZonePoolId.ValueInt64()
		network.ZonePool = &sdk.UpdateNetworkRequestNetworkZonePool{Id: &zonePoolID}
	}

	if !plan.AllowStaticOverride.IsNull() && !plan.AllowStaticOverride.IsUnknown() {
		allowStaticOverride := plan.AllowStaticOverride.ValueBool()
		network.AllowStaticOverride = &allowStaticOverride
	}

	if !plan.AssignPublicIp.IsNull() && !plan.AssignPublicIp.IsUnknown() {
		assignPublicIP := plan.AssignPublicIp.ValueBool()
		network.AssignPublicIp = &assignPublicIP
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		active := plan.Active.ValueBool()
		network.Active = &active
	}

	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		dhcpServer := plan.DhcpServer.ValueBool()
		network.DhcpServer = &dhcpServer
	}

	if !plan.SearchDomains.IsNull() && !plan.SearchDomains.IsUnknown() {
		searchDomains := plan.SearchDomains.ValueString()
		network.SearchDomains = &searchDomains
	}

	if !plan.SwitchId.IsNull() && !plan.SwitchId.IsUnknown() {
		switchID := plan.SwitchId.ValueString()
		network.SwitchId = &switchID
	}

	if !plan.ApplianceUrlProxyBypass.IsNull() && !plan.ApplianceUrlProxyBypass.IsUnknown() {
		applianceURLProxyBypass := plan.ApplianceUrlProxyBypass.ValueBool()
		network.ApplianceUrlProxyBypass = &applianceURLProxyBypass
	}

	if !plan.NoProxy.IsNull() && !plan.NoProxy.IsUnknown() {
		noProxy := plan.NoProxy.ValueString()
		network.NoProxy.Set(&noProxy)
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		visibility := plan.Visibility.ValueString()
		network.Visibility = &visibility
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
		network.Labels = labelStrings
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
			network.Config = configDataMap
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

		var tenants []sdk.UpdateNetworkRequestNetworkTenantsInner
		for _, tenantID := range tenantIDs {
			if !tenantID.IsNull() {
				id := tenantID.ValueInt64()
				tenant := sdk.UpdateNetworkRequestNetworkTenantsInner{Id: &id}
				tenants = append(tenants, tenant)
			}
		}
		network.Tenants = tenants
	}

	if !plan.NetworkDomainId.IsNull() && !plan.NetworkDomainId.IsUnknown() {
		networkDomainID := plan.NetworkDomainId.ValueInt64()
		network.NetworkDomain = &sdk.UpdateNetworkRequestNetworkNetworkDomain{Id: &networkDomainID}
	}

	if !plan.NetworkProxyId.IsNull() && !plan.NetworkProxyId.IsUnknown() {
		networkProxyID := plan.NetworkProxyId.ValueInt64()
		network.NetworkProxy = &sdk.UpdateNetworkRequestNetworkNetworkProxy{Id: &networkProxyID}
	}

	updateNetworkReq := &sdk.UpdateNetworkRequest{Network: network}

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
				id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	networkState, diags := getNetworkAsState(ctx, id, client, plan)
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
