// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
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

	// Set Config - convert dynamic to SDK config structure
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create policy resource",
				"policy "+name+": failed to convert config: "+err.Error(),
			)

			return
		}

		// Check if config is empty
		if configMapTyped, ok := configMap.(map[string]interface{}); ok && len(configMapTyped) == 0 {
			resp.Diagnostics.AddError(
				"create policy resource",
				fmt.Sprintf("policy %s: config cannot be empty for policy type '%s'. "+
					"Please provide the required configuration fields for this policy type.", name, policyTypeCode),
			)

			return
		}

		// Marshal to JSON then unmarshal to SDK config structure
		// This allows the SDK's UnmarshalJSON to handle the oneOf structure
		configJSON, err := json.Marshal(configMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"create policy resource",
				"policy "+name+": failed to marshal config to JSON: "+err.Error(),
			)

			return
		}

		var sdkConfig sdk.AddPoliciesRequestPolicyConfig
		if err := json.Unmarshal(configJSON, &sdkConfig); err != nil {
			resp.Diagnostics.AddError(
				"create policy resource",
				fmt.Sprintf("policy %s: invalid config for policy type '%s': %s", name, policyTypeCode, err.Error()),
			)

			return
		}

		addPolicy.SetConfig(sdkConfig)
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

	// Write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the created policy to get full state
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
