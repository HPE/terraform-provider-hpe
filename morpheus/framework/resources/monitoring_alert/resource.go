package monitoring_alert

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
	_ resource.Resource                = &monitoringAlertResource{}
	_ resource.ResourceWithConfigure   = &monitoringAlertResource{}
	_ resource.ResourceWithImportState = &monitoringAlertResource{}
)

type monitoringAlertResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &monitoringAlertResource{}
}

func (r *monitoringAlertResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_monitoring_alert"
}

func (r *monitoringAlertResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = MonitoringAlertSchema(ctx)
}

func (r *monitoringAlertResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringAlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddAlertsRequestAlert{
		Name: plan.Name.ValueString(),
	}
	if !plan.MinSeverity.IsNull() {
		body.MinSeverity = plan.MinSeverity.ValueStringPointer()
	}
	if !plan.MinDuration.IsNull() {
		dur := int32(plan.MinDuration.ValueInt64()) //nolint:gosec // value range is safe
		body.MinDuration = &dur
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.AllChecks.IsNull() {
		body.AllChecks = plan.AllChecks.ValueBoolPointer()
	}
	if !plan.AllGroups.IsNull() {
		body.AllGroups = plan.AllGroups.ValueBoolPointer()
	}

	result, httpResp, err := client.AlertsAPI.AddAlerts(ctx).AddAlertsRequest(sdk.AddAlertsRequest{
		Alert: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_alert", plan.Name.ValueString(), err, httpResp)

		return
	}

	alert := result.Alert
	mapAlertResponseToModel(&plan, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.AlertsAPI.GetAlerts(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_alert", "", err, httpResp)

		return
	}

	alert := result.Alert
	mapGetAlertResponseToModel(&state, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *monitoringAlertResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringAlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateAlertsRequestAlert{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.MinSeverity.IsNull() {
		body.MinSeverity = plan.MinSeverity.ValueStringPointer()
	}
	if !plan.MinDuration.IsNull() {
		dur := int32(plan.MinDuration.ValueInt64()) //nolint:gosec // value range is safe
		body.MinDuration = &dur
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.AllChecks.IsNull() {
		body.AllChecks = plan.AllChecks.ValueBoolPointer()
	}
	if !plan.AllGroups.IsNull() {
		body.AllGroups = plan.AllGroups.ValueBoolPointer()
	}

	result, httpResp, err := client.AlertsAPI.UpdateAlerts(ctx, id).UpdateAlertsRequest(sdk.UpdateAlertsRequest{
		Alert: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "monitoring_alert", plan.Name.ValueString(), err, httpResp)

		return
	}

	alert := result.Alert
	mapGetAlertResponseToModel(&plan, alert)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringAlertResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.AlertsAPI.DeleteAlerts(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "monitoring_alert", "", err, httpResp)

		return
	}
}

func (r *monitoringAlertResource) ImportState(
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

func mapAlertResponseToModel(model *monitoringAlertModel, alert *sdk.AddAlerts200ResponseAllOfAlert) {
	if alert.Id != nil {
		model.ID = types.Int64Value(*alert.Id)
	}
	if alert.Name != nil {
		model.Name = types.StringValue(*alert.Name)
	}
	if alert.MinSeverity != nil {
		model.MinSeverity = types.StringValue(*alert.MinSeverity)
	}
	if alert.MinDuration != nil {
		model.MinDuration = types.Int64Value(*alert.MinDuration)
	}
	if alert.Active != nil {
		model.Active = types.BoolValue(*alert.Active)
	}
	if alert.AllChecks != nil {
		model.AllChecks = types.BoolValue(*alert.AllChecks)
	}
	if alert.AllGroups != nil {
		model.AllGroups = types.BoolValue(*alert.AllGroups)
	}
}

func mapGetAlertResponseToModel(model *monitoringAlertModel, alert *sdk.GetAlerts200ResponseAllOfAlert) {
	if alert.Id != nil {
		model.ID = types.Int64Value(*alert.Id)
	}
	if alert.Name != nil {
		model.Name = types.StringValue(*alert.Name)
	}
	if alert.MinSeverity != nil {
		model.MinSeverity = types.StringValue(*alert.MinSeverity)
	}
	if alert.MinDuration != nil {
		model.MinDuration = types.Int64Value(*alert.MinDuration)
	}
	if alert.Active != nil {
		model.Active = types.BoolValue(*alert.Active)
	}
	if alert.AllChecks != nil {
		model.AllChecks = types.BoolValue(*alert.AllChecks)
	}
	if alert.AllGroups != nil {
		model.AllGroups = types.BoolValue(*alert.AllGroups)
	}
}
