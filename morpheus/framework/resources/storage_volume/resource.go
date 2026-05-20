package storage_volume

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
	_ resource.Resource                = &storageVolumeResource{}
	_ resource.ResourceWithConfigure   = &storageVolumeResource{}
	_ resource.ResourceWithImportState = &storageVolumeResource{}
)

type storageVolumeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &storageVolumeResource{}
}

func (r *storageVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_storage_volume"
}

func (r *storageVolumeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageVolumeSchema(ctx)
}

func (r *storageVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan storageVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeType := plan.Type.ValueString()
	if volumeType == "" {
		volumeType = "standard"
	}

	storageServerID := plan.StorageServerID.ValueInt64()

	body := sdk.AddStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueString(),
		Type: volumeType,
		StorageServer: sdk.AddStorageVolumesRequestStorageVolumeStorageServer{
			Id: storageServerID,
		},
	}
	if !plan.MaxStorage.IsNull() {
		config := map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64(),
		}
		body.Config = config
	}

	result, httpResp, err := client.StorageAPI.AddStorageVolumes(ctx).AddStorageVolumesRequest(sdk.AddStorageVolumesRequest{
		StorageVolume: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_volume", plan.Name.ValueString(), err, httpResp)
		return
	}

	sv := result.GetStorageVolume()
	mapCreateResponseToModel(&plan, &sv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state storageVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()
	idParam := sdk.Int64AsGetStorageVolumesIdParameter(&id)

	result, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, idParam).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_volume", "", err, httpResp)
		return
	}

	sv := result.GetStorageVolume()
	mapGetResponseToModel(&state, &sv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *storageVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan storageVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()
	idParam := sdk.Int64AsUpdateStorageVolumesIdParameter(&id)

	body := sdk.UpdateStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Type.IsNull() {
		body.Type = plan.Type.ValueStringPointer()
	}
	if !plan.MaxStorage.IsNull() {
		config := map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64(),
		}
		body.Config = config
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageVolumes(ctx, idParam).UpdateStorageVolumesRequest(sdk.UpdateStorageVolumesRequest{
		StorageVolume: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_volume", plan.Name.ValueString(), err, httpResp)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state storageVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()
	idParam := sdk.Int64AsUpdateStorageVolumesIdParameter(&id)

	_, httpResp, err := client.StorageAPI.RemoveStorageVolumes(ctx, idParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_volume", "", err, httpResp)
		return
	}
}

func (r *storageVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *storageVolumeModel, sv *sdk.AddStorageVolumes200ResponseAllOfStorageVolume) {
	if sv.Id != nil {
		model.ID = types.Int64Value(*sv.Id)
	}
	if sv.Name != nil {
		model.Name = types.StringValue(*sv.Name)
	}
	if sv.MaxStorage != nil {
		model.MaxStorage = types.Int64Value(*sv.MaxStorage)
	}
	if sv.Status != nil {
		model.Status = types.StringValue(*sv.Status)
	}
}

func mapGetResponseToModel(model *storageVolumeModel, sv *sdk.GetStorageVolumes200ResponseStorageVolume) {
	if sv.Id != nil {
		model.ID = types.Int64Value(*sv.Id)
	}
	if sv.Name != nil {
		model.Name = types.StringValue(*sv.Name)
	}
	if sv.MaxStorage != nil {
		model.MaxStorage = types.Int64Value(*sv.MaxStorage)
	}
	if sv.Status != nil {
		model.Status = types.StringValue(*sv.Status)
	}
}
