package library_instance_type

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &libraryInstanceTypeResource{}
	_ resource.ResourceWithConfigure   = &libraryInstanceTypeResource{}
	_ resource.ResourceWithImportState = &libraryInstanceTypeResource{}
)

type libraryInstanceTypeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &libraryInstanceTypeResource{}
}

func (r *libraryInstanceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_library_instance_type"
}

func (r *libraryInstanceTypeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LibraryInstanceTypeSchema(ctx)
}

func (r *libraryInstanceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryInstanceTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddInstanceTypeRequestInstanceType{
		Name: plan.Name.ValueString(),
	}
	if !plan.Code.IsNull() {
		body.Code = plan.Code.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Featured.IsNull() {
		body.Featured = plan.Featured.ValueBoolPointer()
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}

	_, httpResp, err := client.LibraryAPI.AddInstanceType(ctx).AddInstanceTypeRequest(sdk.AddInstanceTypeRequest{
		InstanceType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "library_instance_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back to get the created resource
	// The create response only returns success, so we need to list and find by name
	// For now, use the list endpoint to find the newly created instance type
	listResult, httpResp, err := client.LibraryAPI.ListInstanceTypes(ctx).Name(plan.Name.ValueString()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_instance_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	instanceTypes := listResult.GetInstanceTypes()
	if len(instanceTypes) == 0 {
		resp.Diagnostics.AddError("Not Found", "Instance type not found after creation")
		return
	}

	it := instanceTypes[0]
	mapListInstanceTypeToModel(ctx, &plan, &it, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryInstanceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryInstanceTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.LibraryAPI.GetInstanceType(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_instance_type", "", err, httpResp)
		return
	}

	it := result.GetInstanceType()
	mapGetInstanceTypeToModel(ctx, &state, &it, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *libraryInstanceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan libraryInstanceTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateInstanceTypeRequestInstanceType{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Code.IsNull() {
		body.Code = plan.Code.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Category.IsNull() {
		body.Category = plan.Category.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Featured.IsNull() {
		body.Featured = plan.Featured.ValueBoolPointer()
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}

	_, httpResp, err := client.LibraryAPI.UpdateInstanceType(ctx, id).UpdateInstanceTypeRequest(sdk.UpdateInstanceTypeRequest{
		InstanceType: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "library_instance_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back
	result, httpResp, err := client.LibraryAPI.GetInstanceType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "library_instance_type", "", err, httpResp)
		return
	}

	it := result.GetInstanceType()
	mapGetInstanceTypeToModel(ctx, &plan, &it, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *libraryInstanceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state libraryInstanceTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.LibraryAPI.DeleteInstanceType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "library_instance_type", "", err, httpResp)
		return
	}
}

func (r *libraryInstanceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapListInstanceTypeToModel(ctx context.Context, model *libraryInstanceTypeModel, it *sdk.ListInstanceTypes200ResponseAllOfInstanceTypesInner, diags *diag.Diagnostics) {
	if it.Id != nil {
		model.ID = types.Int64Value(*it.Id)
	}
	if it.Name != nil {
		model.Name = types.StringValue(*it.Name)
	}
	if it.Code != nil {
		model.Code = types.StringValue(*it.Code)
	} else {
		model.Code = types.StringNull()
	}
	if v := it.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if it.Category != nil {
		model.Category = types.StringValue(*it.Category)
	} else {
		model.Category = types.StringNull()
	}
	if it.Visibility != nil {
		model.Visibility = types.StringValue(*it.Visibility)
	}
	if it.Featured != nil {
		model.Featured = types.BoolValue(*it.Featured)
	}
	if it.Labels != nil {
		labelValues, d := types.ListValueFrom(ctx, types.StringType, it.Labels)
		diags.Append(d...)
		model.Labels = labelValues
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
}

func mapGetInstanceTypeToModel(ctx context.Context, model *libraryInstanceTypeModel, it *sdk.GetInstanceType200ResponseInstanceType, diags *diag.Diagnostics) {
	if it.Id != nil {
		model.ID = types.Int64Value(*it.Id)
	}
	if it.Name != nil {
		model.Name = types.StringValue(*it.Name)
	}
	if it.Code != nil {
		model.Code = types.StringValue(*it.Code)
	} else {
		model.Code = types.StringNull()
	}
	if v := it.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if it.Category != nil {
		model.Category = types.StringValue(*it.Category)
	} else {
		model.Category = types.StringNull()
	}
	if it.Visibility != nil {
		model.Visibility = types.StringValue(*it.Visibility)
	}
	if it.Featured != nil {
		model.Featured = types.BoolValue(*it.Featured)
	}
	if it.Labels != nil {
		labelValues, d := types.ListValueFrom(ctx, types.StringType, it.Labels)
		diags.Append(d...)
		model.Labels = labelValues
	} else {
		model.Labels = types.ListNull(types.StringType)
	}
}
