package backuphost

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
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

func (r *backupHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := sdk.BackupServerHost{
		Name:         plan.Name.ValueString(),
		LocationType: "server",
		BackupType:   plan.BackupTypeCode.ValueString(),
		JobAction:    "addTo",
		ServerId:     plan.HostId.ValueInt64(),
		JobId:        plan.JobId.ValueInt64Pointer(),
	}
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		host.Target = plan.StorageProviderId.ValueInt64Pointer()
	}
	if !plan.Path.IsNull() {
		host.TargetPath = plan.Path.ValueStringPointer()
	}

	backup := sdk.AddBackupsRequestBackup{
		BackupServerHost: &host,
	}

	result, httpResp, err := client.BackupsAPI.AddBackups(ctx).AddBackupsRequest(sdk.AddBackupsRequest{
		Backup: backup,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(createOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	if result.Backup == nil || result.Backup.Id == nil {
		resp.Diagnostics.AddError(createOperation, "Backup ID is nil in the create response")

		return
	}
	id := *result.Backup.Id

	state, diags := getBackupAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var prior BackupHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := prior.Id.ValueInt64()

	state, diags := getBackupAsState(ctx, id, client, prior)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan BackupHostModel
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
	if !plan.Path.IsNull() {
		body.TargetPath = plan.Path.ValueStringPointer()
	}

	_, httpResp, err := client.BackupsAPI.UpdateBackups(ctx, id).UpdateBackupsRequest(sdk.UpdateBackupsRequest{
		Backup: body,
	}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(updateOperation, errfmt.ErrMsg(err, httpResp))

		return
	}

	state, diags := getBackupAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state BackupHostModel
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

func (r *backupHostResource) ImportState(
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
	plan BackupHostModel,
) (BackupHostModel, diag.Diagnostics) {
	var state BackupHostModel
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

	if b.BackupType != nil {
		state.BackupTypeCode = convert.StrToType(b.BackupType.Code)
	}

	if b.Job != nil {
		state.JobId = convert.Int64ToType(b.Job.Id)
	}

	if b.StorageProvider != nil {
		state.StorageProviderId = convert.Int64ToType(b.StorageProvider.Id)
	}

	// host_id and path are not present on every backup type, so fall back to the
	// planned value when the API omits them.
	state.HostId = plan.HostId
	if b.Server != nil {
		state.HostId = convert.Int64ToType(b.Server.Id)
	}

	state.Path = plan.Path
	if v := b.TargetPath.Get(); v != nil {
		state.Path = convert.StrToType(v)
	}

	return state, diags
}
