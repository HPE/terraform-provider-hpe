// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState LoadBalancerProfileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	profile := &sdk.UpdateLoadBalancerProfileRequestLoadBalancerProfile{}
	profile.Name = sdk.PtrString(plan.Name.ValueString())

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		profile.Description = sdk.PtrString(plan.Description.ValueString())
	}

	profile.ServiceType = sdk.PtrString(plan.ServiceType.ValueString())

	cfg, diags := buildUpdateConfig(ctx, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	profile.Config = cfg

	loadBalancerID := currentState.LoadBalancerId.ValueInt64()
	id := currentState.Id.ValueInt64()

	updateReq := sdk.UpdateLoadBalancerProfileRequest{
		LoadBalancerProfile: profile,
	}

	_, httpResp, err := client.LoadBalancersAPI.
		UpdateLoadBalancerProfile(ctx, loadBalancerID, id).
		UpdateLoadBalancerProfileRequest(updateReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating load balancer profile",
			"load balancer profile "+plan.Name.ValueString()+" PUT failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, stateDiags := getLoadBalancerProfileAsState(
		ctx, loadBalancerID, id, client, plan,
	)
	if resp.Diagnostics.Append(stateDiags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
