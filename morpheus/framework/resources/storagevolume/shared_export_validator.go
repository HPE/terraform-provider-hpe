// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sharedAttrPath points at config_alletramp_bmaas.shared, which the export-target
// validators below consult to enforce the shared/target pairing. shared is a
// write-only attribute, but its value is present in the configuration and so is
// readable during validation.
func sharedAttrPath() path.Path {
	return path.Root("config_alletramp_bmaas").AtName("shared")
}

// requireSharedForInstanceIDsValidator enforces shared = true whenever
// instance_ids is set. instance_ids exports the volume to every node of each
// listed instance (multi-attach), so the volume must be marked shared.
type requireSharedForInstanceIDsValidator struct{}

// requireSharedForInstanceIDs returns a validator that requires shared = true
// when instance_ids is set.
func requireSharedForInstanceIDs() validator.List {
	return requireSharedForInstanceIDsValidator{}
}

func (v requireSharedForInstanceIDsValidator) Description(_ context.Context) string {
	return "shared must be true when instance_ids is set"
}

func (v requireSharedForInstanceIDsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requireSharedForInstanceIDsValidator) ValidateList(
	ctx context.Context,
	req validator.ListRequest,
	resp *validator.ListResponse,
) {
	// Only applies when instance_ids is actually set.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() ||
		len(req.ConfigValue.Elements()) == 0 {
		return
	}

	var shared types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, sharedAttrPath(), &shared)...)
	if resp.Diagnostics.HasError() || shared.IsUnknown() {
		return
	}

	if shared.IsNull() || !shared.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"shared is required with instance_ids",
			"instance_ids exports the volume to every node of each instance "+
				"(multi-attach); set shared = true.",
		)
	}
}

// forbidSharedForComputeServerValidator enforces that shared is not true when
// compute_server_id is set. compute_server_id exports the volume to a single
// host (single-attach), which is incompatible with a shared (multi-attach)
// volume.
type forbidSharedForComputeServerValidator struct{}

// forbidSharedForComputeServer returns a validator that rejects shared = true
// when compute_server_id is set.
func forbidSharedForComputeServer() validator.Int64 {
	return forbidSharedForComputeServerValidator{}
}

func (v forbidSharedForComputeServerValidator) Description(_ context.Context) string {
	return "shared must not be true when compute_server_id is set"
}

func (v forbidSharedForComputeServerValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v forbidSharedForComputeServerValidator) ValidateInt64(
	ctx context.Context,
	req validator.Int64Request,
	resp *validator.Int64Response,
) {
	// Only applies when compute_server_id is actually set.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var shared types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, sharedAttrPath(), &shared)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !shared.IsNull() && !shared.IsUnknown() && shared.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"shared conflicts with compute_server_id",
			"compute_server_id exports the volume to a single host "+
				"(single-attach); shared must not be true.",
		)
	}
}
