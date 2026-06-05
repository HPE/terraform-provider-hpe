package vdi_gateway

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
	_ resource.Resource                = &vdiGatewayResource{}
	_ resource.ResourceWithConfigure   = &vdiGatewayResource{}
	_ resource.ResourceWithImportState = &vdiGatewayResource{}
)

type vdiGatewayResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &vdiGatewayResource{}
}

func (r *vdiGatewayResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_vdi_gateway"
}

func (r *vdiGatewayResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VdiGatewaySchema(ctx)
}

func (r *vdiGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan vdiGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddVDIGatewaysRequestVdiGatewayOneOf{
		Name:       plan.Name.ValueString(),
		GatewayUrl: plan.GatewayUrl.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	result, httpResp, err := client.VDIAPI.AddVDIGateways(ctx).AddVDIGatewaysRequest(sdk.AddVDIGatewaysRequest{
		VdiGateway: sdk.AddVDIGatewaysRequestVdiGatewayOneOfAsAddVDIGatewaysRequestVdiGateway(&body),
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "vdi_gateway", plan.Name.ValueString(), err, httpResp)

		return
	}

	createGw := result.AddVDIGateways200ResponseAnyOf.GetVdiGateway()
	id := (&createGw).GetId()

	readResult, httpResp, err := client.VDIAPI.GetVDIGateways(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "vdi_gateway", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "vdi_gateway",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readGw := readResult.GetVdiGateway()
	mapGetResponseToModel(&plan, &readGw)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vdiGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state vdiGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.VDIAPI.GetVDIGateways(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "vdi_gateway", "", err, httpResp)

		return
	}

	gw := result.GetVdiGateway()
	mapGetResponseToModel(&state, &gw)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vdiGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan vdiGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateVDIGatewaysRequestVdiGatewayOneOf{
		Name:       plan.Name.ValueStringPointer(),
		GatewayUrl: plan.GatewayUrl.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	_, httpResp, err := client.VDIAPI.UpdateVDIGateways(ctx, id).
		UpdateVDIGatewaysRequest(sdk.UpdateVDIGatewaysRequest{
			VdiGateway: sdk.UpdateVDIGatewaysRequestVdiGatewayOneOfAsUpdateVDIGatewaysRequestVdiGateway(&body),
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "vdi_gateway", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.VDIAPI.GetVDIGateways(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "vdi_gateway", plan.Name.ValueString(), err, httpResp)

		return
	}

	readGw := readResult.GetVdiGateway()
	mapGetResponseToModel(&plan, &readGw)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vdiGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state vdiGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.VDIAPI.RemoveVDIGateways(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "vdi_gateway", "", err, httpResp)

		return
	}
}

func (r *vdiGatewayResource) ImportState(
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

func mapCreateResponseToModel(model *vdiGatewayModel, gw *sdk.AddVDIGateways200ResponseAnyOfVdiGateway) {
	if gw.Id != nil {
		model.ID = types.Int64Value(*gw.Id)
	}
	if gw.Name != nil {
		model.Name = types.StringValue(*gw.Name)
	}
	if v := gw.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if v := gw.GatewayUrl.Get(); v != nil {
		model.GatewayUrl = types.StringValue(*v)
	}
}

func mapGetResponseToModel(model *vdiGatewayModel, gw *sdk.GetVDIGateways200ResponseVdiGateway) {
	if gw.Id != nil {
		model.ID = types.Int64Value(*gw.Id)
	}
	if gw.Name != nil {
		model.Name = types.StringValue(*gw.Name)
	}
	if v := gw.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if v := gw.GatewayUrl.Get(); v != nil {
		model.GatewayUrl = types.StringValue(*v)
	}
}

func mapUpdateResponseToModel(model *vdiGatewayModel, gw *sdk.UpdateVDIGateways200ResponseAnyOfVdiGateway) {
	if gw.Id != nil {
		model.ID = types.Int64Value(*gw.Id)
	}
	if gw.Name != nil {
		model.Name = types.StringValue(*gw.Name)
	}
	if v := gw.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if v := gw.GatewayUrl.Get(); v != nil {
		model.GatewayUrl = types.StringValue(*v)
	}
}
