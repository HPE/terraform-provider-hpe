// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

func (r *storageVolumeResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = StorageVolumeResourceSchema(ctx)
}

func (r *storageVolumeResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// config (config, config_alletramp_bmaas) is write-only: its values are present
	// in the request config but null in the plan/state. Read them from req.Config.
	var config StorageVolumeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API "type" field accepts either a storage volume type code or id.
	var volumeType string
	if !plan.TypeCode.IsNull() && !plan.TypeCode.IsUnknown() {
		volumeType = plan.TypeCode.ValueString()
	} else {
		volumeType = strconv.FormatInt(plan.TypeId.ValueInt64(), 10)
	}

	body := sdk.AddStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueString(),
		Type: volumeType,
	}

	// max_storage is expressed in GiB. The Morpheus API stores and interprets
	// storageVolume.maxStorage in bytes (it does not convert), so convert GiB to
	// bytes here — mirroring the bytes-to-GiB conversion in the read path. It is
	// carried in AdditionalProperties so it serialises at the top level of the
	// storageVolume object (not inside config).
	if !plan.MaxStorage.IsNull() && !plan.MaxStorage.IsUnknown() {
		body.AdditionalProperties = map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64() * oneGibibyte,
		}
	}

	if !plan.StorageServerId.IsNull() && !plan.StorageServerId.IsUnknown() {
		body.StorageServer = &sdk.AddStorageVolumesRequestStorageVolumeStorageServer{
			Id: plan.StorageServerId.ValueInt64(),
		}
	}

	if !plan.StorageGroupId.IsNull() && !plan.StorageGroupId.IsUnknown() {
		body.StorageGroup = &sdk.AddStorageVolumesRequestStorageVolumeStorageGroup{
			Id: plan.StorageGroupId.ValueInt64(),
		}
	}

	if !plan.ProvisionType.IsNull() && !plan.ProvisionType.IsUnknown() {
		body.ProvisionType = plan.ProvisionType.ValueStringPointer()
	}

	// Build the write-only Config union from req.Config (not plan).
	body.Config = buildCreateConfig(ctx, &config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := client.StorageAPI.AddStorageVolumes(ctx).
		AddStorageVolumesRequest(sdk.AddStorageVolumesRequest{
			StorageVolume: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "storage_volume",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	if result.StorageVolume == nil || result.StorageVolume.Id == nil {
		resp.Diagnostics.AddError(
			"API returned nil ID",
			"StorageVolume ID is nil in the create response",
		)

		return
	}

	id := *result.StorageVolume.Id
	idParam := sdk.GetStorageVolumesIdParameter{Int64: &id}

	readResult, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, idParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "storage_volume",
			plan.Name.ValueString(), err, httpResp,
		)
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

func (r *storageVolumeResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()
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

func (r *storageVolumeResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()
	idParam := sdk.UpdateStorageVolumesIdParameter{Int64: &id}

	body := sdk.UpdateStorageVolumesRequestStorageVolume{
		Name: plan.Name.ValueStringPointer(),
	}

	// max_storage is expressed in GiB; convert to bytes (the API unit), matching
	// the create path.
	if !plan.MaxStorage.IsNull() && !plan.MaxStorage.IsUnknown() {
		body.AdditionalProperties = map[string]interface{}{
			"maxStorage": plan.MaxStorage.ValueInt64() * oneGibibyte,
		}
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageVolumes(ctx, idParam).
		UpdateStorageVolumesRequest(sdk.UpdateStorageVolumesRequest{
			StorageVolume: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpUpdate, "storage_volume",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	readIdParam := sdk.GetStorageVolumesIdParameter{Int64: &id}

	readResult, httpResp, err := client.StorageAPI.GetStorageVolumes(ctx, readIdParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "storage_volume",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	if readResult.StorageVolume == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageVolume is nil in the response")

		return
	}

	mapGetResponseToModel(&plan, readResult.StorageVolume)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageVolumeResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()
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
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// oneGibibyte is the number of bytes in 1 GiB. The Morpheus API stores and
// returns storage volume sizes in bytes, while the resource expresses
// max_storage in GiB.
const oneGibibyte int64 = 1024 * 1024 * 1024

// buildCreateConfig constructs the typed Config union for the create request.
// It handles two scenarios:
//  1. Typed Alletra MP BMaaS block is set → build AlletraMPBMaaSVolumeConfiguration
//  2. Generic dynamic config is set → build map[string]interface{} fallback
//
// max_storage is not part of config; it is sent as the top-level
// storageVolume.maxStorage field (see Create).
func buildCreateConfig(
	ctx context.Context,
	plan *StorageVolumeModel,
	diags *diag.Diagnostics,
) *sdk.AddStorageVolumesRequestStorageVolumeConfig {
	alletraSet := !plan.ConfigAlletrampBmaas.IsNull() && !plan.ConfigAlletrampBmaas.IsUnknown()
	dynamicSet := !plan.Config.IsNull() && !plan.Config.IsUnknown()

	// Scenario 1: Typed Alletra MP BMaaS config block.
	if alletraSet {
		block := plan.ConfigAlletrampBmaas

		variant := sdk.AlletraMPBMaaSVolumeConfiguration{
			HpeStorageDatastore: block.DatastoreId.ValueInt64(),
		}

		if !block.Shared.IsNull() && !block.Shared.IsUnknown() {
			val := boolToOnOff(block.Shared.ValueBool())
			variant.HpeStorageVolumeShared = &val
		}

		if !block.ComputeServerId.IsNull() && !block.ComputeServerId.IsUnknown() {
			val := fmt.Sprintf("[id: %d]", block.ComputeServerId.ValueInt64())
			variant.HpeStorageComputeServer = &val
		}

		if !block.InstanceIds.IsNull() && !block.InstanceIds.IsUnknown() {
			var instanceIDs []int64
			diags.Append(block.InstanceIds.ElementsAs(ctx, &instanceIDs, false)...)
			if diags.HasError() {
				return nil
			}

			instances := make([]string, 0, len(instanceIDs))
			for _, id := range instanceIDs {
				instances = append(instances, fmt.Sprintf("[id: %d]", id))
			}

			variant.HpeStorageInstances = instances
		}

		if !block.RemoteCopyTargetId.IsNull() && !block.RemoteCopyTargetId.IsUnknown() {
			variant.HpeStorageRemotecopytargetId = block.RemoteCopyTargetId.ValueStringPointer()
		}

		if !block.UseExistingVolumeSet.IsNull() && !block.UseExistingVolumeSet.IsUnknown() {
			val := boolToOnOff(block.UseExistingVolumeSet.ValueBool())
			variant.HpeStorageExistingVolumeSet = &val
		}

		if !block.VolumeSetId.IsNull() && !block.VolumeSetId.IsUnknown() {
			variant.HpeStorageVolumesetId = block.VolumeSetId.ValueStringPointer()
		}

		if !block.VolumeSetName.IsNull() && !block.VolumeSetName.IsUnknown() {
			variant.HpeStorageVolumeSetName = block.VolumeSetName.ValueStringPointer()
		}

		return &sdk.AddStorageVolumesRequestStorageVolumeConfig{
			AlletraMPBMaaSVolumeConfiguration: &variant,
		}
	}

	// Scenario 2: Generic dynamic config map.
	if dynamicSet {
		configValue := plan.Config.UnderlyingValue()

		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			diags.AddError(
				"storage_volume config conversion",
				"Failed to convert config dynamic value: "+err.Error(),
			)

			return nil
		}

		configMap, ok := configAny.(map[string]interface{})
		if !ok {
			diags.AddError(
				"storage_volume config type",
				"config must be a valid object/map",
			)

			return nil
		}

		return &sdk.AddStorageVolumesRequestStorageVolumeConfig{
			MapmapOfStringAny: &configMap,
		}
	}

	// Nothing set — leave Config nil.
	return nil
}

// boolToOnOff converts a Go bool to the "on"/"off" string the Alletra plugin expects.
func boolToOnOff(b bool) string {
	if b {
		return "on"
	}

	return "off"
}

// objectID extracts the numeric "id" from a Morpheus relation object that the
// read API encodes as a generic JSON object (e.g. storageServer, storageGroup).
// JSON numbers decode into float64, so the value is asserted as float64 and
// narrowed to int64. The bool reports whether a usable id was present; a nil map
// or a missing/non-numeric "id" yields false, letting callers leave the
// corresponding optional attribute untouched.
func objectID(obj map[string]interface{}) (int64, bool) {
	id, ok := obj["id"].(float64)
	if !ok {
		return 0, false
	}

	return int64(id), true
}

// mapGetResponseToModel maps the API read response onto the Terraform model.
// Write-only config attributes (config, config_alletramp_bmaas) are NOT overwritten —
// the API read does not return config, so we preserve whatever the plan/state holds.
func mapGetResponseToModel(
	model *StorageVolumeModel,
	sv *sdk.GetStorageVolumes200ResponseStorageVolume,
) {
	if sv.Id != nil {
		model.Id = types.Int64Value(*sv.Id)
	}

	if sv.Name != nil {
		model.Name = types.StringValue(*sv.Name)
	}

	// type_id and type_code are both computed_optional and mutually exclusive.
	// Populate both from the API so the plan is clean regardless of which one the
	// user configured — the unconfigured (computed) value is not compared against
	// the null config, so it does not cause spurious drift or replacement.
	if t := sv.Type; t != nil {
		model.TypeId = convert.Int64ToType(t.Id)
		model.TypeCode = convert.StrToType(t.Code)
	} else if sv.TypeId != nil {
		model.TypeId = types.Int64Value(*sv.TypeId)
	}

	if id, ok := objectID(sv.StorageServer); ok {
		model.StorageServerId = types.Int64Value(id)
	}

	// The API returns maxStorage in bytes; the resource expresses it in GiB.
	if sv.MaxStorage != nil {
		model.MaxStorage = types.Int64Value(*sv.MaxStorage / oneGibibyte)
	}

	if sv.Status != nil {
		model.Status = types.StringValue(*sv.Status)
	}

	// ProvisionType is a NullableString in the read model.
	if sv.ProvisionType.IsSet() {
		model.ProvisionType = convert.StrToType(sv.ProvisionType.Get())
	}

	// wwn is a computed, read-only identifier assigned by the storage system.
	model.Wwn = convert.StrToType(sv.Wwn.Get())

	// StorageGroup is an object in the read model; map its id back to
	// storage_group_id so it round-trips (and imports) cleanly.
	if id, ok := objectID(sv.StorageGroup); ok {
		model.StorageGroupId = types.Int64Value(id)
	}
}
