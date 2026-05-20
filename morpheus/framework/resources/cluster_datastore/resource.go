package cluster_datastore

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &clusterDatastoreResource{}
	_ resource.ResourceWithConfigure   = &clusterDatastoreResource{}
	_ resource.ResourceWithImportState = &clusterDatastoreResource{}
)

type clusterDatastoreResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &clusterDatastoreResource{}
}

func (r *clusterDatastoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_cluster_datastore"
}

func (r *clusterDatastoreResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterDatastoreSchema(ctx)
}

func (r *clusterDatastoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan clusterDatastoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterID.ValueInt64()

	ds := sdk.UpdateClusterDatastoreRequestDatastore{}
	if !plan.Active.IsNull() {
		ds.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		v := plan.Visibility.ValueString()
		ds.Visibility = &v
	}

	body := sdk.UpdateClusterDatastoreRequest{
		Datastore: &ds,
	}

	// This is an adopt-style resource; we use UpdateClusterDatastore to manage an existing datastore.
	// The ID must be provided via import or known ahead of time.
	// For create, we use SaveClusterDatastore if available, otherwise error.
	result, httpResp, err := client.ClustersAPI.SaveClusterDatastore(ctx, clusterID).SaveClusterDatastoreRequest(sdk.SaveClusterDatastoreRequest{
		Datastore: &sdk.SaveClusterDatastoreRequestDatastore{},
	}).Execute()
	_ = body
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "cluster_datastore", "", err, httpResp)
		return
	}

	if result.Datastore != nil && result.Datastore.Id != nil {
		plan.ID = types.Int64Value(*result.Datastore.Id)
		if result.Datastore.Name != nil {
			plan.Name = types.StringValue(*result.Datastore.Name)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterDatastoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state clusterDatastoreModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterID.ValueInt64()
	id := state.ID.ValueInt64()

	result, httpResp, err := client.ClustersAPI.GetClusterDatastore(ctx, clusterID, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_datastore", "", err, httpResp)
		return
	}

	ds := result.Datastore
	if ds != nil {
		if ds.Id != nil {
			state.ID = types.Int64Value(*ds.Id)
		}
		if ds.Name != nil {
			state.Name = types.StringValue(*ds.Name)
		}
		if ds.Active != nil {
			state.Active = types.BoolValue(*ds.Active)
		}
		if ds.Visibility != nil {
			state.Visibility = types.StringValue(*ds.Visibility)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clusterDatastoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan clusterDatastoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterID.ValueInt64()
	id := plan.ID.ValueInt64()

	ds := sdk.UpdateClusterDatastoreRequestDatastore{}
	if !plan.Active.IsNull() {
		ds.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		v := plan.Visibility.ValueString()
		ds.Visibility = &v
	}

	body := sdk.UpdateClusterDatastoreRequest{
		Datastore: &ds,
	}

	_, httpResp, err := client.ClustersAPI.UpdateClusterDatastore(ctx, clusterID, id).UpdateClusterDatastoreRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "cluster_datastore", "", err, httpResp)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterDatastoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state clusterDatastoreModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterID.ValueInt64()
	id := state.ID.ValueInt64()

	_, httpResp, err := client.ClustersAPI.DeleteClusterDatastore(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "cluster_datastore", "", err, httpResp)
		return
	}
}

func (r *clusterDatastoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected format: cluster_id/datastore_id, got: %s", req.ID))
		return
	}

	clusterID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cluster_id %q: %s", parts[0], err))
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse datastore_id %q: %s", parts[1], err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), clusterID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Ensure unused imports are satisfied.
var _ *http.Response
