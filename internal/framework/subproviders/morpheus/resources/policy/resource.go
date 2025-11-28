// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
	_ resource.ResourceWithModifyPlan     = &Resource{}
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
		// For other types (User, Role, Network, Plan), pass through as-is
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
		// For other types (User, Role, Network, Plan), pass through as-is
		return apiType
	}
}

// ValidateConfig validates that associated_resource_id is set when
// associated_resource_type is not "Global" and that each_user is only set
// when associated_resource_type is "Role"
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

	// Validate each_user is only set when associated_resource_type is "Role"
	if !config.EachUser.IsNull() {
		if resourceType != "Role" {
			resp.Diagnostics.AddAttributeError(
				path.Root("each_user"),
				"Invalid attribute combination",
				fmt.Sprintf(
					"each_user can only be set when associated_resource_type is 'Role'. "+
						"Current associated_resource_type is '%s'. "+
						"Either remove each_user or set associated_resource_type to 'Role'.",
					resourceType,
				),
			)
		}
	}
}

// ModifyPlan validates policy type compatibility during the plan phase when the provider is configured
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Only run validation if we have a plan (not during destroy)
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan PolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate policy type is compatible with the associated resource type
	if plan.PolicyType.IsNull() || plan.PolicyType.IsUnknown() {
		return
	}

	policyTypeCode := plan.PolicyType.Code.ValueString()
	resourceType := plan.AssociatedResourceType.ValueString()

	// Get API client - provider is configured at this point
	client, err := r.NewClient(ctx)
	if err != nil {
		// If we can't get a client, skip validation - will be caught during apply
		return
	}

	// Fetch policy types from the API
	policyTypesResp, httpResp, err := client.OptionsAPI.GetOptionSourceData(ctx, "policyTypes").Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != 200 {
		return
	}

	// Find the matching policy type
	var matchingPolicyType map[string]interface{}
	if policyTypesResp != nil && policyTypesResp.Data != nil {
		for _, pt := range policyTypesResp.Data {
			if code, ok := pt["code"].(string); ok && code == policyTypeCode {
				matchingPolicyType = pt
				break
			}
		}
	}

	if matchingPolicyType == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("policy_type").AtName("code"),
			"Invalid policy type",
			fmt.Sprintf("Policy type with code '%s' not found", policyTypeCode),
		)
		return
	}

	// Map resource type to the corresponding "allowOn" field
	var allowFieldName string
	switch resourceType {
	case "Global":
		allowFieldName = "allowOnGlobal"
	case "Group":
		allowFieldName = "allowOnSite"
	case "Cloud":
		allowFieldName = "allowOnZone"
	case "User":
		allowFieldName = "allowOnUser"
	case "Role":
		allowFieldName = "allowOnRole"
	case "Network":
		allowFieldName = "allowOnNetwork"
	case "Plan":
		allowFieldName = "allowOnPlan"
	case "Label":
		allowFieldName = "allowOnLabel"
	default:
		// Unknown resource type - should not happen due to schema validation
		resp.Diagnostics.AddAttributeError(
			path.Root("associated_resource_type"),
			"Unknown resource type",
			fmt.Sprintf("Resource type '%s' is not recognized", resourceType),
		)
		return
	}

	// Check if the policy type allows this resource type
	// Treat Global as always allowed (API may return nil for allowOnGlobal)
	allowed := false
	if resourceType == "Global" {
		allowed = true
	} else if allowValue, ok := matchingPolicyType[allowFieldName]; ok {
		if allowBool, ok := allowValue.(bool); ok {
			allowed = allowBool
		}
	}

	// Validate tenants field is only set when allowOnTenant is true
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		allowOnTenant := false
		if allowValue, ok := matchingPolicyType["allowOnTenant"]; ok {
			if allowBool, ok := allowValue.(bool); ok {
				allowOnTenant = allowBool
			}
		}

		if !allowOnTenant {
			policyTypeName := ""
			if name, ok := matchingPolicyType["name"].(string); ok {
				policyTypeName = name
			}

			resp.Diagnostics.AddAttributeError(
				path.Root("tenants"),
				"Tenants not supported for this policy type",
				fmt.Sprintf(
					"Policy type '%s' (%s) does not support specifying tenants. "+
						"Remove the tenants attribute or choose a different policy type.",
					policyTypeName, policyTypeCode,
				),
			)
		}
	}

	if !allowed {
		policyTypeName := ""
		if name, ok := matchingPolicyType["name"].(string); ok {
			policyTypeName = name
		}

		// Build a list of allowed scopes for this policy type
		var allowedScopes []string
		scopeMapping := map[string]string{
			"allowOnGlobal":  "Global",
			"allowOnSite":    "Group",
			"allowOnZone":    "Cloud",
			"allowOnUser":    "User",
			"allowOnRole":    "Role",
			"allowOnNetwork": "Network",
			"allowOnPlan":    "Plan",
			"allowOnLabel":   "Label",
		}

		for field, scope := range scopeMapping {
			// Always treat Global as allowed (API may return nil for allowOnGlobal)
			if scope == "Global" {
				if !contains(allowedScopes, scope) {
					allowedScopes = append(allowedScopes, scope)
				}
			} else if val, ok := matchingPolicyType[field].(bool); ok && val {
				if !contains(allowedScopes, scope) {
					allowedScopes = append(allowedScopes, scope)
				}
			}
		}

		availableMsg := ""
		if len(allowedScopes) > 0 {
			availableMsg = fmt.Sprintf("\n\nThis policy type can be applied to the following resource types: %v", allowedScopes)
		}

		resp.Diagnostics.AddAttributeError(
			path.Root("associated_resource_type"),
			"Incompatible policy type and resource type",
			fmt.Sprintf(
				"Policy type '%s' (%s) cannot be applied to resource type '%s'. "+
					"This policy type does not support the selected scope.%s",
				policyTypeName, policyTypeCode, resourceType, availableMsg,
			),
		)
	}
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
