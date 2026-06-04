package provisioning_license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type provisioningLicenseModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	LicenseType         types.String `tfsdk:"license_type"`
	LicenseKeyWo        types.String `tfsdk:"license_key_wo"`
	LicenseKeyWoVersion types.Int64  `tfsdk:"license_key_wo_version"`
	Description         types.String `tfsdk:"description"`
	VirtualImages       types.List   `tfsdk:"virtual_images"`
	Tenants             types.List   `tfsdk:"tenants"`
}

func ProvisioningLicenseSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Provisioning License resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the provisioning license.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the provisioning license.",
			},
			"license_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the license.",
			},
			"license_key_wo": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "The license key.",
			},
			"license_key_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "License key version. Used to determine if license_key_wo has been updated.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the provisioning license.",
			},
			"virtual_images": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of virtual image IDs associated with the license.",
			},
			"tenants": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
				Description: "List of tenant IDs associated with the license.",
			},
		},
	}
}
