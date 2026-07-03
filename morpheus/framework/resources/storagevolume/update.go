// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

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
		if body.AdditionalProperties == nil {
			body.AdditionalProperties = make(map[string]interface{})
		}
		body.AdditionalProperties["maxStorage"] = plan.MaxStorage.ValueInt64() * oneGibibyte
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
