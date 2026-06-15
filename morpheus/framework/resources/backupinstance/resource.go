// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &backupInstanceResource{}
	_ resource.ResourceWithConfigure   = &backupInstanceResource{}
	_ resource.ResourceWithImportState = &backupInstanceResource{}
)

const (
	createOperation = "create instance backup"
	readOperation   = "read instance backup"
	updateOperation = "update instance backup"
	deleteOperation = "delete instance backup"
)

type backupInstanceResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &backupInstanceResource{}
}

func (r *backupInstanceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup_instance"
}

func (r *backupInstanceResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = BackupInstanceResourceSchema(ctx)
}
