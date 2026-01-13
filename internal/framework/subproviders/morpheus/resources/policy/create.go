// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HPE/terraform-provider-hpe/internal/framework/utils"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	addPolicy := sdk.NewAddPoliciesRequestPolicyWithDefaults()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	// Set required fields
	addPolicy.SetName(name)

	policyTypeCode := plan.PolicyType.Code.ValueString()
	policyType := sdk.NewAddPoliciesRequestPolicyPolicyTypeWithDefaults()
	policyType.SetCode(policyTypeCode)
	addPolicy.SetPolicyType(*policyType)

	// When associated_resource_type is "Global", leave refType null (null defaults to Global)
	associatedResourceType := plan.AssociatedResourceType.ValueString()
	if associatedResourceType != AssociatedResourceTypeGlobal {
		// Convert user-facing resource type to API type
		apiType := resourceTypeToAPIType(associatedResourceType)
		addPolicy.SetRefType(apiType)
	}

	// Set AssociatedResourceId if provided - required when RefType is not Global
	// This is validated in ValidateConfig
	if !plan.AssociatedResourceId.IsNull() && !plan.AssociatedResourceId.IsUnknown() {
		addPolicy.SetRefId(plan.AssociatedResourceId.ValueInt64())
	}

	// Set optional fields
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		addPolicy.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		addPolicy.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.EachUser.IsNull() && !plan.EachUser.IsUnknown() {
		addPolicy.SetEachUser(plan.EachUser.ValueBool())
	}

	// Set tenant IDs if provided
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantIDs []int64
		diags := plan.Tenants.ElementsAs(ctx, &tenantIDs, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)

			return
		}
		addPolicy.SetAccounts(tenantIDs)
	}

	// Set Config - convert state config fields to SDK config structure
	sdkConfig, configDiags := mapStateToAddPolicyConfig(ctx, &plan)
	if configDiags.HasError() {
		resp.Diagnostics.Append(configDiags...)

		return
	}
	if sdkConfig != nil {
		addPolicy.SetConfig(*sdkConfig)
	}

	addPolicyRequest := sdk.NewAddPoliciesRequest(*addPolicy)

	policy, hresp, err := client.PoliciesAPI.AddPolicies(ctx).AddPoliciesRequest(*addPolicyRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if policy.Policy == nil || policy.Policy.Id == nil {
		resp.Diagnostics.AddError(
			"create policy resource",
			"policy "+name+" id is nil",
		)

		return
	}

	id := *policy.Policy.Id
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		utils.TaintResourceState(ctx, utils.TaintResourceStateConfig{
			ResourceType: "policy",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Read the created policy to get full state
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"failed to read policy state",
			fmt.Sprintf("Policy %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set policy state",
			fmt.Sprintf("Policy %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}
