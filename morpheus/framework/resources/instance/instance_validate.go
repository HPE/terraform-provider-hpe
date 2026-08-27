// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// hvmHostAffinityConflictValidator rejects configs that set both
// config_hvm.kvm_host_id and config_hvm.affinity_group_id simultaneously.
// Morpheus resolves host selection by applying the pinned host first and then
// overriding it with the affinity group's existing host, so setting both means
// the host pin is silently discarded.
type hvmHostAffinityConflictValidator struct{}

func (v hvmHostAffinityConflictValidator) Description(_ context.Context) string {
	return "config_hvm.kvm_host_id and config_hvm.affinity_group_id are mutually exclusive"
}

func (v hvmHostAffinityConflictValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v hvmHostAffinityConflictValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config InstanceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigHvm.IsNull() || config.ConfigHvm.IsUnknown() {
		return
	}

	hostSet := !config.ConfigHvm.KvmHostId.IsNull() && !config.ConfigHvm.KvmHostId.IsUnknown()
	affinitySet := !config.ConfigHvm.AffinityGroupId.IsNull() && !config.ConfigHvm.AffinityGroupId.IsUnknown()

	if hostSet && affinitySet {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_hvm").AtName("kvm_host_id"),
			"Conflicting host placement attributes",
			"config_hvm.kvm_host_id and config_hvm.affinity_group_id cannot both be set. "+
				"Morpheus lets the affinity group silently override an explicit host pin — "+
				"the instance would land on a host chosen by the affinity group, not the one "+
				"you specified. Remove one of the two attributes.",
		)
	}
}
