// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute service_password_wo: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value. The deprecated
// *_version marker keeps working for the whole 2.2 line.
//
// APIHash additionally detects a change made outside Terraform. Morpheus
// returns a deterministic hash of the stored secret; we record it right after a
// send and compare it only for equality across later reads -- the hash
// algorithm is never assumed. OnDrift re-sends the configured value when the
// appliance's stored secret has changed.
var wo = wodigest.Tracker{
	Resource: "network_pool_server",
	Attributes: map[string]string{
		"service_password_wo": "service_password_wo_version",
	},
	APIHash: map[string]wodigest.HashExtractor{
		"service_password_wo": servicePasswordHash,
	},
	OnDrift: wodigest.ResendPolicy,
}

// servicePasswordHash pulls the appliance's hash of the stored service password
// from the read-back object. The field is a plain *string, which cannot tell
// an omitted field from a null one, so a nil pointer is read as not-present:
// a pool server always has a service password, so a missing hash means the
// appliance did not render the field, and nothing can be compared.
func servicePasswordHash(resp any) (string, bool) {
	server, ok := resp.(*sdk.GetNetworkPoolServer200ResponseNetworkPoolServer)
	if !ok || server == nil || server.ServicePasswordHash == nil {
		return "", false
	}

	return *server.ServicePasswordHash, true
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did).
func (r *networkPoolServerResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
