package power_schedule

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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &powerScheduleResource{}
	_ resource.ResourceWithConfigure   = &powerScheduleResource{}
	_ resource.ResourceWithImportState = &powerScheduleResource{}
)

type powerScheduleResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &powerScheduleResource{}
}

func (r *powerScheduleResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_power_schedule"
}

func (r *powerScheduleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PowerScheduleSchema(ctx)
}

func (r *powerScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan powerScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddPowerSchedulesRequestSchedule{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ScheduleType.IsNull() && !plan.ScheduleType.IsUnknown() {
		body.ScheduleType = plan.ScheduleType.ValueStringPointer()
	}
	if !plan.ScheduleTimezone.IsNull() && !plan.ScheduleTimezone.IsUnknown() {
		body.ScheduleTimezone = plan.ScheduleTimezone.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.MondayOnTime.IsNull() && !plan.MondayOnTime.IsUnknown() {
		body.MondayOnTime = plan.MondayOnTime.ValueStringPointer()
	}
	if !plan.MondayOffTime.IsNull() && !plan.MondayOffTime.IsUnknown() {
		body.MondayOffTime = plan.MondayOffTime.ValueStringPointer()
	}
	if !plan.TuesdayOnTime.IsNull() && !plan.TuesdayOnTime.IsUnknown() {
		body.TuesdayOnTime = plan.TuesdayOnTime.ValueStringPointer()
	}
	if !plan.TuesdayOffTime.IsNull() && !plan.TuesdayOffTime.IsUnknown() {
		body.TuesdayOffTime = plan.TuesdayOffTime.ValueStringPointer()
	}
	if !plan.WednesdayOnTime.IsNull() && !plan.WednesdayOnTime.IsUnknown() {
		body.WednesdayOnTime = plan.WednesdayOnTime.ValueStringPointer()
	}
	if !plan.WednesdayOffTime.IsNull() && !plan.WednesdayOffTime.IsUnknown() {
		body.WednesdayOffTime = plan.WednesdayOffTime.ValueStringPointer()
	}
	if !plan.ThursdayOnTime.IsNull() && !plan.ThursdayOnTime.IsUnknown() {
		body.ThursdayOnTime = plan.ThursdayOnTime.ValueStringPointer()
	}
	if !plan.ThursdayOffTime.IsNull() && !plan.ThursdayOffTime.IsUnknown() {
		body.ThursdayOffTime = plan.ThursdayOffTime.ValueStringPointer()
	}
	if !plan.FridayOnTime.IsNull() && !plan.FridayOnTime.IsUnknown() {
		body.FridayOnTime = plan.FridayOnTime.ValueStringPointer()
	}
	if !plan.FridayOffTime.IsNull() && !plan.FridayOffTime.IsUnknown() {
		body.FridayOffTime = plan.FridayOffTime.ValueStringPointer()
	}
	if !plan.SaturdayOnTime.IsNull() && !plan.SaturdayOnTime.IsUnknown() {
		body.SaturdayOnTime = plan.SaturdayOnTime.ValueStringPointer()
	}
	if !plan.SaturdayOffTime.IsNull() && !plan.SaturdayOffTime.IsUnknown() {
		body.SaturdayOffTime = plan.SaturdayOffTime.ValueStringPointer()
	}
	if !plan.SundayOnTime.IsNull() && !plan.SundayOnTime.IsUnknown() {
		body.SundayOnTime = plan.SundayOnTime.ValueStringPointer()
	}
	if !plan.SundayOffTime.IsNull() && !plan.SundayOffTime.IsUnknown() {
		body.SundayOffTime = plan.SundayOffTime.ValueStringPointer()
	}

	reqBody := sdk.AddPowerSchedulesRequest{Schedule: body}
	result, httpResp, err := client.AutomationAPI.AddPowerSchedules(ctx).AddPowerSchedulesRequest(reqBody).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "power_schedule", plan.Name.ValueString(), err, httpResp)

		return
	}

	var id int64
	if result.Schedule != nil && result.Schedule.Id != nil {
		id = *result.Schedule.Id
	}

	readResult, httpResp, err := client.AutomationAPI.GetPowerSchedules(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "power_schedule", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "power_schedule",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.Schedule == nil {
		resp.Diagnostics.AddError("API returned nil", "Schedule is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.Schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *powerScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state powerScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.AutomationAPI.GetPowerSchedules(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "power_schedule", "", err, httpResp)

		return
	}

	schedule := result.Schedule
	if schedule == nil {
		resp.Diagnostics.AddError("API returned nil", "Schedule is nil in the response")

		return
	}

	mapGetResponseToModel(&state, schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *powerScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan powerScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdatePowerSchedulesRequestSchedule{
		Name: plan.Name.ValueStringPointer(),
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
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.MondayOnTime.IsNull() {
		body.MondayOnTime = plan.MondayOnTime.ValueStringPointer()
	}
	if !plan.MondayOffTime.IsNull() {
		body.MondayOffTime = plan.MondayOffTime.ValueStringPointer()
	}
	if !plan.TuesdayOnTime.IsNull() {
		body.TuesdayOnTime = plan.TuesdayOnTime.ValueStringPointer()
	}
	if !plan.TuesdayOffTime.IsNull() {
		body.TuesdayOffTime = plan.TuesdayOffTime.ValueStringPointer()
	}
	if !plan.WednesdayOnTime.IsNull() {
		body.WednesdayOnTime = plan.WednesdayOnTime.ValueStringPointer()
	}
	if !plan.WednesdayOffTime.IsNull() {
		body.WednesdayOffTime = plan.WednesdayOffTime.ValueStringPointer()
	}
	if !plan.ThursdayOnTime.IsNull() {
		body.ThursdayOnTime = plan.ThursdayOnTime.ValueStringPointer()
	}
	if !plan.ThursdayOffTime.IsNull() {
		body.ThursdayOffTime = plan.ThursdayOffTime.ValueStringPointer()
	}
	if !plan.FridayOnTime.IsNull() {
		body.FridayOnTime = plan.FridayOnTime.ValueStringPointer()
	}
	if !plan.FridayOffTime.IsNull() {
		body.FridayOffTime = plan.FridayOffTime.ValueStringPointer()
	}
	if !plan.SaturdayOnTime.IsNull() {
		body.SaturdayOnTime = plan.SaturdayOnTime.ValueStringPointer()
	}
	if !plan.SaturdayOffTime.IsNull() {
		body.SaturdayOffTime = plan.SaturdayOffTime.ValueStringPointer()
	}
	if !plan.SundayOnTime.IsNull() {
		body.SundayOnTime = plan.SundayOnTime.ValueStringPointer()
	}
	if !plan.SundayOffTime.IsNull() {
		body.SundayOffTime = plan.SundayOffTime.ValueStringPointer()
	}

	_, httpResp, err := client.AutomationAPI.UpdatePowerSchedules(ctx, id).
		UpdatePowerSchedulesRequest(sdk.UpdatePowerSchedulesRequest{
			Schedule: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "power_schedule", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.AutomationAPI.GetPowerSchedules(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "power_schedule", plan.Name.ValueString(), err, httpResp)

		return
	}

	if readResult.Schedule == nil {
		resp.Diagnostics.AddError("API returned nil", "Schedule is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.Schedule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *powerScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state powerScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.AutomationAPI.RemovePowerSchedules(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "power_schedule", "", err, httpResp)

		return
	}
}

func (r *powerScheduleResource) ImportState(
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

func mapAddResponseToModel(model *powerScheduleModel, schedule *sdk.AddPowerSchedules200ResponseAllOfSchedule) {
	if schedule.Id != nil {
		model.ID = types.Int64Value(*schedule.Id)
	}
	if schedule.Name != nil {
		model.Name = types.StringValue(*schedule.Name)
	}
	if desc := schedule.Description.Get(); desc != nil {
		model.Description = types.StringValue(*desc)
	} else {
		model.Description = types.StringNull()
	}
	if schedule.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*schedule.ScheduleType)
	}
	if schedule.ScheduleTimezone != nil {
		model.ScheduleTimezone = types.StringValue(*schedule.ScheduleTimezone)
	}
	if schedule.Enabled != nil {
		model.Enabled = types.BoolValue(*schedule.Enabled)
	}
	mapTimeFields(model,
		schedule.MondayOnTime, schedule.MondayOffTime,
		schedule.TuesdayOnTime, schedule.TuesdayOffTime,
		schedule.WednesdayOnTime, schedule.WednesdayOffTime,
		schedule.ThursdayOnTime, schedule.ThursdayOffTime,
		schedule.FridayOnTime, schedule.FridayOffTime,
		schedule.SaturdayOnTime, schedule.SaturdayOffTime,
		schedule.SundayOnTime, schedule.SundayOffTime,
	)
	if schedule.TotalMonthlyHoursSaved != nil {
		model.TotalMonthlyHoursSaved = types.Float64Value(float64(*schedule.TotalMonthlyHoursSaved))
	} else {
		model.TotalMonthlyHoursSaved = types.Float64Null()
	}
}

func mapGetResponseToModel(model *powerScheduleModel, schedule *sdk.GetPowerSchedules200ResponseAllOfSchedule) {
	if schedule.Id != nil {
		model.ID = types.Int64Value(*schedule.Id)
	}
	if schedule.Name != nil {
		model.Name = types.StringValue(*schedule.Name)
	}
	if desc := schedule.Description.Get(); desc != nil {
		model.Description = types.StringValue(*desc)
	} else {
		model.Description = types.StringNull()
	}
	if schedule.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*schedule.ScheduleType)
	}
	if schedule.ScheduleTimezone != nil {
		model.ScheduleTimezone = types.StringValue(*schedule.ScheduleTimezone)
	}
	if schedule.Enabled != nil {
		model.Enabled = types.BoolValue(*schedule.Enabled)
	}
	mapTimeFields(model,
		schedule.MondayOnTime, schedule.MondayOffTime,
		schedule.TuesdayOnTime, schedule.TuesdayOffTime,
		schedule.WednesdayOnTime, schedule.WednesdayOffTime,
		schedule.ThursdayOnTime, schedule.ThursdayOffTime,
		schedule.FridayOnTime, schedule.FridayOffTime,
		schedule.SaturdayOnTime, schedule.SaturdayOffTime,
		schedule.SundayOnTime, schedule.SundayOffTime,
	)
	if schedule.TotalMonthlyHoursSaved != nil {
		model.TotalMonthlyHoursSaved = types.Float64Value(float64(*schedule.TotalMonthlyHoursSaved))
	} else {
		model.TotalMonthlyHoursSaved = types.Float64Null()
	}
}

func mapUpdateResponseToModel(model *powerScheduleModel, schedule *sdk.UpdatePowerSchedules200ResponseAllOfSchedule) {
	if schedule.Id != nil {
		model.ID = types.Int64Value(*schedule.Id)
	}
	if schedule.Name != nil {
		model.Name = types.StringValue(*schedule.Name)
	}
	if desc := schedule.Description.Get(); desc != nil {
		model.Description = types.StringValue(*desc)
	} else {
		model.Description = types.StringNull()
	}
	if schedule.ScheduleType != nil {
		model.ScheduleType = types.StringValue(*schedule.ScheduleType)
	}
	if schedule.ScheduleTimezone != nil {
		model.ScheduleTimezone = types.StringValue(*schedule.ScheduleTimezone)
	}
	if schedule.Enabled != nil {
		model.Enabled = types.BoolValue(*schedule.Enabled)
	}
	mapTimeFields(model,
		schedule.MondayOnTime, schedule.MondayOffTime,
		schedule.TuesdayOnTime, schedule.TuesdayOffTime,
		schedule.WednesdayOnTime, schedule.WednesdayOffTime,
		schedule.ThursdayOnTime, schedule.ThursdayOffTime,
		schedule.FridayOnTime, schedule.FridayOffTime,
		schedule.SaturdayOnTime, schedule.SaturdayOffTime,
		schedule.SundayOnTime, schedule.SundayOffTime,
	)
	if schedule.TotalMonthlyHoursSaved != nil {
		model.TotalMonthlyHoursSaved = types.Float64Value(float64(*schedule.TotalMonthlyHoursSaved))
	} else {
		model.TotalMonthlyHoursSaved = types.Float64Null()
	}
}

func mapTimeFields(
	model *powerScheduleModel,
	mondayOn,
	mondayOff,
	tuesdayOn,
	tuesdayOff,
	wednesdayOn,
	wednesdayOff,
	thursdayOn,
	thursdayOff,
	fridayOn,
	fridayOff,
	saturdayOn,
	saturdayOff,
	sundayOn,
	sundayOff *string,
) {
	model.MondayOnTime = ptrToStringValue(mondayOn)
	model.MondayOffTime = ptrToStringValue(mondayOff)
	model.TuesdayOnTime = ptrToStringValue(tuesdayOn)
	model.TuesdayOffTime = ptrToStringValue(tuesdayOff)
	model.WednesdayOnTime = ptrToStringValue(wednesdayOn)
	model.WednesdayOffTime = ptrToStringValue(wednesdayOff)
	model.ThursdayOnTime = ptrToStringValue(thursdayOn)
	model.ThursdayOffTime = ptrToStringValue(thursdayOff)
	model.FridayOnTime = ptrToStringValue(fridayOn)
	model.FridayOffTime = ptrToStringValue(fridayOff)
	model.SaturdayOnTime = ptrToStringValue(saturdayOn)
	model.SaturdayOffTime = ptrToStringValue(saturdayOff)
	model.SundayOnTime = ptrToStringValue(sundayOn)
	model.SundayOffTime = ptrToStringValue(sundayOff)
}

func ptrToStringValue(s *string) types.String {
	if s != nil {
		return types.StringValue(*s)
	}

	return types.StringNull()
}
