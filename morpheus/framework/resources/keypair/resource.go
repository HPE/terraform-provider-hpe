package keypair

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
	_ resource.Resource                = &keyPairResource{}
	_ resource.ResourceWithConfigure   = &keyPairResource{}
	_ resource.ResourceWithImportState = &keyPairResource{}
)

type keyPairResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &keyPairResource{}
}

func (r *keyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_key_pair"
}

func (r *keyPairResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = KeyPairSchema(ctx)
}

func (r *keyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan keyPairModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Build request body from plan and call API
	// result, httpResp, err := client.KeyPairsAPI.AddKeyPair(ctx).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "key_pair", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = polling.ForCreate

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *keyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state keyPairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to read resource
	// result, httpResp, err := client.KeyPairsAPI.GetKeyPair(ctx, id).Execute()
	var httpResp *http.Response
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "key_pair", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	// TODO: Map response to state model
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *keyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan keyPairModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	// TODO: Build request body from plan and call API
	// _, httpResp, err := client.KeyPairsAPI.UpdateKeyPair(ctx, id).Body(body).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "key_pair", plan.Name.ValueString(), err, httpResp)
	//     return
	// }
	_ = client
	_ = id

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *keyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state keyPairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	// TODO: Call API to delete resource
	// httpResp, err := client.KeyPairsAPI.DeleteKeyPair(ctx, id).Execute()
	// if err := errfmt.CheckResponse(err, httpResp); err != nil {
	//     errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "key_pair", "", err, httpResp)
	//     return
	// }
	_ = client
	_ = id
	_ = fmt.Sprintf
	_ = polling.ForDelete
}

func (r *keyPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
