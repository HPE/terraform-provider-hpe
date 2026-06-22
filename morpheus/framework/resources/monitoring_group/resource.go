package monitoring_group

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
	_ resource.Resource                = &monitoringGroupResource{}
	_ resource.ResourceWithConfigure   = &monitoringGroupResource{}
	_ resource.ResourceWithImportState = &monitoringGroupResource{}
)

type monitoringGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &monitoringGroupResource{}
}

func (r *monitoringGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_monitoring_group"
}

func (r *monitoringGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = MonitoringGroupResourceSchema(ctx)
}

func (r *monitoringGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan MonitoringGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddCheckGroupsRequestCheckGroup{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.MinHappy.IsNull() {
		mh := int32(plan.MinHappy.ValueInt64()) //nolint:gosec // value range is safe
		body.MinHappy = &mh
	}
	if !plan.Severity.IsNull() {
		body.Severity = plan.Severity.ValueStringPointer()
	}
	if !plan.InUptime.IsNull() {
		body.InUptime = plan.InUptime.ValueBoolPointer()
	}

	result, httpResp, err := client.ChecksAPI.AddCheckGroups(ctx).AddCheckGroupsRequest(sdk.AddCheckGroupsRequest{
		CheckGroup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.CheckGroup == nil || result.CheckGroup.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "CheckGroup ID is nil in the create response")

		return
	}

	id := *result.CheckGroup.Id

	readResult, httpResp, err := client.ChecksAPI.GetCheckGroups(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_group", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "monitoring_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.CheckGroup == nil {
		resp.Diagnostics.AddError("API returned nil", "CheckGroup is nil in the response")

		return
	}
	mapGetGroupResponseToModel(&plan, readResult.CheckGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state MonitoringGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.ChecksAPI.GetCheckGroups(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_group", "", err, httpResp)

		return
	}

	group := result.CheckGroup
	if group == nil {
		resp.Diagnostics.AddError("API returned nil", "CheckGroup is nil in the response")

		return
	}
	mapGetGroupResponseToModel(&state, group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *monitoringGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan MonitoringGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateCheckGroupsRequestCheckGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.MinHappy.IsNull() {
		mh := int32(plan.MinHappy.ValueInt64()) //nolint:gosec // value range is safe
		body.MinHappy = &mh
	}
	if !plan.Severity.IsNull() {
		body.Severity = plan.Severity.ValueStringPointer()
	}
	if !plan.InUptime.IsNull() {
		body.InUptime = plan.InUptime.ValueBoolPointer()
	}

	_, httpResp, err := client.ChecksAPI.UpdateCheckGroups(ctx, id).
		UpdateCheckGroupsRequest(sdk.UpdateCheckGroupsRequest{
			CheckGroup: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "monitoring_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.ChecksAPI.GetCheckGroups(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	if readResult.CheckGroup == nil {
		resp.Diagnostics.AddError("API returned nil", "CheckGroup is nil in the response")

		return
	}
	mapGetGroupResponseToModel(&plan, readResult.CheckGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state MonitoringGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.ChecksAPI.DeleteCheckGroups(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "monitoring_group", "", err, httpResp)

		return
	}
}

func (r *monitoringGroupResource) ImportState(
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

func mapGetGroupResponseToModel(model *MonitoringGroupModel, group *sdk.GetCheckGroups200ResponseCheckGroup) {
	if group.Id != nil {
		model.Id = types.Int64Value(*group.Id)
	}
	if group.Name != nil {
		model.Name = types.StringValue(*group.Name)
	}
	if group.Description.IsSet() && group.Description.Get() != nil {
		model.Description = types.StringValue(*group.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if group.MinHappy != nil {
		model.MinHappy = types.Int64Value(*group.MinHappy)
	} else {
		model.MinHappy = types.Int64Null()
	}
	if group.Severity != nil {
		model.Severity = types.StringValue(*group.Severity)
	} else {
		model.Severity = types.StringNull()
	}
	if group.InUptime != nil {
		model.InUptime = types.BoolValue(*group.InUptime)
	}
}
