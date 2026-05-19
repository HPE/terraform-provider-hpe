package deployment

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
	_ resource.Resource                = &deploymentResource{}
	_ resource.ResourceWithConfigure   = &deploymentResource{}
	_ resource.ResourceWithImportState = &deploymentResource{}
)

type deploymentResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &deploymentResource{}
}

func (r *deploymentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_deployment"
}

func (r *deploymentResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DeploymentSchema(ctx)
}

func (r *deploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan deploymentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddDeploymentsRequestDeployment{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	result, httpResp, err := client.DeploymentsAPI.AddDeployments(ctx).AddDeploymentsRequest(sdk.AddDeploymentsRequest{
		Deployment: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "deployment", plan.Name.ValueString(), err, httpResp)
		return
	}

	dep := result.GetDeployment()
	mapAddResponseToModel(&plan, &dep)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state deploymentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.DeploymentsAPI.GetDeployment(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "deployment", "", err, httpResp)
		return
	}

	dep := result.GetDeployment()
	mapGetResponseToModel(&state, &dep)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan deploymentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateDeploymentRequestDeployment{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	result, httpResp, err := client.DeploymentsAPI.UpdateDeployment(ctx, id).UpdateDeploymentRequest(sdk.UpdateDeploymentRequest{
		Deployment: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "deployment", plan.Name.ValueString(), err, httpResp)
		return
	}

	dep := result.GetDeployment()
	mapUpdateResponseToModel(&plan, &dep)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *deploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state deploymentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.DeploymentsAPI.DeleteDeployment(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "deployment", "", err, httpResp)
		return
	}
}

func (r *deploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *deploymentModel, dep *sdk.AddDeployments200ResponseAllOfDeployment) {
	if dep.Id != nil {
		model.ID = types.Int64Value(*dep.Id)
	}
	if dep.Name != nil {
		model.Name = types.StringValue(*dep.Name)
	}
	if dep.Description != nil {
		model.Description = types.StringValue(*dep.Description)
	} else {
		model.Description = types.StringNull()
	}
}

func mapGetResponseToModel(model *deploymentModel, dep *sdk.GetDeployment200ResponseDeployment) {
	if dep.Id != nil {
		model.ID = types.Int64Value(*dep.Id)
	}
	if dep.Name != nil {
		model.Name = types.StringValue(*dep.Name)
	}
	if dep.Description != nil {
		model.Description = types.StringValue(*dep.Description)
	} else {
		model.Description = types.StringNull()
	}
}

func mapUpdateResponseToModel(model *deploymentModel, dep *sdk.UpdateDeployment200ResponseAllOfDeployment) {
	if dep.Id != nil {
		model.ID = types.Int64Value(*dep.Id)
	}
	if dep.Name != nil {
		model.Name = types.StringValue(*dep.Name)
	}
	if dep.Description != nil {
		model.Description = types.StringValue(*dep.Description)
	} else {
		model.Description = types.StringNull()
	}
}
