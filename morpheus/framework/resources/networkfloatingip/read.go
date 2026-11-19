// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *networkFloatingIpResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkFloatingIpModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkFloatingIp(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_floating_ip", "", err, httpResp)

		return
	}

	if result.NetworkFloatingIp == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkFloatingIp is nil in the response")

		return
	}

	// network_server_id and floating_ip_pool_id are not returned by the API, so
	// carry them forward from the prior state.
	newState := floatingIpAsState(result.NetworkFloatingIp, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// floatingIpAsState maps the API response onto the resource model. prior
// supplies the two inputs the API never returns (network_server_id and
// floating_ip_pool_id).
func floatingIpAsState(
	fip *sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIp,
	prior *NetworkFloatingIpModel,
) NetworkFloatingIpModel {
	state := NetworkFloatingIpModel{
		Id:               convert.Int64ToType(fip.Id),
		NetworkServerId:  prior.NetworkServerId,
		FloatingIpPoolId: prior.FloatingIpPoolId,
		IpAddress:        convert.StrToType(fip.IpAddress),
		IpStatus:         convert.StrToType(fip.IpStatus),
		ExternalId:       convert.StrToType(fip.ExternalId.Get()),
		PtrId:            convert.StrToType(fip.PtrId.Get()),
		DateCreated:      convert.StrToType(timePtr(fip.DateCreated)),
		LastUpdated:      convert.StrToType(timePtr(fip.LastUpdated)),
	}

	if fip.Cloud != nil {
		state.CloudId = convert.Int64ToType(fip.Cloud.Id)
		state.CloudName = convert.StrToType(fip.Cloud.Name)
	} else {
		state.CloudId = types.Int64Null()
		state.CloudName = types.StringNull()
	}

	if fip.Server != nil {
		state.ServerId = convert.Int64ToType(fip.Server.Id)
		state.ServerName = convert.StrToType(fip.Server.Name)
	} else {
		state.ServerId = types.Int64Null()
		state.ServerName = types.StringNull()
	}

	if fip.IpRange != nil {
		state.IpRangeStartAddress = convert.StrToType(fip.IpRange.StartAddress)
		state.IpRangeEndAddress = convert.StrToType(fip.IpRange.EndAddress)
	} else {
		state.IpRangeStartAddress = types.StringNull()
		state.IpRangeEndAddress = types.StringNull()
	}

	if fip.NetworkDomain != nil {
		state.NetworkDomainId = convert.Int64ToType(fip.NetworkDomain.Id)
		state.NetworkDomainName = convert.StrToType(fip.NetworkDomain.Name)
	} else {
		state.NetworkDomainId = types.Int64Null()
		state.NetworkDomainName = types.StringNull()
	}

	return state
}

// timePtr formats a time pointer as an RFC 3339 string pointer, or nil when the
// API omitted it.
func timePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(time.RFC3339)

	return &s
}
