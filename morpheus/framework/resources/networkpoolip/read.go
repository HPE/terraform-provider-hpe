// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data NetworkPoolIpModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, found, diags := getNetworkPoolIpAsState(
		ctx, data.NetworkPoolId.ValueInt64(), data.Id.ValueInt64(), client, data,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// The reservation was deleted out of band (Morpheus returned 404); remove
	// it from state so the next plan recreates it instead of erroring.
	if !found {
		tflog.Info(ctx, fmt.Sprintf("network pool IP %d not found, removing from state", data.Id.ValueInt64()))
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made: create_dns is
	// never returned by the API, so after an import it would otherwise stay
	// null rather than its default of false. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, NetworkPoolIpResourceSchema(ctx), &resp.State)...,
	)
}

// getNetworkPoolIpAsState reads the reservation and maps it to model state.
// prior carries the attributes Morpheus never returns (domain_id, create_dns)
// and the parent pool id when the response omits it, so a read does not drop
// what only the configuration knows.
func getNetworkPoolIpAsState(
	ctx context.Context,
	networkPoolID int64,
	id int64,
	client *sdk.APIClient,
	prior NetworkPoolIpModel,
) (NetworkPoolIpModel, bool, diag.Diagnostics) {
	var state NetworkPoolIpModel
	var diags diag.Diagnostics

	ipResp, httpResp, err := client.NetworksAPI.
		GetNetworkPoolIp(ctx, networkPoolID, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return state, false, diags
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading network pool IP",
			fmt.Sprintf("network pool IP %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	p := ipResp.NetworkPoolIp
	if p == nil {
		diags.AddError(
			"error reading network pool IP",
			fmt.Sprintf("network pool IP %d GET returned no reservation", id),
		)

		return state, false, diags
	}

	state.Id = convert.Int64ToType(p.Id)

	// network_pool_id: use the response value when present, else carry the
	// parent id from prior plan/state (it is RequiresReplace, so it is stable).
	if p.NetworkPoolId != nil {
		state.NetworkPoolId = types.Int64Value(*p.NetworkPoolId)
	} else {
		state.NetworkPoolId = prior.NetworkPoolId
	}

	// Nullable response fields: the convert helpers handle nil.
	state.IpType = convert.StrToType(p.IpType)
	state.IpAddress = convert.StrToType(p.IpAddress)
	state.GatewayAddress = convert.StrToType(p.GatewayAddress.Get())
	state.SubnetMask = convert.StrToType(p.SubnetMask.Get())
	state.DnsServer = convert.StrToType(p.DnsServer.Get())
	state.InterfaceName = convert.StrToType(p.InterfaceName.Get())
	state.Description = convert.StrToType(p.Description.Get())
	state.Active = convert.BoolToType(p.Active)
	state.StaticIp = convert.BoolToType(p.StaticIp)
	state.Fqdn = convert.StrToType(p.Fqdn)
	state.DomainName = convert.StrToType(p.DomainName.Get())
	state.Hostname = convert.StrToType(p.Hostname)
	state.ExternalId = convert.StrToType(p.ExternalId.Get())
	state.RefType = convert.StrToType(p.RefType.Get())
	state.RefId = convert.Int64ToType(p.RefId.Get())

	// Date fields: *time.Time formatted RFC3339 by convert; EndDate is
	// NullableTime so its *time.Time comes from .Get().
	state.DateCreated = convert.TimeToType(p.DateCreated)
	state.LastUpdated = convert.TimeToType(p.LastUpdated)
	state.StartDate = convert.TimeToType(p.StartDate)
	state.EndDate = convert.TimeToType(p.EndDate.Get())

	// domain_id and create_dns are create-only inputs Morpheus never returns;
	// carry them from prior plan/state. On import they are null/false.
	state.DomainId = prior.DomainId
	state.CreateDns = prior.CreateDns

	return state, true, diags
}
