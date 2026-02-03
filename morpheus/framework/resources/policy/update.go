// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()
	name := plan.Name.ValueString()

	updatePolicy := sdk.NewUpdatePoliciesRequestPolicyWithDefaults()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)

		return
	}

	// Set required fields
	updatePolicy.SetName(name)

	// Note: PolicyType, AssociatedResourceType, and AssociatedResourceId are not included
	// in updates. PolicyType is not updatable via the API (not present in UpdatePoliciesRequestPolicy).
	// AssociatedResourceType and AssociatedResourceId require replacement
	// While updatePolicy.SetRefType and updatePolicy.SetRefId exist, they are ineffectual and
	// neither refType nor refId can be updated via the API.

	// Set optional fields
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updatePolicy.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		updatePolicy.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.EachUser.IsNull() && !plan.EachUser.IsUnknown() {
		updatePolicy.SetEachUser(plan.EachUser.ValueBool())
	}

	// Set tenant IDs if provided
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantIDs []int64
		diags := plan.Tenants.ElementsAs(ctx, &tenantIDs, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}
		updatePolicy.SetAccounts(tenantIDs)
	}

	// Set Config - convert state config fields to SDK config structure
	sdkConfig, configDiags := mapStateToUpdatePolicyConfig(ctx, &plan)
	if configDiags.HasError() {
		resp.Diagnostics.Append(configDiags...)

		return
	}
	if sdkConfig != nil {
		updatePolicy.SetConfig(*sdkConfig)
	}

	updatePolicyRequest := sdk.NewUpdatePoliciesRequest(*updatePolicy)

	policy, hresp, err := client.PoliciesAPI.UpdatePolicies(ctx, id).
		UpdatePoliciesRequest(*updatePolicyRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d PUT failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if policy.Policy == nil || policy.Policy.Id == nil {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: id is nil", id),
		)

		return
	}

	newID := *policy.Policy.Id
	if newID != id {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: id mismatch %d != %d", id, id, newID),
		)

		return
	}

	// Read the updated policy to get full state
	updatedState, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}
