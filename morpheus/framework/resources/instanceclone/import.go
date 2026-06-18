// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	instanceID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import instance clone",
			fmt.Sprintf("invalid instance id %q: %v", req.ID, err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), instanceID)...)
	// source_instance_id is not returned by the clone API as a first-class
	// field. Set it null here; Read makes a best-effort attempt to recover it
	// from the clone's config (cloneInstanceId). If the platform does not expose
	// that field, set source_instance_id in configuration after import.
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("source_instance_id"), types.Int64Null(),
	)...)
}

// sourceInstanceIDFromConfig returns cloneInstanceId - the source instance id
// that Morpheus stamps onto a clone's config during cloning - from the instance
// read response. It is absent for instances that were not created via the clone
// endpoint, in which case the second return value is false.
func sourceInstanceIDFromConfig(inst *sdk.GetInstance200ResponseInstance) (int64, bool) {
	if inst == nil || inst.Config == nil || inst.Config.CloneInstanceId == nil {
		return 0, false
	}

	return *inst.Config.CloneInstanceId, true
}
