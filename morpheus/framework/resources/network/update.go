// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
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

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	// A clearable attribute is sent when known, sent as its cleared form when
	// pending (an empty string, or null for a nullable field — Morpheus
	// stores both as null), and not sent otherwise.
	network.DisplayName = pending.String("display_name", plan.DisplayName)

	if !plan.Ipv4enabled.IsNull() && !plan.Ipv4enabled.IsUnknown() {
		network.Ipv4Enabled = plan.Ipv4enabled.ValueBoolPointer()
	}

	if !plan.Ipv6enabled.IsNull() && !plan.Ipv6enabled.IsUnknown() {
		network.Ipv6Enabled = plan.Ipv6enabled.ValueBoolPointer()
	}

	network.CidrIPv6 = pending.String("cidr_ipv6", plan.CidrIpv6)
	clearable.NullableString(pending, "gateway_ipv6", plan.GatewayIpv6, &network.GatewayIPv6)
	clearable.NullableString(pending, "netmask_ipv6", plan.NetmaskIpv6, &network.NetmaskIPv6)
	clearable.NullableString(pending, "dns_primary_ipv6", plan.DnsPrimaryIpv6, &network.DnsPrimaryIPv6)
	clearable.NullableString(pending, "dns_secondary_ipv6", plan.DnsSecondaryIpv6, &network.DnsSecondaryIPv6)

	if !plan.DhcpServerIpv6.IsNull() && !plan.DhcpServerIpv6.IsUnknown() {
		network.DhcpServerIPv6 = plan.DhcpServerIpv6.ValueBoolPointer()
	}

	clearable.NullableInt64(pending, "pool_ipv6_id", plan.PoolIpv6Id, &network.PoolIPv6)
	clearable.NullableString(pending, "description", plan.Description, &network.Description)
	network.Cidr = pending.String("cidr", plan.Cidr)
	network.Gateway = pending.String("gateway", plan.Gateway)
	network.DnsPrimary = pending.String("dns_primary", plan.DnsPrimary)
	network.DnsSecondary = pending.String("dns_secondary", plan.DnsSecondary)

	clearable.NullableInt64(pending, "vlan_id", plan.VlanId, &network.VlanId)
	clearable.NullableInt64(pending, "pool_id", plan.PoolId, &network.Pool)

	// Associations are sent as {"id": n}; {"id": null} removes one.
	if id, ok := idRef(pending, "zone_pool_id", plan.ZonePoolId); ok {
		network.ZonePool = &sdk.UpdateNetworkRequestNetworkZonePool{Id: id}
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

	network.SearchDomains = pending.String("search_domains", plan.SearchDomains)
	network.SwitchId = pending.String("switch_id", plan.SwitchId)

	if !plan.ApplianceUrlProxyBypass.IsNull() && !plan.ApplianceUrlProxyBypass.IsUnknown() {
		network.ApplianceUrlProxyBypass = plan.ApplianceUrlProxyBypass.ValueBoolPointer()
	}

	clearable.NullableString(pending, "no_proxy", plan.NoProxy, &network.NoProxy)

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		network.Visibility = plan.Visibility.ValueStringPointer()
	}

	labels, sendLabels, ldiags := pending.Strings(ctx, "labels", plan.Labels)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if sendLabels {
		network.Labels = labels
	}

	// Type-specific configuration, resolved from the write-only inputs. Morpheus
	// merges it into the network's existing config map: keys not sent are kept,
	// so a removed connected_gateway is cleared by sending it empty. The map is
	// copied: inputs.all is what the stored digests describe.
	if inputs.all != nil || pending.Clears("connected_gateway") {
		network.Config = make(map[string]any, len(inputs.all)+1)
		for k, v := range inputs.all {
			network.Config[k] = v
		}

		if pending.Clears("connected_gateway") {
			network.Config["connectedGateway"] = ""
		}
	}

	tenantIDs, sendTenants, tdiags := pending.Int64s(ctx, "tenant_ids", plan.TenantIds)
	resp.Diagnostics.Append(tdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if sendTenants {
		tenants := make([]sdk.UpdateNetworkRequestNetworkTenantsInner, 0, len(tenantIDs))
		for _, tenantID := range tenantIDs {
			tenants = append(tenants, sdk.UpdateNetworkRequestNetworkTenantsInner{Id: &tenantID})
		}

		network.Tenants = tenants
	}

	if id, ok := idRef(pending, "network_domain_id", plan.NetworkDomainId); ok {
		network.NetworkDomain = &sdk.UpdateNetworkRequestNetworkNetworkDomain{Id: id}
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

	if id, ok := idRef(pending, "network_proxy_id", plan.NetworkProxyId); ok {
		network.NetworkProxy = &sdk.UpdateNetworkRequestNetworkNetworkProxy{Id: id}
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

	networkState, found, diags := getNetworkAsState(ctx, id, client, plan)
	if !found && !diags.HasError() {
		diags.AddError("update network resource", fmt.Sprintf("network %d not found after update", id))
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update network resource",
			fmt.Sprintf("network %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(networkconfig.Store(ctx, resp.Private, inputs.digests)...)
	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &networkState)...)
}

// idRef is the id of an association object in the update request: null when
// the attribute is being cleared, the planned id when known, and false — send
// no object — otherwise.
func idRef(pending clearable.Pending, name string, plan types.Int64) (sdk.NullableInt64, bool) {
	var id sdk.NullableInt64

	ok := clearable.NullableInt64(pending, name, plan, &id)

	return id, ok
}
