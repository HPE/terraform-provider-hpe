//go:build experimental

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
