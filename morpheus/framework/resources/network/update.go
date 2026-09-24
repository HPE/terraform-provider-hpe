// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config NetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// The write-only configuration is only present in the configuration.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Routed for the network's type — the type cannot change in place, and
	// this is what ModifyPlan routed for when it compared the digests.
	inputs, idiags := resolveConfigInputs(ctx, config, r.planTypeSpec(ctx, config, state))
	resp.Diagnostics.Append(idiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	network := &sdk.UpdateNetworkRequestNetwork{}

	// Set all updateable fields from plan. Whether Morpheus applies a change
	// to name, cidr, gateway, vlanId, dhcpServer or dns depends on the network
	// type's flags; ModifyPlan has already planned a replacement where it
	// would not, so what reaches here is applied.
	// name only when it changes: the API sets a null displayName to the name it
	// is sent, so sending an unchanged name would alter display_name behind
	// the plan's back.
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() && plan.Name.ValueString() != state.Name.ValueString() {
		network.Name = plan.Name.ValueStringPointer()
	}

	if !plan.DisplayName.IsNull() && !plan.DisplayName.IsUnknown() {
		network.DisplayName = plan.DisplayName.ValueStringPointer()
	}

	if !plan.Ipv4enabled.IsNull() && !plan.Ipv4enabled.IsUnknown() {
		network.Ipv4Enabled = plan.Ipv4enabled.ValueBoolPointer()
	}

	if !plan.Ipv6enabled.IsNull() && !plan.Ipv6enabled.IsUnknown() {
		network.Ipv6Enabled = plan.Ipv6enabled.ValueBoolPointer()
	}

	if !plan.CidrIpv6.IsNull() && !plan.CidrIpv6.IsUnknown() {
		network.CidrIPv6 = plan.CidrIpv6.ValueStringPointer()
	}

	if !plan.GatewayIpv6.IsNull() && !plan.GatewayIpv6.IsUnknown() {
		gatewayIpv6 := plan.GatewayIpv6.ValueString()
		network.GatewayIPv6.Set(&gatewayIpv6)
	}

	if !plan.NetmaskIpv6.IsNull() && !plan.NetmaskIpv6.IsUnknown() {
		netmaskIpv6 := plan.NetmaskIpv6.ValueString()
		network.NetmaskIPv6.Set(&netmaskIpv6)
	}

	if !plan.DnsPrimaryIpv6.IsNull() && !plan.DnsPrimaryIpv6.IsUnknown() {
		dnsPrimaryIpv6 := plan.DnsPrimaryIpv6.ValueString()
		network.DnsPrimaryIPv6.Set(&dnsPrimaryIpv6)
	}

	if !plan.DnsSecondaryIpv6.IsNull() && !plan.DnsSecondaryIpv6.IsUnknown() {
		dnsSecondaryIpv6 := plan.DnsSecondaryIpv6.ValueString()
		network.DnsSecondaryIPv6.Set(&dnsSecondaryIpv6)
	}

	if !plan.DhcpServerIpv6.IsNull() && !plan.DhcpServerIpv6.IsUnknown() {
		network.DhcpServerIPv6 = plan.DhcpServerIpv6.ValueBoolPointer()
	}

	if !plan.PoolIpv6Id.IsNull() && !plan.PoolIpv6Id.IsUnknown() {
		poolIpv6ID := plan.PoolIpv6Id.ValueInt64()
		network.PoolIPv6.Set(&poolIpv6ID)
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		network.Description.Set(&description)
	}

	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		network.Cidr = plan.Cidr.ValueStringPointer()
	}

	if !plan.Gateway.IsNull() && !plan.Gateway.IsUnknown() {
		network.Gateway = plan.Gateway.ValueStringPointer()
	}

	if !plan.DnsPrimary.IsNull() && !plan.DnsPrimary.IsUnknown() {
		network.DnsPrimary = plan.DnsPrimary.ValueStringPointer()
	}

	if !plan.DnsSecondary.IsNull() && !plan.DnsSecondary.IsUnknown() {
		network.DnsSecondary = plan.DnsSecondary.ValueStringPointer()
	}

	if !plan.VlanId.IsNull() && !plan.VlanId.IsUnknown() {
		network.VlanId = plan.VlanId.ValueInt64Pointer()
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
		network.AllowStaticOverride = plan.AllowStaticOverride.ValueBoolPointer()
	}

	if !plan.AssignPublicIp.IsNull() && !plan.AssignPublicIp.IsUnknown() {
		network.AssignPublicIp = plan.AssignPublicIp.ValueBoolPointer()
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		network.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		network.DhcpServer = plan.DhcpServer.ValueBoolPointer()
	}

	if !plan.SearchDomains.IsNull() && !plan.SearchDomains.IsUnknown() {
		network.SearchDomains = plan.SearchDomains.ValueStringPointer()
	}

	if !plan.SwitchId.IsNull() && !plan.SwitchId.IsUnknown() {
		network.SwitchId = plan.SwitchId.ValueStringPointer()
	}

	if !plan.ApplianceUrlProxyBypass.IsNull() && !plan.ApplianceUrlProxyBypass.IsUnknown() {
		network.ApplianceUrlProxyBypass = plan.ApplianceUrlProxyBypass.ValueBoolPointer()
	}

	if !plan.NoProxy.IsNull() && !plan.NoProxy.IsUnknown() {
		noProxy := plan.NoProxy.ValueString()
		network.NoProxy.Set(&noProxy)
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		network.Visibility = plan.Visibility.ValueStringPointer()
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

	// Type-specific configuration, resolved from the write-only inputs. Morpheus
	// merges it into the network's existing config map: keys not sent are kept.
	if inputs.all != nil {
		network.Config = inputs.all
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

	// The network server and a child network's parent; the type-specific
	// fields Morpheus does not edit after create (scopeId, vxlanId) are not
	// sent, since a change to them is planned as a replacement.
	if v := knownInt64(plan.NetworkServerId); v != 0 {
		network.NetworkServer = &sdk.UpdateNetworkRequestNetworkNetworkServer{Id: &v}
	}

	if v := knownInt64(plan.ParentNetworkId); v != 0 {
		network.ParentNetwork = &sdk.UpdateNetworkRequestNetworkParentNetwork{Id: &v}
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

	resp.Diagnostics.Append(networkconfig.Store(ctx, resp.Private, inputs.digests)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &networkState)...)
}
