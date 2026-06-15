// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

func (r *backupInstanceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// instance_id is required
	instanceID := plan.InstanceId.ValueInt64()

	instance := sdk.BackupInstance{
		Name:         plan.Name.ValueString(),
		InstanceId:   instanceID,
		LocationType: "instance",
		JobAction:    "addTo",
		JobId:        plan.JobId.ValueInt64Pointer(),
	}
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		instance.Target = plan.StorageProviderId.ValueInt64Pointer()
	}

	// container_id is optional/computed. Use the user-supplied value when set;
	// otherwise resolve it from the instance's container list (the same thing
	// the API does on a dry-run create to /api/backups/create).
	if !plan.ContainerId.IsNull() && !plan.ContainerId.IsUnknown() {
		instance.ContainerId = plan.ContainerId.ValueInt64()
	} else {
		instResult, instHResp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if err != nil || instHResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, instHResp))

			return
		}

		if instResult.Instance == nil {
			resp.Diagnostics.AddError(
				createOperation,
				fmt.Sprintf("instance %d returned no instance", instanceID),
			)

			return
		}

		if len(instResult.Instance.Containers) == 0 {
			resp.Diagnostics.AddError(
				createOperation,
				fmt.Sprintf("instance %d has no containers", instanceID),
			)

			return
		}

		instance.ContainerId = instResult.Instance.Containers[0]
	}

	backup := sdk.AddBackupsRequestBackup{
		BackupInstance: &instance,
	}

	result, httpResp, err := client.BackupsAPI.AddBackups(ctx).AddBackupsRequest(sdk.AddBackupsRequest{
		Backup: backup,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	if result.Backup == nil || result.Backup.Id == nil {
		resp.Diagnostics.AddError("API returned nil", "Backup ID is nil in the create response")

		return
	}
	id := *result.Backup.Id

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		// The backup was created, but reading it back failed. Taint the resource
		// so the created backup is not leaked from Terraform's perspective.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "backup_instance",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state BackupInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupInstanceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateBackupsRequestBackup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.JobId.IsNull() {
		body.BackupJobId = plan.JobId.ValueInt64Pointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		body.StorageProviderId = *sdk.NewNullableInt64(plan.StorageProviderId.ValueInt64Pointer())
	}

	_, httpResp, err := client.BackupsAPI.UpdateBackups(ctx, id).UpdateBackupsRequest(sdk.UpdateBackupsRequest{
		Backup: body,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(updateOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	state, diags := getBackupAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupInstanceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state BackupInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.BackupsAPI.RemoveBackups(ctx, id).Execute()
	// A 404 means the backup is already gone, so treat it as a successful
	// delete rather than an error.
	if errfmt.IsNotFound(httpResp) {
		return
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(deleteOperation, errfmt.ErrMsg(err, httpResp))

		return
	}
}

func (r *backupInstanceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// getBackupAsState performs a read of the backup by ID and returns the
// resulting state. It is used by Create, Read and Update so that state is
// always populated from a fresh read rather than from the create/update
// responses.
func getBackupAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (BackupInstanceModel, diag.Diagnostics) {
	var state BackupInstanceModel
	var diags diag.Diagnostics

	result, hresp, err := client.BackupsAPI.GetBackups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(readOperation, errfmt.ErrMsg(err, hresp))

		return state, diags
	}

	b := result.Backup
	if b == nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("backup %d GET returned no backup", id),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(b.Id)
	state.Name = convert.StrToType(b.Name)
	state.Enabled = convert.BoolToType(b.Enabled)

	if v := b.ContainerId.Get(); v != nil {
		state.ContainerId = convert.Int64ToType(v)
	}

	if b.BackupType != nil {
		state.BackupTypeCode = convert.StrToType(b.BackupType.Code)
	}

	if b.Instance != nil {
		state.InstanceId = convert.Int64ToType(b.Instance.Id)
	}

	if b.Job != nil {
		state.JobId = convert.Int64ToType(b.Job.Id)
	}

	if b.StorageProvider != nil {
		state.StorageProviderId = convert.Int64ToType(b.StorageProvider.Id)
	}

	return state, diags
}
