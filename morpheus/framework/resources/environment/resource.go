package environment

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
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

type environmentResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &environmentResource{}
}

func (r *environmentResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_environment"
}

func (r *environmentResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = EnvironmentSchema(ctx)
}

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddEnvironmentsRequestEnvironment{
		Name: plan.Name.ValueString(),
		Code: plan.Code.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.SortOrder.IsNull() {
		body.SortOrder = plan.SortOrder.ValueInt64Pointer()
	}

	result, httpResp, err := client.EnvironmentsAPI.AddEnvironments(ctx).AddEnvironmentsRequest(sdk.AddEnvironmentsRequest{
		Environment: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "environment", plan.Name.ValueString(), err, httpResp)

		return
	}

	env := result.GetEnvironment()
	mapAddResponseToModel(&plan, &env)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.EnvironmentsAPI.GetEnvironments(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "environment", "", err, httpResp)

		return
	}

	env := result.GetEnvironment()
	mapGetResponseToModel(&state, &env)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateEnvironmentsRequestEnvironment{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.SortOrder.IsNull() {
		body.SortOrder = plan.SortOrder.ValueInt64Pointer()
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}

	result, httpResp, err := client.EnvironmentsAPI.UpdateEnvironments(ctx, id).
		UpdateEnvironmentsRequest(sdk.UpdateEnvironmentsRequest{
			Environment: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "environment", plan.Name.ValueString(), err, httpResp)

		return
	}

	env := result.GetEnvironment()
	mapUpdateResponseToModel(&plan, &env)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.EnvironmentsAPI.DeleteEnvironments(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "environment", "", err, httpResp)

		return
	}
}

func (r *environmentResource) ImportState(
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

func mapAddResponseToModel(model *environmentModel, env *sdk.AddEnvironments200ResponseAllOfEnvironment) {
	if env.Id != nil {
		model.ID = types.Int64Value(*env.Id)
	}
	if env.Name != nil {
		model.Name = types.StringValue(*env.Name)
	}
	if env.Code != nil {
		model.Code = types.StringValue(*env.Code)
	}
	if env.Description != nil {
		model.Description = types.StringValue(*env.Description)
	} else {
		model.Description = types.StringNull()
	}
	if env.Visibility != nil {
		model.Visibility = types.StringValue(*env.Visibility)
	}
	if env.Active != nil {
		model.Active = types.BoolValue(*env.Active)
	}
	if env.SortOrder != nil {
		model.SortOrder = types.Int64Value(*env.SortOrder)
	} else {
		model.SortOrder = types.Int64Null()
	}
}

func mapGetResponseToModel(model *environmentModel, env *sdk.GetEnvironments200ResponseEnvironment) {
	if env.Id != nil {
		model.ID = types.Int64Value(*env.Id)
	}
	if env.Name != nil {
		model.Name = types.StringValue(*env.Name)
	}
	if env.Code != nil {
		model.Code = types.StringValue(*env.Code)
	}
	if env.Description != nil {
		model.Description = types.StringValue(*env.Description)
	} else {
		model.Description = types.StringNull()
	}
	if env.Visibility != nil {
		model.Visibility = types.StringValue(*env.Visibility)
	}
	if env.Active != nil {
		model.Active = types.BoolValue(*env.Active)
	}
	if env.SortOrder != nil {
		model.SortOrder = types.Int64Value(*env.SortOrder)
	} else {
		model.SortOrder = types.Int64Null()
	}
}

func mapUpdateResponseToModel(model *environmentModel, env *sdk.UpdateEnvironments200ResponseAllOfEnvironment) {
	if env.Id != nil {
		model.ID = types.Int64Value(*env.Id)
	}
	if env.Name != nil {
		model.Name = types.StringValue(*env.Name)
	}
	if env.Code != nil {
		model.Code = types.StringValue(*env.Code)
	}
	if env.Description != nil {
		model.Description = types.StringValue(*env.Description)
	} else {
		model.Description = types.StringNull()
	}
	if env.Visibility != nil {
		model.Visibility = types.StringValue(*env.Visibility)
	}
	if env.Active != nil {
		model.Active = types.BoolValue(*env.Active)
	}
	if env.SortOrder != nil {
		model.SortOrder = types.Int64Value(*env.SortOrder)
	} else {
		model.SortOrder = types.Int64Null()
	}
}
