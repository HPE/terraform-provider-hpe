// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// storageVolumeAsState maps an API storage volume into the datasource model.
func storageVolumeAsState(
	sv *sdk.GetStorageVolumes200ResponseStorageVolume,
) StorageVolumeModel {
	state := StorageVolumeModel{
		Id:                   convert.Int64ToType(sv.Id),
		Name:                 convert.StrToType(sv.Name),
		Description:          convert.StrToType(sv.Description.Get()),
		ControllerId:         convert.Int64ToType(sv.ControllerId.Get()),
		ControllerMountPoint: convert.StrToType(sv.ControllerMountPoint.Get()),
		Resizeable:           convert.BoolToType(sv.Resizeable.Get()),
		RootVolume:           convert.BoolToType(sv.RootVolume),
		UnitNumber:           convert.StrToType(sv.UnitNumber.Get()),
		DeviceName:           convert.StrToType(sv.DeviceName),
		DeviceDisplayName:    convert.StrToType(sv.DeviceDisplayName),
		VolumeName:           convert.StrToType(sv.VolumeName),
		VolumePath:           convert.StrToType(sv.VolumePath),
		VolumeType:           convert.StrToType(sv.VolumeType),
		RefType:              convert.StrToType(sv.RefType),
		RefId:                convert.Int64ToType(sv.RefId),
		DiskMode:             convert.StrToType(sv.DiskMode),
		DiskType:             convert.StrToType(sv.DiskType),
		TypeId:               convert.Int64ToType(sv.TypeId),
		Category:             convert.StrToType(sv.Category),
		Status:               convert.StrToType(sv.Status),
		StatusMessage:        convert.StrToType(sv.StatusMessage.Get()),
		ConfigurableIops:     convert.BoolToType(sv.ConfigurableIOPS),
		MaxStorage:           convert.Int64ToType(sv.MaxStorage),
		UsedStorage:          convert.Int64ToType(sv.UsedStorage),
		DisplayOrder:         convert.Int64ToType(sv.DisplayOrder),
		MaxIops:              convert.StrToType(sv.MaxIOPS.Get()),
		Uuid:                 convert.StrToType(sv.Uuid),
		Active:               convert.BoolToType(sv.Active),
		ReadOnly:             convert.BoolToType(sv.ReadOnly),
		Removable:            convert.BoolToType(sv.Removable),
		Online:               convert.BoolToType(sv.Online),
		PoolName:             convert.StrToType(sv.PoolName),
		DatastoreId:          convert.Int64ToType(sv.DatastoreId.Get()),
		DatastoreOption:      convert.StrToType(sv.DatastoreOption),
		Namespace:            convert.StrToType(sv.Namespace.Get()),
		Source:               convert.StrToType(sv.Source),
		UniqueId:             convert.StrToType(sv.UniqueId.Get()),
		InternalId:           convert.StrToType(sv.InternalId.Get()),
		ExternalId:           convert.StrToType(sv.ExternalId),
		ProvisionType:        convert.StrToType(sv.ProvisionType.Get()),
		CopyType:             convert.StrToType(sv.CopyType.Get()),
		FiberWwn:             convert.StrToType(sv.FiberWwn.Get()),
		Wwn:                  convert.StrToType(sv.Wwn.Get()),
		FileName:             convert.StrToType(sv.FileName.Get()),
		ClaimName:            convert.StrToType(sv.ClaimName.Get()),
		SharePath:            convert.StrToType(sv.SharePath.Get()),
		SourceId:             convert.StrToType(sv.SourceId),
		ImageType:            convert.StrToType(sv.ImageType),
		CreateForMultiAttach: convert.BoolToType(sv.CreateForMultiAttach),
		IsMultiAttach:        convert.BoolToType(sv.IsMultiAttach),
		StorageProfile:       convert.StrToType(sv.StorageProfile.Get()),
	}

	// Type nested object
	state.Type = types.StringNull()
	state.TypeCode = types.StringNull()
	state.TypeName = types.StringNull()
	if sv.Type != nil {
		state.Type = convert.StrToType(sv.Type.Name)
		state.TypeCode = convert.StrToType(sv.Type.Code)
		state.TypeName = convert.StrToType(sv.Type.Name)
	}

	// Zone -> cloud_id / cloud_name
	state.CloudId = convert.Int64ToType(sv.ZoneId)
	state.CloudName = types.StringNull()
	if sv.Zone != nil {
		state.CloudName = convert.StrToType(sv.Zone.Name)
	}

	// Datastore -> datastore_name
	state.DatastoreName = types.StringNull()
	if sv.Datastore != nil {
		state.DatastoreName = convert.StrToType(sv.Datastore.Name)
	}

	// StorageServer -> storage_server_id / storage_server_name
	state.StorageServerId = types.Int64Null()
	state.StorageServerName = types.StringNull()
	if sv.StorageServer != nil {
		state.StorageServerId = convert.Int64ToType(sv.StorageServer.Id)
		state.StorageServerName = convert.StrToType(sv.StorageServer.Name)
	}

	// StorageGroup -> storage_group_id / storage_group_name
	state.StorageGroupId = types.Int64Null()
	state.StorageGroupName = types.StringNull()
	if sv.StorageGroup != nil {
		state.StorageGroupId = convert.Int64ToType(sv.StorageGroup.Id)
		state.StorageGroupName = convert.StrToType(sv.StorageGroup.Name)
	}

	return state
}

func getStorageVolumeByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageVolumes200ResponseStorageVolume, error) {
	idParam := sdk.GetStorageVolumesIdParameter{Int64: &id}
	r, hresp, err := apiClient.StorageAPI.GetStorageVolumes(ctx, idParam).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for storage volume %d: %s",
			id, providererrors.ErrMsg(err, hresp))
	}

	if r.StorageVolume == nil {
		return nil, fmt.Errorf(
			"GET failed for storage volume %d: response missing storageVolume", id,
		)
	}

	return r.StorageVolume, nil
}

func getStorageVolumeByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageVolumes200ResponseStorageVolume, error) {
	rs, hresp, err := apiClient.StorageAPI.ListStorageVolumes(ctx).Name(name).Max(listMax).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for storage volume %s: %s",
			name, providererrors.ErrMsg(err, hresp))
	}

	var matchedID int64

	var matchCount int

	for _, sv := range rs.StorageVolumes {
		if sv.Name != nil && *sv.Name == name {
			if sv.Id != nil {
				matchedID = *sv.Id
				matchCount++
			}
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoStorageVolumeFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleStorageVolumes)
	}

	return getStorageVolumeByID(ctx, matchedID, apiClient)
}

func getStorageVolume(
	ctx context.Context,
	config *StorageVolumeModel,
	apiClient *sdk.APIClient,
) (*sdk.GetStorageVolumes200ResponseStorageVolume, error) {
	if !config.Id.IsNull() {
		return getStorageVolumeByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getStorageVolumeByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config StorageVolumeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	sv, err := getStorageVolume(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := storageVolumeAsState(sv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
