package storage_volume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type storageVolumeModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	StorageServerID types.Int64  `tfsdk:"storage_server_id"`
	MaxStorage      types.Int64  `tfsdk:"max_storage"`
	Status          types.String `tfsdk:"status"`
}

func StorageVolumeSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Storage Volume resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the storage volume.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the storage volume.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The storage type code or ID.",
			},
			"storage_server_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The ID of the storage server.",
			},
			"max_storage": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum storage size in bytes.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the storage volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
