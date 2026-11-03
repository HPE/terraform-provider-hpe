// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
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

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePolicy := &sdk.UpdatePoliciesRequestPolicy{}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)

		return
	}

	// Set required fields
	updatePolicy.Name = sdk.PtrString(name)

	// Note: PolicyType, AssociatedResourceType, and AssociatedResourceId are not included
	// in updates. PolicyType is not updatable via the API (not present in UpdatePoliciesRequestPolicy).
	// AssociatedResourceType and AssociatedResourceId require replacement
	// While updatePolicy.SetRefType and updatePolicy.SetRefId exist, they are ineffectual and
	// neither refType nor refId can be updated via the API.

	// Set optional fields
	clearable.NullableString(pending, "description", plan.Description, &updatePolicy.Description)

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		updatePolicy.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.EachUser.IsNull() && !plan.EachUser.IsUnknown() {
		updatePolicy.EachUser = plan.EachUser.ValueBoolPointer()
	}

	// tenants: the planned list, or an empty list when being cleared
	if tenantIDs, ok, ldiags := pending.Int64s(ctx, "tenants", plan.Tenants); ok {
		resp.Diagnostics.Append(ldiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updatePolicy.Accounts = tenantIDs
	}

	// Set Config - convert state config fields to SDK config structure
	sdkConfig, configDiags := mapStateToUpdatePolicyConfig(ctx, &plan)
	if configDiags.HasError() {
		resp.Diagnostics.Append(configDiags...)

		return
	}
	if sdkConfig != nil {
		updatePolicy.Config = sdkConfig
	}

	updatePolicyRequest := &sdk.UpdatePoliciesRequest{Policy: *updatePolicy}

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
	updatedState, found, diags := getPolicyAsState(ctx, id, client, &plan)
	if !found && !diags.HasError() {
		diags.AddError("update policy resource", fmt.Sprintf("policy %d not found after update", id))
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}
