// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// policy provides the package for hpe_morpheus_policy
package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_policy"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = PolicyResourceSchema(ctx)
}

// resourceTypeToAPIType converts user-facing resource types to API types
func resourceTypeToAPIType(resourceType string) string {
	switch resourceType {
	case "Cloud":
		return "ComputeZone"
	case "Group":
		return "ComputeSite"
	default:
		// For other types (Global, User, Role, Network, Plan), pass through as-is
		return resourceType
	}
}

// apiTypeToResourceType converts API types back to user-facing resource types
func apiTypeToResourceType(apiType string) string {
	switch apiType {
	case "ComputeZone":
		return "Cloud"
	case "ComputeSite":
		return "Group"
	default:
		// For other types (Global, User, Role, Network, Plan), pass through as-is
		return apiType
	}
}

// populate policy resource model with current API values
func getPolicyAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan *PolicyModel,
) (PolicyModel, diag.Diagnostics) {
	var state PolicyModel
	var diags diag.Diagnostics

	policy, hresp, err := client.PoliciesAPI.GetPolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK || policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d GET failed", id)+errors.ErrMsg(err, hresp),
		)
		return state, diags
	}

	if policy.Policy == nil {
		diags.AddError(
			"populate policy resource",
			fmt.Sprintf("policy %d is nil", id),
		)
		return state, diags
	}

	p := policy.Policy

	// Set basic fields
	state.Id = convert.Int64ToType(p.Id)
	state.Name = convert.StrToType(p.Name)

	// Handle nullable fields properly
	if p.Description.IsSet() {
		state.Description = convert.StrToType(p.Description.Get())
	}

	state.Enabled = convert.BoolToType(p.Enabled)

	if p.EachUser.IsSet() {
		state.EachUser = convert.BoolToType(p.EachUser.Get())
	}

	// Handle RefId - convert string to int64
	if p.RefId.IsSet() && p.RefId.Get() != nil {
		state.AssociatedResourceId = types.Int64Value(p.GetRefId())
	}

	// Handle RefType
	// If RefType is null or not set, it means it's a Global policy
	if p.RefType.IsSet() && p.RefType.Get() != nil {
		apiType := *p.RefType.Get()
		// Convert API type to user-facing resource type
		state.AssociatedResourceType = types.StringValue(apiTypeToResourceType(apiType))
	} else {
		state.AssociatedResourceType = types.StringValue("Global")
	}

	// Set Tenant IDs
	if len(p.Accounts) > 0 {
		tenantIDs := make([]int64, 0, len(p.Accounts))
		for _, acc := range p.Accounts {
			if acc.Id != nil {
				tenantIDs = append(tenantIDs, *acc.Id)
			}
		}
		tenantsSet, setDiags := types.SetValueFrom(ctx, types.Int64Type, tenantIDs)
		if setDiags.HasError() {
			diags.Append(setDiags...)
			return state, diags
		}
		state.Tenants = tenantsSet
	} else {
		// Set to null if no accounts
		state.Tenants = types.SetNull(types.Int64Type)
	}

	// Set PolicyType
	if p.PolicyType != nil {
		policyTypeAttrs := map[string]attr.Value{}
		if p.PolicyType.Code != nil {
			policyTypeAttrs["code"] = types.StringValue(*p.PolicyType.Code)
		} else {
			policyTypeAttrs["code"] = types.StringNull()
		}
		if p.PolicyType.Id != nil {
			policyTypeAttrs["id"] = types.Int64Value(*p.PolicyType.Id)
		} else {
			policyTypeAttrs["id"] = types.Int64Null()
		}
		if p.PolicyType.Name != nil {
			policyTypeAttrs["name"] = types.StringValue(*p.PolicyType.Name)
		} else {
			policyTypeAttrs["name"] = types.StringNull()
		}

		policyTypeValue, policyTypeDiags := NewPolicyTypeValue(
			PolicyTypeValue{}.AttributeTypes(ctx), policyTypeAttrs)
		if policyTypeDiags.HasError() {
			diags.Append(policyTypeDiags...)
			return state, diags
		}
		state.PolicyType = policyTypeValue
	}

	// Handle Config - preserve from plan or convert from API
	state.Config = types.DynamicNull()

	if plan != nil && !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		// Preserve config from plan - the API doesn't always return the full config structure
		state.Config = plan.Config
	} else if p.Config != nil {
		// Convert API config to dynamic type
		var err error
		state.Config, err = convert.StructToDynamic(ctx, p.Config)
		if err != nil {
			diags.AddError(
				"populate policy resource",
				fmt.Sprintf("policy %d: failed to convert config: %s", id, err.Error()),
			)
			return state, diags
		}
	}

	return state, diags
}

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

	if !plan.AssociatedResourceId.IsNull() && !plan.AssociatedResourceId.IsUnknown() {
		addPolicy.SetRefId(plan.AssociatedResourceId.ValueInt64())
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

	// Set PolicyType - required field
	policyTypeCode := plan.PolicyType.Code.ValueString()

	policyType := sdk.NewAddPoliciesRequestPolicyPolicyTypeWithDefaults()
	policyType.SetCode(policyTypeCode)
	addPolicy.SetPolicyType(*policyType)

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

	// Set RefType if provided and not "Global"
	// When associated_resource_type is "Global", we don't set RefType (leave it null)
	refType := plan.AssociatedResourceType.ValueString()
	if refType != "Global" {
		// Convert user-facing resource type to API type
		apiType := resourceTypeToAPIType(refType)
		addPolicy.SetRefType(apiType)
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

	if !plan.AssociatedResourceId.IsNull() && !plan.AssociatedResourceId.IsUnknown() {
		updatePolicy.SetRefId(plan.AssociatedResourceId.ValueInt64())
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

	// Set Config - convert dynamic to SDK config structure
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configMap, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"update policy resource",
				fmt.Sprintf("policy %d: failed to convert config: %s", id, err.Error()),
			)
			return
		}

		// Check if config is empty
		if configMapTyped, ok := configMap.(map[string]interface{}); ok && len(configMapTyped) == 0 {
			resp.Diagnostics.AddError(
				"update policy resource",
				fmt.Sprintf("policy %d: config cannot be empty. "+
					"Please provide the required configuration fields for this policy type.", id),
			)
			return
		}

		// Marshal to JSON then unmarshal to SDK config structure
		// This allows the SDK's UnmarshalJSON to handle the oneOf structure
		configJSON, err := json.Marshal(configMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"update policy resource",
				fmt.Sprintf("policy %d: failed to marshal config to JSON: %s", id, err.Error()),
			)
			return
		}

		var sdkConfig sdk.UpdatePoliciesRequestPolicyConfig
		if err := json.Unmarshal(configJSON, &sdkConfig); err != nil {
			resp.Diagnostics.AddError(
				"update policy resource",
				fmt.Sprintf("policy %d: invalid config: %s", id, err.Error()),
			)
			return
		}

		updatePolicy.SetConfig(sdkConfig)
	}

	// Set RefType if provided and not "Global"
	// When associated_resource_type is "Global", we don't set RefType (leave it null)
	refType := plan.AssociatedResourceType.ValueString()
	if refType != "Global" {
		// Convert user-facing resource type to API type
		apiType := resourceTypeToAPIType(refType)
		updatePolicy.SetRefType(apiType)
	}

	updatePolicyRequest := sdk.NewUpdatePoliciesRequest(*updatePolicy)

	policy, hresp, err := client.PoliciesAPI.UpdatePolicies(ctx, id).
		UpdatePoliciesRequest(*updatePolicyRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update policy resource",
			fmt.Sprintf("policy %d PUT failed: ", id)+errors.ErrMsg(err, hresp),
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

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan PolicyModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read policy resource",
			"new client call failed with "+err.Error(),
		)
		return
	}

	id := plan.Id.ValueInt64()
	state, diags := getPolicyAsState(ctx, id, client, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: failed to create client: %s", id, err.Error()),
		)
		return
	}

	_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete policy resource",
			fmt.Sprintf("policy %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)
		return
	}
}

// ValidateConfig validates that associated_resource_id is set when
// associated_resource_type is not "Global"
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config PolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If associated_resource_type is not "Global", associated_resource_id must be set
	resourceType := config.AssociatedResourceType.ValueString()

	if resourceType != "Global" {
		// Check if associated_resource_id is set
		if config.AssociatedResourceId.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("associated_resource_id"),
				"Missing required attribute",
				fmt.Sprintf(
					"associated_resource_id is required when associated_resource_type is '%s'. "+
						"Set associated_resource_id to the ID of the %s resource, or set associated_resource_type to 'Global'.",
					resourceType, resourceType,
				),
			)
		}
	}
}
