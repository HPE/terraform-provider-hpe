// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// ModifyPlan enforces the HPE Alletra MP / 9000 max_storage upper bound (65536
// GiB) at plan time for volumes whose type is supplied as type_id.
//
// The static max_storage attribute validator (maxStorageSizeValidator) can only
// apply the Alletra ceiling when the type is given as type_code, because a
// schema validator has no API client to resolve a numeric type_id. Here the
// resource is configured, so when only type_id is set we resolve it to its type
// code and apply the same bound (via maxStorageSizeError) -- surfacing an
// out-of-range size at plan time instead of as a 400 from the API at apply.
//
// This reads only the planned config (never req.State), so no req.State guard is
// required and the check runs on create as well as update. The lower bound
// (>= 1 GiB) is already enforced for every type by the attribute validator.
//
// Plan-time validation is best effort: if the client is unavailable or the type
// lookup fails, the plan is not blocked -- Morpheus still enforces the limit on
// apply.
func (r *storageVolumeResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Nothing to validate on destroy (no planned config).
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan StorageVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only resolve a type when max_storage is a known value and the type is
	// given as type_id. A known type_code is already validated statically by the
	// attribute validator (no API call needed).
	if plan.MaxStorage.IsNull() || plan.MaxStorage.IsUnknown() {
		return
	}
	if !plan.TypeCode.IsNull() && !plan.TypeCode.IsUnknown() {
		return
	}
	if plan.TypeId.IsNull() || plan.TypeId.IsUnknown() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		// Best effort: let create/update surface any client configuration error.
		return
	}

	typeID := plan.TypeId.ValueInt64()

	typeResult, httpResp, err := client.StorageAPI.
		GetStorageVolumeTypes(ctx, typeID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("max_storage"),
			"Could not verify max_storage against the storage volume type",
			"Unable to look up storage volume type "+strconv.FormatInt(typeID, 10)+
				" to validate max_storage at plan time; Morpheus will enforce the "+
				"size limit on apply. Error: "+err.Error(),
		)

		return
	}

	if typeResult == nil || typeResult.StorageVolumeType == nil ||
		typeResult.StorageVolumeType.Code == nil {
		return
	}

	if title, detail := maxStorageSizeError(
		plan.MaxStorage.ValueInt64(), *typeResult.StorageVolumeType.Code, true,
	); title != "" {
		resp.Diagnostics.AddAttributeError(path.Root("max_storage"), title, detail)
	}
}
