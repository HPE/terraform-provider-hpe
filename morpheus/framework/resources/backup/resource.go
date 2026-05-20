package backup

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
	_ resource.Resource                = &backupResource{}
	_ resource.ResourceWithConfigure   = &backupResource{}
	_ resource.ResourceWithImportState = &backupResource{}
)

type backupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &backupResource{}
}

func (r *backupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_backup"
}

func (r *backupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BackupSchema(ctx)
}

func (r *backupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan backupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance := sdk.BackupInstance{
		Name:         plan.Name.ValueString(),
		LocationType: "instance",
		JobAction:    "new",
	}
	if !plan.InstanceID.IsNull() {
		instance.InstanceId = plan.InstanceID.ValueInt64()
	}
	if !plan.ContainerID.IsNull() {
		instance.ContainerId = plan.ContainerID.ValueInt64()
	}
	if !plan.BackupType.IsNull() {
		instance.BackupType = plan.BackupType.ValueString()
	}
	if !plan.ScheduleID.IsNull() {
		instance.JobSchedule = plan.ScheduleID.ValueInt64Pointer()
	}
	if !plan.RetentionCount.IsNull() {
		instance.RetentionCount = plan.RetentionCount.ValueInt64Pointer()
	}

	backup := sdk.AddBackupsRequestBackup{
		BackupInstance: &instance,
	}

	result, httpResp, err := client.BackupsAPI.AddBackups(ctx).AddBackupsRequest(sdk.AddBackupsRequest{
		Backup: backup,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "backup", plan.Name.ValueString(), err, httpResp)
		return
	}

	b := result.GetBackup()
	mapAddResponseToModel(&plan, &b)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state backupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.BackupsAPI.GetBackups(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "backup", "", err, httpResp)
		return
	}

	b := result.GetBackup()
	mapGetResponseToModel(&state, &b)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan backupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateBackupsRequestBackup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	result, httpResp, err := client.BackupsAPI.UpdateBackups(ctx, id).UpdateBackupsRequest(sdk.UpdateBackupsRequest{
		Backup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "backup", plan.Name.ValueString(), err, httpResp)
		return
	}

	b := result.GetBackup()
	mapUpdateResponseToModel(&plan, &b)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state backupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.BackupsAPI.RemoveBackups(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "backup", "", err, httpResp)
		return
	}
}

func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *backupModel, b *sdk.AddBackups200ResponseAllOfBackup) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
	if v := b.RetentionCount.Get(); v != nil {
		model.RetentionCount = types.Int64Value(*v)
	} else {
		model.RetentionCount = types.Int64Null()
	}
}

func mapGetResponseToModel(model *backupModel, b *sdk.GetBackups200ResponseBackup) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
	if v := b.ContainerId.Get(); v != nil {
		model.ContainerID = types.Int64Value(*v)
	} else {
		model.ContainerID = types.Int64Null()
	}
	if v := b.RetentionCount.Get(); v != nil {
		model.RetentionCount = types.Int64Value(*v)
	} else {
		model.RetentionCount = types.Int64Null()
	}
}

func mapUpdateResponseToModel(model *backupModel, b *sdk.UpdateBackups200ResponseAllOfBackup) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
}
