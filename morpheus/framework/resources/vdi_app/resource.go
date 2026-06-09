package vdi_app

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
	_ resource.Resource                = &vdiAppResource{}
	_ resource.ResourceWithConfigure   = &vdiAppResource{}
	_ resource.ResourceWithImportState = &vdiAppResource{}
)

type vdiAppResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &vdiAppResource{}
}

func (r *vdiAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_vdi_app"
}

func (r *vdiAppResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VdiAppSchema(ctx)
}

func (r *vdiAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan vdiAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddVDIAppsRequestVdiAppOneOf{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.LaunchPrefix.IsNull() {
		body.LaunchPrefix = plan.LaunchPrefix.ValueStringPointer()
	}

	result, httpResp, err := client.VDIAPI.AddVDIApps(ctx).AddVDIAppsRequest(sdk.AddVDIAppsRequest{
		VdiApp: sdk.AddVDIAppsRequestVdiApp{AddVDIAppsRequestVdiAppOneOf: &body},
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "vdi_app", plan.Name.ValueString(), err, httpResp)

		return
	}

	app := result.AddVDIApps200ResponseAnyOf.VdiApp
	if app == nil {
		resp.Diagnostics.AddError("API returned nil", "VdiApp is nil in the response")

		return
	}
	mapCreateResponseToModel(&plan, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vdiAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state vdiAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.VDIAPI.GetVDIApps(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "vdi_app", "", err, httpResp)

		return
	}

	app := result.VdiApp
	if app == nil {
		resp.Diagnostics.AddError("API returned nil", "VdiApp is nil in the response")

		return
	}
	mapGetResponseToModel(&state, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vdiAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan vdiAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateVDIAppsRequestVdiAppOneOf{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.LaunchPrefix.IsNull() {
		body.LaunchPrefix = plan.LaunchPrefix.ValueStringPointer()
	}

	result, httpResp, err := client.VDIAPI.UpdateVDIApps(ctx, id).UpdateVDIAppsRequest(sdk.UpdateVDIAppsRequest{
		VdiApp: sdk.UpdateVDIAppsRequestVdiApp{UpdateVDIAppsRequestVdiAppOneOf: &body},
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "vdi_app", plan.Name.ValueString(), err, httpResp)

		return
	}

	app := result.UpdateVDIApps200ResponseAnyOf.VdiApp
	if app == nil {
		resp.Diagnostics.AddError("API returned nil", "VdiApp is nil in the response")

		return
	}
	mapUpdateResponseToModel(&plan, app)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vdiAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state vdiAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.VDIAPI.RemoveVDIApps(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "vdi_app", "", err, httpResp)

		return
	}
}

func (r *vdiAppResource) ImportState(
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

func mapCreateResponseToModel(model *vdiAppModel, app *sdk.AddVDIApps200ResponseAnyOfVdiApp) {
	if app.Id != nil {
		model.ID = types.Int64Value(*app.Id)
	}
	if app.Name != nil {
		model.Name = types.StringValue(*app.Name)
	}
	if v := app.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if app.LaunchPrefix != nil {
		model.LaunchPrefix = types.StringValue(*app.LaunchPrefix)
	} else {
		model.LaunchPrefix = types.StringNull()
	}
}

func mapGetResponseToModel(model *vdiAppModel, app *sdk.GetVDIApps200ResponseVdiApp) {
	if app.Id != nil {
		model.ID = types.Int64Value(*app.Id)
	}
	if app.Name != nil {
		model.Name = types.StringValue(*app.Name)
	}
	if v := app.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if app.LaunchPrefix != nil {
		model.LaunchPrefix = types.StringValue(*app.LaunchPrefix)
	} else {
		model.LaunchPrefix = types.StringNull()
	}
}

func mapUpdateResponseToModel(model *vdiAppModel, app *sdk.UpdateVDIApps200ResponseAnyOfVdiApp) {
	if app.Id != nil {
		model.ID = types.Int64Value(*app.Id)
	}
	if app.Name != nil {
		model.Name = types.StringValue(*app.Name)
	}
	if v := app.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if app.LaunchPrefix != nil {
		model.LaunchPrefix = types.StringValue(*app.LaunchPrefix)
	} else {
		model.LaunchPrefix = types.StringNull()
	}
}
