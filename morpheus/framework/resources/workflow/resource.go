package workflow

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &workflowResource{}
	_ resource.ResourceWithConfigure   = &workflowResource{}
	_ resource.ResourceWithImportState = &workflowResource{}
)

type workflowResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &workflowResource{}
}

func (r *workflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_workflow"
}

func (r *workflowResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = WorkflowSchema(ctx)
}

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan workflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddWorkflowsRequestTaskSet{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Type.IsNull() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Labels.IsNull() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}

	result, httpResp, err := client.AutomationAPI.AddWorkflows(ctx).AddWorkflowsRequest(sdk.AddWorkflowsRequest{
		TaskSet: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "workflow", plan.Name.ValueString(), err, httpResp)
		return
	}

	ts := result.GetTaskSet()
	mapAddResponseToModel(ctx, &plan, &ts, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state workflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.AutomationAPI.GetWorkflows(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "workflow", "", err, httpResp)
		return
	}

	ts := result.GetTaskSet()
	mapGetResponseToModel(ctx, &state, &ts, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan workflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateWorkflowsRequestTaskSet{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Type.IsNull() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.Labels.IsNull() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}

	_, httpResp, err := client.AutomationAPI.UpdateWorkflows(ctx, id).UpdateWorkflowsRequest(sdk.UpdateWorkflowsRequest{
		TaskSet: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "workflow", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get current state
	getResult, httpResp, err := client.AutomationAPI.GetWorkflows(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "workflow", "", err, httpResp)
		return
	}

	ts := getResult.GetTaskSet()
	mapGetResponseToModel(ctx, &plan, &ts, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state workflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.AutomationAPI.RemoveWorkflows(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "workflow", "", err, httpResp)
		return
	}
}

func (r *workflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(ctx context.Context, model *workflowModel, ts *sdk.AddWorkflows200ResponseAllOfTaskSet, diags *diag.Diagnostics) {
	if ts.Id != nil {
		model.ID = types.Int64Value(*ts.Id)
	}
	if ts.Name != nil {
		model.Name = types.StringValue(*ts.Name)
	}
	if ts.Description.IsSet() && ts.Description.Get() != nil {
		model.Description = types.StringValue(*ts.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if ts.Type != nil {
		model.Type = types.StringValue(*ts.Type)
	}
	if ts.Platform.IsSet() && ts.Platform.Get() != nil {
		model.Platform = types.StringValue(*ts.Platform.Get())
	} else {
		model.Platform = types.StringNull()
	}
	if ts.Visibility != nil {
		model.Visibility = types.StringValue(*ts.Visibility)
	}
	if len(ts.Labels) > 0 {
		labelList, d := types.ListValueFrom(ctx, types.StringType, ts.Labels)
		diags.Append(d...)
		model.Labels = labelList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
}

func mapGetResponseToModel(ctx context.Context, model *workflowModel, ts *sdk.GetWorkflows200ResponseAllOfTaskSet, diags *diag.Diagnostics) {
	if ts.Id != nil {
		model.ID = types.Int64Value(*ts.Id)
	}
	if ts.Name != nil {
		model.Name = types.StringValue(*ts.Name)
	}
	if ts.Description.IsSet() && ts.Description.Get() != nil {
		model.Description = types.StringValue(*ts.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if ts.Type != nil {
		model.Type = types.StringValue(*ts.Type)
	}
	if ts.Platform.IsSet() && ts.Platform.Get() != nil {
		model.Platform = types.StringValue(*ts.Platform.Get())
	} else {
		model.Platform = types.StringNull()
	}
	if ts.Visibility != nil {
		model.Visibility = types.StringValue(*ts.Visibility)
	}
	if len(ts.Labels) > 0 {
		labelList, d := types.ListValueFrom(ctx, types.StringType, ts.Labels)
		diags.Append(d...)
		model.Labels = labelList
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
}
