// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

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
	model.Id = convert.Int64ToType(sv.Id)
	model.Name = convert.StrToType(sv.Name)

	// type_id and type_code are both computed_optional and mutually exclusive.
	// Populate both from the API so the plan is clean regardless of which one the
	// user configured — the unconfigured (computed) value is not compared against
	// the null config, so it does not cause spurious drift or replacement. A value
	// the API does not return is mapped to a known null rather than left unknown
	// (which would fail apply for a computed attribute).
	if t := sv.Type; t != nil {
		model.TypeId = convert.Int64ToType(t.Id)
		model.TypeCode = convert.StrToType(t.Code)
	} else if sv.TypeId != nil {
		model.TypeId = convert.Int64ToType(sv.TypeId)
	}
	if model.TypeId.IsUnknown() {
		model.TypeId = types.Int64Null()
	}
	if model.TypeCode.IsUnknown() {
		model.TypeCode = types.StringNull()
	}

	if id, ok := objectID(sv.StorageServer); ok {
		model.StorageServerId = types.Int64Value(id)
	}

	// The API returns maxStorage in bytes; the resource expresses it in GiB.
	// max_storage is optional+computed, so an unset value must map to a known
	// null rather than remain unknown (which would fail apply with "computed
	// attribute remained unknown").
	if sv.MaxStorage != nil {
		model.MaxStorage = types.Int64Value(*sv.MaxStorage / oneGibibyte)
	} else {
		model.MaxStorage = types.Int64Null()
	}

	// status is computed; an unset value maps to a known null.
	model.Status = convert.StrToType(sv.Status)

	// ProvisionType is an optional, computed NullableString. As with wwn below, an
	// unset value must be mapped to a known null rather than left unknown —
	// otherwise apply fails with "computed attribute remained unknown" when the
	// user did not configure provision_type and the API omits it.
	if sv.ProvisionType.IsSet() {
		model.ProvisionType = convert.StrToType(sv.ProvisionType.Get())
	} else {
		model.ProvisionType = types.StringNull()
	}

	// wwn is a computed, read-only identifier assigned by the storage system.
	// The storage system may assign it asynchronously, so an unset value is
	// mapped to a known null rather than left unknown (which would fail apply for
	// this computed-only attribute).
	if sv.Wwn.IsSet() {
		model.Wwn = convert.StrToType(sv.Wwn.Get())
	} else {
		model.Wwn = types.StringNull()
	}

	// StorageGroup is an object in the read model; map its id back to
	// storage_group_id so it round-trips (and imports) cleanly.
	if id, ok := objectID(sv.StorageGroup); ok {
		model.StorageGroupId = types.Int64Value(id)
	}
}
