// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkPoolIpModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	networkPoolID := plan.NetworkPoolId.ValueInt64()

	ip := &sdk.CreateNetworkPoolIpRequestNetworkPoolIp{
		Hostname: plan.Hostname.ValueString(),
	}

	// Omit ipAddress when unset or unknown so Morpheus allocates the next free
	// address from the pool and returns it in the response. This is the main
	// feature of the resource.
	if !plan.IpAddress.IsNull() && !plan.IpAddress.IsUnknown() {
		ip.IpAddress = plan.IpAddress.ValueStringPointer()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ip.Description = plan.Description.ValueStringPointer()
	}

	if !plan.RefType.IsNull() && !plan.RefType.IsUnknown() {
		ip.RefType = plan.RefType.ValueStringPointer()
	}

	if !plan.RefId.IsNull() && !plan.RefId.IsUnknown() {
		ip.RefId = plan.RefId.ValueInt64Pointer()
	}

	// Domain is set only when domain_id is given; AlsoRequires(domain_id)
	// guards create_dns, so create_dns being true implies a domain.
	if !plan.DomainId.IsNull() && !plan.DomainId.IsUnknown() {
		ip.Domain = &sdk.CreateNetworkPoolIpRequestNetworkPoolIpDomain{
			Id: plan.DomainId.ValueInt64Pointer(),
		}
	}

	createReq := sdk.CreateNetworkPoolIpRequest{
		NetworkPoolIp: ip,
	}

	// Send createDns only when true; the default is false and the server need
	// not be told to skip DNS.
	if plan.CreateDns.ValueBool() {
		createReq.CreateDns = plan.CreateDns.ValueBoolPointer()
	}

	createResp, httpResp, err := client.NetworksAPI.
		CreateNetworkPoolIp(ctx, networkPoolID).
		CreateNetworkPoolIpRequest(createReq).
		Execute()
	// A 400 carries Morpheus's errors map (for example ipAddress "not within
	// the pool's IP ranges"); errfmt.ErrMsg includes the response body
	// verbatim so the diagnostic surfaces that message.
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating network pool IP",
			fmt.Sprintf("network pool %d IP POST failed: ", networkPoolID)+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	created := createResp.NetworkPoolIp
	if created == nil || created.Id == nil {
		resp.Diagnostics.AddError(
			"error creating network pool IP",
			fmt.Sprintf("network pool %d IP POST returned no id", networkPoolID),
		)

		return
	}

	id := *created.Id

	// Read back the full reservation so every computed attribute is populated
	// from the authoritative GET; the allocated ip_address is on the create
	// response too, but the read keeps all fields consistent.
	state, found, diags := getNetworkPoolIpAsState(
		ctx, networkPoolID, id, client, plan,
	)
	if !found && !diags.HasError() {
		diags.AddError("read network pool IP", "reservation not found immediately after create")
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_pool_ip",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
