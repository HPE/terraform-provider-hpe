package backup_job

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &backupJobResource{}
	_ resource.ResourceWithConfigure   = &backupJobResource{}
	_ resource.ResourceWithImportState = &backupJobResource{}
)

type backupJobResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &backupJobResource{}
}

func (r *backupJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup_job"
}

func (r *backupJobResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupJobSchema(ctx)
}

func (r *backupJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan backupJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddBackupJobsRequestJob{
		Name: plan.Name.ValueString(),
	}
	if !plan.Code.IsNull() {
		body.Code = plan.Code.ValueStringPointer()
	}
	if !plan.RetentionCount.IsNull() {
		body.RetentionCount = plan.RetentionCount.ValueInt64Pointer()
	}
	if !plan.ScheduleID.IsNull() {
		scheduleID := plan.ScheduleID.ValueInt64()
		body.ScheduleId = *sdk.NewNullableInt64(&scheduleID)
	}

	result, httpResp, err := client.BackupsAPI.AddBackupJobs(ctx).AddBackupJobsRequest(sdk.AddBackupJobsRequest{
		Job: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "backup_job", plan.Name.ValueString(), err, httpResp)
		return
	}

	job := result.GetJob()
	mapAddResponseToModel(&plan, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state backupJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.BackupsAPI.GetBackupJobs(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "backup_job", "", err, httpResp)
		return
	}

	job := result.GetJob()
	mapGetResponseToModel(&state, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan backupJobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateBackupJobsRequestJob{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Code.IsNull() {
		body.Code = plan.Code.ValueStringPointer()
	}
	if !plan.RetentionCount.IsNull() {
		body.RetentionCount = plan.RetentionCount.ValueInt64Pointer()
	}
	if !plan.ScheduleID.IsNull() {
		scheduleID := plan.ScheduleID.ValueInt64()
		body.ScheduleId = *sdk.NewNullableInt64(&scheduleID)
	}

	result, httpResp, err := client.BackupsAPI.UpdateBackupJobs(ctx, id).UpdateBackupJobsRequest(sdk.UpdateBackupJobsRequest{
		Job: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "backup_job", plan.Name.ValueString(), err, httpResp)
		return
	}

	job := result.GetJob()
	mapUpdateResponseToModel(&plan, &job)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state backupJobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.BackupsAPI.RemoveBackupJobs(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "backup_job", "", err, httpResp)
		return
	}
}

func (r *backupJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *backupJobModel, job *sdk.AddBackupJobs200ResponseAllOfJob) {
	if job.Id != nil {
		model.ID = types.Int64Value(*job.Id)
	}
	if job.Name != nil {
		model.Name = types.StringValue(*job.Name)
	}
	if v := job.RetentionCount.Get(); v != nil {
		model.RetentionCount = types.Int64Value(*v)
	} else {
		model.RetentionCount = types.Int64Null()
	}
	if job.Enabled != nil {
		model.Enabled = types.BoolValue(*job.Enabled)
	}
}

func mapGetResponseToModel(model *backupJobModel, job *sdk.GetBackupJobs200ResponseJob) {
	if job.Id != nil {
		model.ID = types.Int64Value(*job.Id)
	}
	if job.Name != nil {
		model.Name = types.StringValue(*job.Name)
	}
	if v := job.RetentionCount.Get(); v != nil {
		model.RetentionCount = types.Int64Value(*v)
	} else {
		model.RetentionCount = types.Int64Null()
	}
	if job.Enabled != nil {
		model.Enabled = types.BoolValue(*job.Enabled)
	}
}

func mapUpdateResponseToModel(model *backupJobModel, job *sdk.UpdateBackupJobs200ResponseAllOfJob) {
	if job.Id != nil {
		model.ID = types.Int64Value(*job.Id)
	}
	if job.Name != nil {
		model.Name = types.StringValue(*job.Name)
	}
	if job.Enabled != nil {
		model.Enabled = types.BoolValue(*job.Enabled)
	}
}
