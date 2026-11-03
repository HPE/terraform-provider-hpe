// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
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

	var state StorageVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	// The only in-place change the API offers is a resize; the volume's name
	// cannot be edited after create (the schema replaces the volume when it
	// changes). Morpheus 9.1.0 repurposed PUT /api/storage-volumes/{id} for
	// tenant transfer, so the resize has its own endpoint. max_storage is
	// expressed in GiB; convert to bytes (the API unit), matching create.
	if !plan.MaxStorage.IsNull() && !plan.MaxStorage.IsUnknown() && !plan.MaxStorage.Equal(state.MaxStorage) {
		idParam := sdk.TransferStorageVolumeOwnershipIdParameter{Int64: &id}
		body := sdk.ResizeStorageVolumeRequest{
			StorageVolume: sdk.ResizeStorageVolumeRequestStorageVolume{
				MaxStorage: plan.MaxStorage.ValueInt64() * oneGibibyte,
			},
		}

		_, httpResp, err := client.StorageAPI.ResizeStorageVolume(ctx, idParam).
			ResizeStorageVolumeRequest(body).Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(
				&resp.Diagnostics, errfmt.OpUpdate, "storage_volume",
				plan.Name.ValueString(), err, httpResp,
			)

			return
		}
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

	// The write-only configuration is create-only (a change replaces the
	// volume); an update planned through write_only_revision only records the
	// digests of the configured values.
	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
