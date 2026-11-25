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

// hvmPlacementValidator enforces the HVM host-placement rules at plan time.
//
//   - preferred_host_id set but placement_strategy null -> error
//   - placement_strategy "failover" with preferred_host_id null -> error
//   - placement_strategy "pinned" with affinity_group_id set -> warning
//   - kvm_host_id and preferred_host_id set and different -> warning
//
// Every check is skipped when its operands are unknown.
type hvmPlacementValidator struct{}

func (v hvmPlacementValidator) Description(_ context.Context) string {
	return "config_hvm.placement_strategy and config_hvm.preferred_host_id must be consistent"
}

func (v hvmPlacementValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v hvmPlacementValidator) ValidateResource(
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

	ps := config.ConfigHvm.PlacementStrategy
	ph := config.ConfigHvm.PreferredHostId
	ag := config.ConfigHvm.AffinityGroupId
	kh := config.ConfigHvm.KvmHostId

	strategyKnown := !ps.IsUnknown()
	strategySet := strategyKnown && !ps.IsNull()
	preferredKnown := !ph.IsUnknown()
	preferredSet := preferredKnown && !ph.IsNull()

	// preferred_host_id set but placement_strategy null -> error.
	if preferredSet && strategyKnown && ps.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_hvm").AtName("placement_strategy"),
			"Missing placement strategy",
			"config_hvm.placement_strategy must be set when preferred_host_id is set "+
				"(failover or pinned); a host preference without a strategy would be "+
				"ignored or cause a surprise migration.",
		)
	}

	// failover with no preferred host -> error.
	if strategySet && ps.ValueString() == "failover" && preferredKnown && ph.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_hvm").AtName("preferred_host_id"),
			"Missing preferred host for failover",
			"config_hvm.placement_strategy = \"failover\" requires config_hvm.preferred_host_id "+
				"to name the host to fail over to.",
		)
	}

	// auto with a preferred host -> error: the host would be sent nowhere
	// (an auto placement carries no host) and silently ignored.
	if strategySet && ps.ValueString() == placementStrategyAuto && preferredSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_hvm").AtName("preferred_host_id"),
			"Preferred host is ignored under auto placement",
			"config_hvm.placement_strategy = \"auto\" releases any host preference, so "+
				"config_hvm.preferred_host_id has no effect; remove it or choose \"failover\" or \"pinned\".",
		)
	}

	// pinned AND affinity_group_id -> warning.
	if strategySet && ps.ValueString() == "pinned" &&
		!ag.IsUnknown() && !ag.IsNull() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("config_hvm").AtName("placement_strategy"),
			"Pin outranks affinity group",
			"the pin outranks the affinity group; Morpheus keeps the VM on the pinned "+
				"host and raises an affinity-violation alarm instead of moving it.",
		)
	}

	// kvm_host_id and preferred_host_id set and different, under a strategy
	// that honours the host -> warning (under auto the host is rejected above).
	if !kh.IsUnknown() && !kh.IsNull() && preferredSet &&
		strategySet && ps.ValueString() != placementStrategyAuto &&
		kh.ValueInt64() != ph.ValueInt64() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("config_hvm").AtName("preferred_host_id"),
			"Initial placement differs from steady-state preference",
			"initial placement and steady-state preference differ; the VM will be "+
				"migrated right after provisioning.",
		)
	}
}
