package budget

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
	_ resource.Resource                = &budgetResource{}
	_ resource.ResourceWithConfigure   = &budgetResource{}
	_ resource.ResourceWithImportState = &budgetResource{}
)

type budgetResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &budgetResource{}
}

func (r *budgetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_budget"
}

func (r *budgetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BudgetSchema(ctx)
}

func (r *budgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan budgetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddBudgetsRequestBudget{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Year.IsNull() {
		body.Year = plan.Year.ValueInt64Pointer()
	}
	if !plan.Interval.IsNull() {
		body.Interval = plan.Interval.ValueStringPointer()
	}
	if !plan.Scope.IsNull() {
		body.Scope = plan.Scope.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	result, httpResp, err := client.BudgetsAPI.AddBudgets(ctx).AddBudgetsRequest(sdk.AddBudgetsRequest{
		Budget: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "budget", plan.Name.ValueString(), err, httpResp)

		return
	}

	createBudget := result.GetBudget()
	id := (&createBudget).GetId()

	readResult, httpResp, err := client.BudgetsAPI.GetBudgets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "budget", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "budget",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readBudget := readResult.GetBudget()
	mapGetResponseToModel(&plan, &readBudget)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *budgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state budgetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.BudgetsAPI.GetBudgets(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "budget", "", err, httpResp)

		return
	}

	budget := result.GetBudget()
	mapGetResponseToModel(&state, &budget)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *budgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan budgetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateBudgetsRequestBudget{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Year.IsNull() {
		body.Year = plan.Year.ValueInt64Pointer()
	}
	if !plan.Interval.IsNull() {
		body.Interval = plan.Interval.ValueStringPointer()
	}
	if !plan.Scope.IsNull() {
		body.Scope = plan.Scope.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	_, httpResp, err := client.BudgetsAPI.UpdateBudgets(ctx, id).UpdateBudgetsRequest(sdk.UpdateBudgetsRequest{
		Budget: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "budget", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.BudgetsAPI.GetBudgets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "budget", plan.Name.ValueString(), err, httpResp)

		return
	}

	readBudget := readResult.GetBudget()
	mapGetResponseToModel(&plan, &readBudget)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *budgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state budgetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.BudgetsAPI.RemoveBudgets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "budget", "", err, httpResp)

		return
	}
}

func (r *budgetResource) ImportState(
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

func mapAddResponseToModel(model *budgetModel, b *sdk.AddBudgets200ResponseAllOfBudget) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
	if v := b.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if b.Interval != nil {
		model.Interval = types.StringValue(*b.Interval)
	}
	if b.RefScope != nil {
		model.Scope = types.StringValue(*b.RefScope)
	}
	if b.Enabled != nil {
		model.Enabled = types.BoolValue(*b.Enabled)
	}
}

func mapGetResponseToModel(model *budgetModel, b *sdk.GetBudgets200ResponseAllOfBudget) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
	if v := b.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if b.Interval != nil {
		model.Interval = types.StringValue(*b.Interval)
	}
	if b.RefScope != nil {
		model.Scope = types.StringValue(*b.RefScope)
	}
	if b.Enabled != nil {
		model.Enabled = types.BoolValue(*b.Enabled)
	}
	if b.Year != nil {
		if v, err := strconv.ParseInt(*b.Year, 10, 64); err == nil {
			model.Year = types.Int64Value(v)
		}
	}
}

func mapUpdateResponseToModel(model *budgetModel, b *sdk.UpdateBudgets200ResponseAllOfBudget) {
	if b.Id != nil {
		model.ID = types.Int64Value(*b.Id)
	}
	if b.Name != nil {
		model.Name = types.StringValue(*b.Name)
	}
	if v := b.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if b.Interval != nil {
		model.Interval = types.StringValue(*b.Interval)
	}
	if b.RefScope != nil {
		model.Scope = types.StringValue(*b.RefScope)
	}
	if b.Enabled != nil {
		model.Enabled = types.BoolValue(*b.Enabled)
	}
}
