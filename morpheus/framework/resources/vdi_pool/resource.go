package vdi_pool

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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &vdiPoolResource{}
	_ resource.ResourceWithConfigure   = &vdiPoolResource{}
	_ resource.ResourceWithImportState = &vdiPoolResource{}
)

type vdiPoolResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &vdiPoolResource{}
}

func (r *vdiPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_vdi_pool"
}

func (r *vdiPoolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VdiPoolResourceSchema(ctx)
}

func (r *vdiPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan VdiPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oneOf := sdk.AddVDIPoolsRequestVdiPoolOneOf{
		Name:        plan.Name.ValueString(),
		MaxPoolSize: float32(plan.MaxPoolSize.ValueInt64()),
	}
	if !plan.Description.IsNull() {
		oneOf.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		oneOf.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.PersistentUser.IsNull() {
		oneOf.PersistentUser = plan.PersistentUser.ValueBoolPointer()
	}
	if !plan.Recyclable.IsNull() {
		oneOf.Recyclable = plan.Recyclable.ValueBoolPointer()
	}
	if !plan.AllowCopy.IsNull() {
		oneOf.AllowCopy = plan.AllowCopy.ValueBoolPointer()
	}
	if !plan.AllowPrinter.IsNull() {
		oneOf.AllowPrinter = plan.AllowPrinter.ValueBoolPointer()
	}
	if !plan.AllowFileshare.IsNull() {
		oneOf.AllowFileshare = plan.AllowFileshare.ValueBoolPointer()
	}
	if !plan.AllowHypervisorConsole.IsNull() {
		oneOf.AllowHypervisorConsole = plan.AllowHypervisorConsole.ValueBoolPointer()
	}
	if !plan.AutoCreateLocalUserOnReservation.IsNull() {
		oneOf.AutoCreateLocalUserOnReservation = plan.AutoCreateLocalUserOnReservation.ValueBoolPointer()
	}
	if !plan.MinIdle.IsNull() {
		v := float32(plan.MinIdle.ValueInt64())
		oneOf.MinIdle = &v
	}
	if !plan.InitialPoolSize.IsNull() {
		v := float32(plan.InitialPoolSize.ValueInt64())
		oneOf.InitialPoolSize = &v
	}
	if !plan.MaxIdle.IsNull() {
		v := float32(plan.MaxIdle.ValueInt64())
		oneOf.MaxIdle = &v
	}
	if !plan.IdleTimeout.IsNull() {
		v := float32(plan.IdleTimeout.ValueInt64())
		oneOf.AllocationTimeoutMinutes = &v
	}

	vdiPool := sdk.AddVDIPoolsRequestVdiPool{AddVDIPoolsRequestVdiPoolOneOf: &oneOf}

	_, httpResp, err := client.VDIAPI.AddVDIPools(ctx).AddVDIPoolsRequest(sdk.AddVDIPoolsRequest{
		VdiPool: vdiPool,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "vdi_pool", plan.Name.ValueString(), err, httpResp)

		return
	}

	// The create response is a union type; re-read to get the full object with ID.
	// We need to list and find by name since the response doesn't cleanly expose the ID.
	listResult, httpResp, err := client.VDIAPI.ListVDIPools(ctx).Name(plan.Name.ValueString()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"Read Error After Create",
			"VDI pool was created successfully but could not be read back. "+
				"The resource may exist in Morpheus. Check the Morpheus UI and import manually if needed: "+
				"'terraform import <resource_type>.<name> <id>'",
		)

		return
	}

	pools := listResult.VdiPools
	if len(pools) == 0 {
		resp.Diagnostics.AddError(
			"VDI Pool Not Found",
			"VDI pool was created successfully but could not be found by name. "+
				"The resource may exist in Morpheus. Check the Morpheus UI and import manually if needed: "+
				"'terraform import <resource_type>.<name> <id>'",
		)

		return
	}

	if pools[0].Id == nil {
		resp.Diagnostics.AddError(
			"VDI Pool ID Not Returned",
			"VDI pool was created successfully but the API did not return an ID.",
		)

		return
	}

	poolID := *pools[0].Id

	state, diags := r.getVdiPoolAsState(ctx, poolID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "vdi_pool",
			ResourceID:   poolID,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *vdiPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VdiPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, diags := r.getVdiPoolAsState(ctx, state.Id.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if model == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *vdiPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan VdiPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateVDIPoolsRequestVdiPool{}
	if !plan.Name.IsNull() {
		body.Name = plan.Name.ValueStringPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.PersistentUser.IsNull() {
		body.PersistentUser = plan.PersistentUser.ValueBoolPointer()
	}
	if !plan.Recyclable.IsNull() {
		body.Recyclable = plan.Recyclable.ValueBoolPointer()
	}
	if !plan.AllowCopy.IsNull() {
		body.AllowCopy = plan.AllowCopy.ValueBoolPointer()
	}
	if !plan.AllowPrinter.IsNull() {
		body.AllowPrinter = plan.AllowPrinter.ValueBoolPointer()
	}
	if !plan.AllowFileshare.IsNull() {
		body.AllowFileshare = plan.AllowFileshare.ValueBoolPointer()
	}
	if !plan.AllowHypervisorConsole.IsNull() {
		body.AllowHypervisorConsole = plan.AllowHypervisorConsole.ValueBoolPointer()
	}
	if !plan.AutoCreateLocalUserOnReservation.IsNull() {
		body.AutoCreateLocalUserOnReservation = plan.AutoCreateLocalUserOnReservation.ValueBoolPointer()
	}
	if !plan.MaxPoolSize.IsNull() {
		v := float32(plan.MaxPoolSize.ValueInt64())
		body.MaxPoolSize = &v
	}
	if !plan.MinIdle.IsNull() {
		v := float32(plan.MinIdle.ValueInt64())
		body.MinIdle = &v
	}
	if !plan.MaxIdle.IsNull() {
		v := float32(plan.MaxIdle.ValueInt64())
		body.MaxIdle = &v
	}
	if !plan.InitialPoolSize.IsNull() {
		v := float32(plan.InitialPoolSize.ValueInt64())
		body.InitialPoolSize = &v
	}
	if !plan.IdleTimeout.IsNull() {
		v := float32(plan.IdleTimeout.ValueInt64())
		body.AllocationTimeoutMinutes = &v
	}

	_, httpResp, err := client.VDIAPI.UpdateVDIPools(ctx, id).UpdateVDIPoolsRequest(sdk.UpdateVDIPoolsRequest{
		VdiPool: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "vdi_pool", plan.Name.ValueString(), err, httpResp)

		return
	}

	state, diags := r.getVdiPoolAsState(ctx, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *vdiPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state VdiPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.VDIAPI.RemoveVDIPools(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "vdi_pool", "", err, httpResp)

		return
	}
}

func (r *vdiPoolResource) ImportState(
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

func (r *vdiPoolResource) getVdiPoolAsState(
	ctx context.Context,
	id int64,
) (*VdiPoolModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&diags, err)

		return nil, diags
	}

	result, httpResp, err := client.VDIAPI.GetVDIPools(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return nil, diags
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&diags, errfmt.OpRead, "vdi_pool", "", err, httpResp)

		return nil, diags
	}

	pool := result.VdiPool
	if pool == nil {
		diags.AddError("API returned nil", "VdiPool is nil in the response")

		return nil, diags
	}
	var model VdiPoolModel
	mapGetResponseToModel(&model, pool)

	return &model, diags
}

func mapGetResponseToModel(model *VdiPoolModel, pool *sdk.GetVDIPools200ResponseVdiPool) {
	if pool.Id != nil {
		model.Id = types.Int64Value(*pool.Id)
	}
	if pool.Name != nil {
		model.Name = types.StringValue(*pool.Name)
	}
	if v := pool.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if pool.Enabled != nil {
		model.Enabled = types.BoolValue(*pool.Enabled)
	}
	if pool.AutoCreateLocalUserOnReservation != nil {
		model.AutoCreateLocalUserOnReservation = types.BoolValue(*pool.AutoCreateLocalUserOnReservation)
	} else {
		model.AutoCreateLocalUserOnReservation = types.BoolNull()
	}
	if v := pool.PersistentUser.Get(); v != nil {
		model.PersistentUser = types.BoolValue(*v)
	} else {
		model.PersistentUser = types.BoolNull()
	}
	if v := pool.Recyclable.Get(); v != nil {
		model.Recyclable = types.BoolValue(*v)
	} else {
		model.Recyclable = types.BoolNull()
	}
	if v := pool.AllowHypervisorConsole.Get(); v != nil {
		model.AllowHypervisorConsole = types.BoolValue(*v)
	} else {
		model.AllowHypervisorConsole = types.BoolNull()
	}
	if v := pool.AllowCopy.Get(); v != nil {
		model.AllowCopy = types.BoolValue(*v)
	} else {
		model.AllowCopy = types.BoolNull()
	}
	if v := pool.AllowPrinter.Get(); v != nil {
		model.AllowPrinter = types.BoolValue(*v)
	} else {
		model.AllowPrinter = types.BoolNull()
	}
	if v := pool.AllowFileshare.Get(); v != nil {
		model.AllowFileshare = types.BoolValue(*v)
	} else {
		model.AllowFileshare = types.BoolNull()
	}
	// IdleTimeout maps to AllocationTimeoutMinutes in the API response.
	if pool.AllocationTimeoutMinutes != nil {
		model.IdleTimeout = types.Int64Value(*pool.AllocationTimeoutMinutes)
	} else {
		model.IdleTimeout = types.Int64Null()
	}
	if pool.MaxPoolSize != nil {
		model.MaxPoolSize = types.Int64Value(*pool.MaxPoolSize)
	}
	if pool.MinIdle != nil {
		model.MinIdle = types.Int64Value(*pool.MinIdle)
	} else {
		model.MinIdle = types.Int64Null()
	}
	if pool.InitialPoolSize != nil {
		model.InitialPoolSize = types.Int64Value(*pool.InitialPoolSize)
	} else {
		model.InitialPoolSize = types.Int64Null()
	}
	if pool.MaxIdle != nil {
		model.MaxIdle = types.Int64Value(*pool.MaxIdle)
	} else {
		model.MaxIdle = types.Int64Null()
	}
}
