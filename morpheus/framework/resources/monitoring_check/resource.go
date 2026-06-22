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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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
	resp.Schema = MonitoringCheckResourceSchema(ctx)
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

	var plan MonitoringCheckModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkBody := sdk.WebCheck{
		Name: plan.Name.ValueString(),
	}
	if !plan.CheckTypeId.IsNull() && !plan.CheckTypeId.IsUnknown() {
		// Look up the check type code from the ID since the API requires code
		ctResult, ctHTTPResp, ctErr := client.ChecksAPI.GetCheckTypes(ctx, plan.CheckTypeId.ValueInt64()).Execute()
		if ctErr != nil || (ctHTTPResp != nil && ctHTTPResp.StatusCode >= 400) {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_check", plan.Name.ValueString(),
				fmt.Errorf("failed to look up check type %d", plan.CheckTypeId.ValueInt64()), ctHTTPResp)

			return
		}
		ct := ctResult.CheckType
		if ct == nil {
			resp.Diagnostics.AddError("API returned nil", "CheckType is nil in the response")

			return
		}

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

	if result.Check == nil || result.Check.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "Check ID is nil in the create response")

		return
	}

	id := *result.Check.Id

	readResult, httpResp, err := client.ChecksAPI.GetChecks(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_check", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "monitoring_check",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	check := readResult.Check
	if check == nil {
		resp.Diagnostics.AddError("API returned nil", "MonitoringCheck is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, check)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state MonitoringCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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
	mapGetResponseToModel(&state, check)

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

	var plan MonitoringCheckModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

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

	_, httpResp, err := client.ChecksAPI.UpdateChecks(ctx, id).UpdateChecksRequest(sdk.UpdateChecksRequest{
		Check: checkReq,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "monitoring_check", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.ChecksAPI.GetChecks(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_check", plan.Name.ValueString(), err, httpResp)

		return
	}

	check := readResult.Check
	if check == nil {
		resp.Diagnostics.AddError("API returned nil", "MonitoringCheck is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, check)

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

	var state MonitoringCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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

func mapGetResponseToModel(model *MonitoringCheckModel, check *sdk.GetChecks200ResponseCheck) {
	if check.Id != nil {
		model.Id = types.Int64Value(*check.Id)
	}
	if check.Name != nil {
		model.Name = types.StringValue(*check.Name)
	}
	if check.Description.IsSet() && check.Description.Get() != nil {
		model.Description = types.StringValue(*check.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if check.CheckInterval.IsSet() {
		model.CheckInterval = types.Int64Value(*check.CheckInterval.Get())
	}
	if check.InUptime != nil {
		model.InUptime = types.BoolValue(*check.InUptime)
	}
	if check.Active != nil {
		model.Active = types.BoolValue(*check.Active)
	}
	if check.Severity != nil {
		model.Severity = types.StringValue(*check.Severity)
	}
	if check.CheckType != nil && check.CheckType.Id != nil {
		model.CheckTypeId = types.Int64Value(*check.CheckType.Id)
	}
}
