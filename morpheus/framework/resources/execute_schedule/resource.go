package execute_schedule

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
	_ resource.Resource                = &executeScheduleResource{}
	_ resource.ResourceWithConfigure   = &executeScheduleResource{}
	_ resource.ResourceWithImportState = &executeScheduleResource{}
)

type executeScheduleResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &executeScheduleResource{}
}

func (r *executeScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_execute_schedule"
}

func (r *executeScheduleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ExecuteScheduleSchema(ctx)
}

func (r *executeScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan executeScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddExecuteSchedulesRequestSchedule{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ScheduleType.IsNull() {
		body.ScheduleType = plan.ScheduleType.ValueStringPointer()
	}
	if !plan.ScheduleTimezone.IsNull() {
		body.ScheduleTimezone = plan.ScheduleTimezone.ValueStringPointer()
	}
	if !plan.Cron.IsNull() {
		body.Cron = plan.Cron.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	result, httpResp, err := client.AutomationAPI.AddExecuteSchedules(ctx).AddExecuteSchedulesRequest(sdk.AddExecuteSchedulesRequest{
		Schedule: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "execute_schedule", plan.Name.ValueString(), err, httpResp)
		return
	}

	schedule := result.GetSchedule()
	mapAddResponseToModel(&plan, &schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *executeScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state executeScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.AutomationAPI.GetExecuteSchedules(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "execute_schedule", "", err, httpResp)
		return
	}

	schedule := result.GetSchedule()
	mapGetResponseToModel(&state, &schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *executeScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan executeScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateExecuteSchedulesRequestSchedule{}
	if !plan.Name.IsNull() {
		body.Name = plan.Name.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ScheduleType.IsNull() {
		body.ScheduleType = plan.ScheduleType.ValueStringPointer()
	}
	if !plan.ScheduleTimezone.IsNull() {
		body.ScheduleTimezone = plan.ScheduleTimezone.ValueStringPointer()
	}
	if !plan.Cron.IsNull() {
		body.Cron = plan.Cron.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	result, httpResp, err := client.AutomationAPI.UpdateExecuteSchedules(ctx, id).UpdateExecuteSchedulesRequest(sdk.UpdateExecuteSchedulesRequest{
		Schedule: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "execute_schedule", plan.Name.ValueString(), err, httpResp)
		return
	}

	schedule := result.GetSchedule()
	mapUpdateResponseToModel(&plan, &schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *executeScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state executeScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.AutomationAPI.RemoveExecuteSchedules(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "execute_schedule", "", err, httpResp)
		return
	}
}

func (r *executeScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *executeScheduleModel, s *sdk.AddExecuteSchedules200ResponseAllOfSchedule) {
	if s.Id != nil {
		model.ID = types.Int64Value(*s.Id)
	}
	if s.Name != nil {
		model.Name = types.StringValue(*s.Name)
	}
	if v := s.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if s.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*s.ScheduleType)
	}
	if v := s.ScheduleTimezone.Get(); v != nil {
		model.ScheduleTimezone = types.StringValue(*v)
	}
	if s.Cron != nil {
		model.Cron = types.StringValue(*s.Cron)
	}
	if s.Enabled != nil {
		model.Enabled = types.BoolValue(*s.Enabled)
	}
}

func mapGetResponseToModel(model *executeScheduleModel, s *sdk.GetExecuteSchedules200ResponseAllOfSchedule) {
	if s.Id != nil {
		model.ID = types.Int64Value(*s.Id)
	}
	if s.Name != nil {
		model.Name = types.StringValue(*s.Name)
	}
	if v := s.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if s.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*s.ScheduleType)
	}
	if v := s.ScheduleTimezone.Get(); v != nil {
		model.ScheduleTimezone = types.StringValue(*v)
	}
	if s.Cron != nil {
		model.Cron = types.StringValue(*s.Cron)
	}
	if s.Enabled != nil {
		model.Enabled = types.BoolValue(*s.Enabled)
	}
}

func mapUpdateResponseToModel(model *executeScheduleModel, s *sdk.UpdateExecuteSchedules200ResponseAllOfSchedule) {
	if s.Id != nil {
		model.ID = types.Int64Value(*s.Id)
	}
	if s.Name != nil {
		model.Name = types.StringValue(*s.Name)
	}
	if v := s.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if s.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*s.ScheduleType)
	}
	if v := s.ScheduleTimezone.Get(); v != nil {
		model.ScheduleTimezone = types.StringValue(*v)
	}
	if s.Cron != nil {
		model.Cron = types.StringValue(*s.Cron)
	}
	if s.Enabled != nil {
		model.Enabled = types.BoolValue(*s.Enabled)
	}
}
