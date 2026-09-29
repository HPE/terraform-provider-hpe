// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes access_key, secret_key: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured values. The deprecated
// *_version markers keep working for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "storage_bucket",
	Attributes: map[string]string{
		"access_key": "access_key_version",
		"secret_key": "secret_key_version",
	},
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did).
func (r *storageBucketResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
