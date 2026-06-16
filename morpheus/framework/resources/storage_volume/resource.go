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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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

func (r *storageVolumeResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
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

	volumeType := strconv.FormatInt(plan.TypeId.ValueInt64(), 10)

	body := sdk.AddStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueString(),
		Type: volumeType,
	}
	if !plan.StorageServerID.IsNull() {
		body.StorageServer = sdk.AddStorageVolumesRequestStorageVolumeStorageServer{
			Id: plan.StorageServerID.ValueInt64(),
		}
	}
	if !plan.MaxStorage.IsNull() {
		config := map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64(),
		}
		body.Config = config
	}

	result, httpResp, err := client.StorageAPI.AddStorageVolumes(ctx).
		AddStorageVolumesRequest(sdk.AddStorageVolumesRequest{
			StorageVolume: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_volume", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.StorageVolume == nil || result.StorageVolume.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "StorageVolume ID is nil in the create response")

		return
	}

	id := *result.StorageVolume.Id
	idParam := sdk.GetStorageVolumesIdParameter{Int64: &id}

	readResult, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, idParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_volume", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "storage_volume",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.StorageVolume == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageVolume is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.StorageVolume)

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
	idParam := sdk.GetStorageVolumesIdParameter{Int64: &id}

	result, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, idParam).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_volume", "", err, httpResp)

		return
	}

	sv := result.StorageVolume
	if sv == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageVolume is nil in the response")

		return
	}
	mapGetResponseToModel(&state, sv)

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
	idParam := sdk.UpdateStorageVolumesIdParameter{Int64: &id}

	body := sdk.UpdateStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.MaxStorage.IsNull() {
		config := map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64(),
		}
		body.Config = config
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageVolumes(ctx, idParam).
		UpdateStorageVolumesRequest(sdk.UpdateStorageVolumesRequest{
			StorageVolume: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_volume", plan.Name.ValueString(), err, httpResp)

		return
	}

	readIdParam := sdk.GetStorageVolumesIdParameter{Int64: &id}

	readResult, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, readIdParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_volume", plan.Name.ValueString(), err, httpResp)

		return
	}

	if readResult.StorageVolume == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageVolume is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.StorageVolume)

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
	idParam := sdk.UpdateStorageVolumesIdParameter{Int64: &id}

	_, httpResp, err := client.StorageAPI.RemoveStorageVolumes(ctx, idParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_volume", "", err, httpResp)

		return
	}
}

func (r *storageVolumeResource) ImportState(
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

func mapGetResponseToModel(model *storageVolumeModel, sv *sdk.GetStorageVolumes200ResponseStorageVolume) {
	if sv.Id != nil {
		model.ID = types.Int64Value(*sv.Id)
	}
	if sv.Name != nil {
		model.Name = types.StringValue(*sv.Name)
	}
	if sv.TypeId != nil {
		model.TypeId = types.Int64Value(*sv.TypeId)
	}
	if storageServer := sv.StorageServer; storageServer != nil {
		if id, ok := storageServer["id"].(float64); ok {
			model.StorageServerID = types.Int64Value(int64(id))
		}
	}
	if sv.MaxStorage != nil {
		model.MaxStorage = types.Int64Value(*sv.MaxStorage)
	}
	if sv.Status != nil {
		model.Status = types.StringValue(*sv.Status)
	}
}
