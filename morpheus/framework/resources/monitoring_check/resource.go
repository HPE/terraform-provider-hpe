package monitoring_check

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
	_ resource.Resource                = &monitoringCheckResource{}
	_ resource.ResourceWithConfigure   = &monitoringCheckResource{}
	_ resource.ResourceWithImportState = &monitoringCheckResource{}
)

type monitoringCheckResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &monitoringCheckResource{}
}

func (r *monitoringCheckResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_monitoring_check"
}

func (r *monitoringCheckResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = MonitoringCheckSchema(ctx)
}

func (r *monitoringCheckResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringCheckModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkBody := sdk.WebCheck{
		Name: plan.Name.ValueString(),
	}
	if !plan.CheckTypeID.IsNull() && !plan.CheckTypeID.IsUnknown() {
		// Look up the check type code from the ID since the API requires code
		ctResult, ctHTTPResp, ctErr := client.ChecksAPI.GetCheckTypes(ctx, plan.CheckTypeID.ValueInt64()).Execute()
		if ctErr != nil || (ctHTTPResp != nil && ctHTTPResp.StatusCode >= 400) {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_check", plan.Name.ValueString(),
				fmt.Errorf("failed to look up check type %d", plan.CheckTypeID.ValueInt64()), ctHTTPResp)

			return
		}
		ct := ctResult.CheckType
		code := ct.Code
		checkBody.CheckType = &sdk.WebCheckAllOfCheckType{
			Code: code,
		}
	}
	if !plan.Description.IsNull() {
		desc := sdk.NewNullableString(plan.Description.ValueStringPointer())
		checkBody.Description = *desc
	}
	if !plan.CheckInterval.IsNull() {
		interval := int32(plan.CheckInterval.ValueInt64()) //nolint:gosec // value range is safe
		checkBody.CheckInterval = &interval
	}
	if !plan.InUptime.IsNull() {
		checkBody.InUptime = plan.InUptime.ValueBoolPointer()
	}
	if !plan.Active.IsNull() {
		checkBody.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Severity.IsNull() {
		checkBody.Severity = plan.Severity.ValueStringPointer()
	}

	checkReq := sdk.AddChecksRequestCheck{WebCheck: &checkBody}

	result, httpResp, err := client.ChecksAPI.AddChecks(ctx).AddChecksRequest(sdk.AddChecksRequest{
		Check: checkReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_check", plan.Name.ValueString(), err, httpResp)

		return
	}

	check := result.Check
	if check == nil {
		resp.Diagnostics.AddError("API returned nil", "MonitoringCheck is nil in the response")
		return
	}
	plan.ID = types.Int64Value(*check.Id)
	if check.Name != nil {
		plan.Name = types.StringValue(*check.Name)
	}
	if check.CheckInterval.IsSet() {
		plan.CheckInterval = types.Int64Value(*check.CheckInterval.Get())
	}
	if check.InUptime != nil {
		plan.InUptime = types.BoolValue(*check.InUptime)
	}
	if check.Active != nil {
		plan.Active = types.BoolValue(*check.Active)
	}
	if check.Severity != nil {
		plan.Severity = types.StringValue(*check.Severity)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.ChecksAPI.GetChecks(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_check", "", err, httpResp)

		return
	}

	check := result.Check
	if check == nil {
		resp.Diagnostics.AddError("API returned nil", "MonitoringCheck is nil in the response")
		return
	}
	if check.Id != nil {
		state.ID = types.Int64Value(*check.Id)
	}
	if check.Name != nil {
		state.Name = types.StringValue(*check.Name)
	}
	if check.Description.IsSet() && check.Description.Get() != nil {
		state.Description = types.StringValue(*check.Description.Get())
	} else {
		state.Description = types.StringNull()
	}
	if check.CheckInterval.IsSet() {
		state.CheckInterval = types.Int64Value(*check.CheckInterval.Get())
	}
	if check.InUptime != nil {
		state.InUptime = types.BoolValue(*check.InUptime)
	}
	if check.Active != nil {
		state.Active = types.BoolValue(*check.Active)
	}
	if check.Severity != nil {
		state.Severity = types.StringValue(*check.Severity)
	}
	if check.CheckType != nil && check.CheckType.Id != nil {
		state.CheckTypeID = types.Int64Value(*check.CheckType.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *monitoringCheckResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringCheckModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	checkBody := sdk.WebCheck1{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		desc := sdk.NewNullableString(plan.Description.ValueStringPointer())
		checkBody.Description = *desc
	}
	if !plan.CheckInterval.IsNull() {
		interval := int32(plan.CheckInterval.ValueInt64()) //nolint:gosec // value range is safe
		checkBody.CheckInterval = &interval
	}
	if !plan.InUptime.IsNull() {
		checkBody.InUptime = plan.InUptime.ValueBoolPointer()
	}
	if !plan.Active.IsNull() {
		checkBody.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Severity.IsNull() {
		checkBody.Severity = plan.Severity.ValueStringPointer()
	}

	checkReq := sdk.UpdateChecksRequestCheck{WebCheck1: &checkBody}

	result, httpResp, err := client.ChecksAPI.UpdateChecks(ctx, id).UpdateChecksRequest(sdk.UpdateChecksRequest{
		Check: checkReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "monitoring_check", plan.Name.ValueString(), err, httpResp)

		return
	}

	check := result.Check
	if check == nil {
		resp.Diagnostics.AddError("API returned nil", "MonitoringCheck is nil in the response")
		return
	}
	plan.ID = types.Int64Value(*check.Id)
	if check.Name != nil {
		plan.Name = types.StringValue(*check.Name)
	}
	if check.CheckInterval.IsSet() {
		plan.CheckInterval = types.Int64Value(*check.CheckInterval.Get())
	}
	if check.InUptime != nil {
		plan.InUptime = types.BoolValue(*check.InUptime)
	}
	if check.Active != nil {
		plan.Active = types.BoolValue(*check.Active)
	}
	if check.Severity != nil {
		plan.Severity = types.StringValue(*check.Severity)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringCheckResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.ChecksAPI.DeleteChecks(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "monitoring_check", "", err, httpResp)

		return
	}
}

func (r *monitoringCheckResource) ImportState(
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
