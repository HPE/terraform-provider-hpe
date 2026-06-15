// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &backupHostResource{}
	_ resource.ResourceWithConfigure   = &backupHostResource{}
	_ resource.ResourceWithImportState = &backupHostResource{}
)

const (
	createOperation = "create host backup"
	readOperation   = "read host backup"
	updateOperation = "update host backup"
	deleteOperation = "delete host backup"
)

type backupHostResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &backupHostResource{}
}

func (r *backupHostResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup_host"
}

func (r *backupHostResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupHostResourceSchema(ctx)
}
