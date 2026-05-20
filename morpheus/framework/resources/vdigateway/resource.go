package vdigateway

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/polling"
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

func (r *vdiGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
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

	// TODO: Build request body from plan and call API
	// result, httpResp, err := client.VdiAPI.AddVdiGateway(ctx).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "vdi_gateway", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = polling.ForCreate

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

	// TODO: Call API to read resource
	// result, httpResp, err := client.VdiAPI.GetVdiGateway(ctx, id).Execute()
	var httpResp *http.Response
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "vdi_gateway", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	// TODO: Map response to state model
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

	// TODO: Build request body from plan and call API
	// _, httpResp, err := client.VdiAPI.UpdateVdiGateway(ctx, id).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "vdi_gateway", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = id

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

	// TODO: Call API to delete resource
	// httpResp, err := client.VdiAPI.DeleteVdiGateway(ctx, id).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "vdi_gateway", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id
	_ = fmt.Sprintf
	_ = polling.ForDelete
}

func (r *vdiGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
