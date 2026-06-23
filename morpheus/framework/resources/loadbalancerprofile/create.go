// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan LoadBalancerProfileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	profile := &sdk.CreateLoadBalancerProfileRequestLoadBalancerProfile{}
	profile.Name = sdk.PtrString(plan.Name.ValueString())

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		profile.Description = sdk.PtrString(plan.Description.ValueString())
	}

	profile.ServiceType = sdk.PtrString(plan.ServiceType.ValueString())

	cfg, diags := buildCreateConfig(ctx, plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	profile.Config = cfg

	loadBalancerID := plan.LoadBalancerId.ValueInt64()

	createReq := sdk.CreateLoadBalancerProfileRequest{
		LoadBalancerProfile: profile,
	}

	createResp, httpResp, err := client.LoadBalancersAPI.
		CreateLoadBalancerProfile(ctx, loadBalancerID).
		CreateLoadBalancerProfileRequest(createReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating load balancer profile",
			"load balancer profile "+plan.Name.ValueString()+" POST failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	created := createResp.LoadBalancerProfile
	if created == nil || created.Id == nil {
		resp.Diagnostics.AddError(
			"error creating load balancer profile",
			"load balancer profile "+plan.Name.ValueString()+": id is nil",
		)

		return
	}

	state, stateDiags := getLoadBalancerProfileAsState(ctx, loadBalancerID, *created.Id, client, plan)
	if stateDiags.HasError() {
		resp.Diagnostics.Append(stateDiags...)

		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "load_balancer_profile",
			ResourceID:   *created.Id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
