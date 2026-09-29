// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the nested write-only server.ssh_password_wo: see
// wodigest. Morpheus sets a cluster's SSH password at creation and has no
// update for it, so a change to the configured value plans a replacement of
// the cluster -- the only way the new value can take effect -- rather than an
// update that would silently do nothing. A change to the deprecated
// server.ssh_password_wo_version marker plans the same replacement for the
// whole 2.2 line.
//
// The cluster resource is not importable, so the import-ID digests do not
// apply here.
var wo = wodigest.Tracker{
	Resource: "cluster",
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"server.ssh_password_wo": "server.ssh_password_wo_version",
	},
	Immutable: []string{"server.ssh_password_wo"},
}

// ModifyPlan plans a replacement when server.ssh_password_wo changed (or its
// deprecated marker did).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
