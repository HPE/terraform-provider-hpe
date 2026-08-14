// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Hand-written schema for the hpe_morpheus_storage_controller_type data source.
//
// This lookup has no singular REST endpoint to generate from: it composes a
// controller_mount_point string from a storage controller type id (read from
// /api/provision-types) and the caller's bus and interface numbers. The schema
// is therefore maintained here directly rather than produced by the code-spec
// generator.
//
// It is the hpe equivalent of hpegl_vmaas_instance_storage_controller and
// supplies the controller_mount_point value consumed by volumes on
// hpe_morpheus_instance. hpegl exposes only the composed string (as its id);
// this data source is a superset, exposing both the controller type id and the
// composed controller_mount_point.

package storagecontrollertype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// defaultProvisionTypeCode mirrors hpegl, which resolves storage controller
// types against the VMware provision type. It is exposed as an optional argument
// so non-VMware provision types remain reachable.
const defaultProvisionTypeCode = "vmware"

func StorageControllerTypeDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Composes a controller_mount_point for a Morpheus storage controller type, " +
			"for use with volumes on hpe_morpheus_instance. The mount point has the format " +
			"id:busNumber:typeId:unitNumber, where id is -1 (a new controller).",
		MarkdownDescription: "Composes a `controller_mount_point` for a Morpheus storage controller " +
			"type, for use with volumes on `hpe_morpheus_instance`. The mount point has the format " +
			"`id:busNumber:typeId:unitNumber`, where `id` is `-1` (a new controller).",
		Attributes: map[string]schema.Attribute{
			"controller_name": schema.StringAttribute{
				Required: true,
				Description: "The name of the storage controller type to look up " +
					"(for example \"SCSI VMware Paravirtual\"). Matched case-insensitively " +
					"and whitespace-trimmed.",
				MarkdownDescription: "The name of the storage controller type to look up " +
					"(for example `SCSI VMware Paravirtual`). Matched case-insensitively " +
					"and whitespace-trimmed.",
			},
			"bus_number": schema.Int64Attribute{
				Required:            true,
				Description:         "The controller bus number to embed in the mount point.",
				MarkdownDescription: "The controller bus number to embed in the mount point.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"interface_number": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "The device (unit) number to embed in the mount point. " +
					"Defaults to 0, matching hpegl.",
				MarkdownDescription: "The device (unit) number to embed in the mount point. " +
					"Defaults to `0`, matching hpegl.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"provision_type_code": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The provision type code the storage controller type belongs to. " +
					"Defaults to \"vmware\".",
				MarkdownDescription: "The provision type code the storage controller type belongs to. " +
					"Defaults to `vmware`.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The id of the storage controller type.",
				MarkdownDescription: "The id of the storage controller type.",
			},
			"controller_mount_point": schema.StringAttribute{
				Computed: true,
				Description: "The composed controller mount point, in the format " +
					"id:busNumber:typeId:unitNumber (for example \"-1:1:6:0\").",
				MarkdownDescription: "The composed controller mount point, in the format " +
					"`id:busNumber:typeId:unitNumber` (for example `-1:1:6:0`).",
			},
			"category": schema.StringAttribute{
				Computed:            true,
				Description:         "The category of the storage controller type (for example \"scsi\").",
				MarkdownDescription: "The category of the storage controller type (for example `scsi`).",
			},
			"max_devices": schema.Int64Attribute{
				Computed: true,
				Description: "The base maximum number of devices the controller type supports. " +
					"Informational only: the effective limit can be higher on newer VM " +
					"hardware versions, so it is not enforced here.",
				MarkdownDescription: "The base maximum number of devices the controller type supports. " +
					"Informational only: the effective limit can be higher on newer VM " +
					"hardware versions, so it is not enforced here.",
			},
			"display_order": schema.Int64Attribute{
				Computed:            true,
				Description:         "The display order of the storage controller type in the Morpheus UI.",
				MarkdownDescription: "The display order of the storage controller type in the Morpheus UI.",
			},
		},
	}
}

type StorageControllerTypeModel struct {
	ControllerName       types.String `tfsdk:"controller_name"`
	BusNumber            types.Int64  `tfsdk:"bus_number"`
	InterfaceNumber      types.Int64  `tfsdk:"interface_number"`
	ProvisionTypeCode    types.String `tfsdk:"provision_type_code"`
	Id                   types.Int64  `tfsdk:"id"`
	ControllerMountPoint types.String `tfsdk:"controller_mount_point"`
	Category             types.String `tfsdk:"category"`
	MaxDevices           types.Int64  `tfsdk:"max_devices"`
	DisplayOrder         types.Int64  `tfsdk:"display_order"`
}
