// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

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
