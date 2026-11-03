// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute source_instance_id: see
// wodigest. Morpheus takes the source instance only at creation, so a change plans a replacement --
// the only way the new value can take effect -- rather than an update that would silently do
// nothing. The deprecated *_version marker keeps working for the
// whole 2.2 line.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "instance_clone",
	Attributes: map[string]string{
		"source_instance_id": "source_instance_id_version",
	},
	Immutable: []string{"source_instance_id"},
}

// ModifyPlan plans a replacement when a write-only value changed (or a
// deprecated marker did), and refuses a volume renamed or the volumes
// reordered in place, neither of which can be carried out (see
// rejectVolumeNameChanges in update.go).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
	rejectVolumeNameChanges(ctx, req, resp)
}
