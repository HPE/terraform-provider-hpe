package network_group

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &networkGroupResource{}
	_ resource.ResourceWithConfigure   = &networkGroupResource{}
	_ resource.ResourceWithImportState = &networkGroupResource{}
)

type networkGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkGroupResource{}
}

func (r *networkGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_group"
}

func (r *networkGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkGroupSchema(ctx)
}

func (r *networkGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.CreateNetworkGroupRequestNetworkGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkGroup(ctx).CreateNetworkGroupRequest(sdk.CreateNetworkGroupRequest{
		NetworkGroup: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	// The API returns the id inside {"networkGroup":{"id":...}} but the SDK model
	// only exposes a top-level Id field. Extract from AdditionalProperties.
	var newID int64
	if ngData, ok := result.AdditionalProperties["networkGroup"]; ok {
		if ngMap, ok := ngData.(map[string]interface{}); ok {
			if id, ok := ngMap["id"].(float64); ok {
				newID = int64(id)
			}
		}
	}
	if newID == 0 {
		newID = result.GetId()
	}
	if newID == 0 {
		resp.Diagnostics.AddError("Failed to extract ID", "Could not determine network group ID from create response")
		return
	}

	// Re-read to get full object
	readResult, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, newID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", "", err, httpResp)
		return
	}

	group := readResult.GetNetworkGroup()
	mapResponseToModel(&plan, &group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", "", err, httpResp)
		return
	}

	group := result.GetNetworkGroup()
	mapResponseToModel(&state, &group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateNetworkGroupRequestNetworkGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkGroup(ctx, id).UpdateNetworkGroupRequest(sdk.UpdateNetworkGroupRequest{
		NetworkGroup: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get full object
	readResult, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", "", err, httpResp)
		return
	}

	group := readResult.GetNetworkGroup()
	mapResponseToModel(&plan, &group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteNetworkGroup(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_group", "", err, httpResp)
		return
	}
}

func (r *networkGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapResponseToModel(model *networkGroupModel, group *sdk.GetNetworkGroup200ResponseNetworkGroup) {
	if group.Id != nil {
		model.ID = types.Int64Value(*group.Id)
	}
	if group.Name != nil {
		model.Name = types.StringValue(*group.Name)
	}
	if group.Description != nil {
		model.Description = types.StringValue(*group.Description)
	} else {
		model.Description = types.StringNull()
	}
	if group.Visibility != nil {
		model.Visibility = types.StringValue(*group.Visibility)
	}
	if group.Active != nil {
		model.Active = types.BoolValue(*group.Active)
	}
}
