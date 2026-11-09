// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package affinityshape enforces, at plan time, the membership rule Morpheus
// 9.1 introduced for affinity groups (MORPH-20083).
//
// From 9.1.0 the AffinityGroup domain validates its "rule shape": a group is
// either a VM-to-VM rule with at least one server, or a VM-to-Host rule with a
// VM Group and a Host Group and no servers. An empty group is neither and is
// refused -- on create ("An affinity rule requires either a Servers list ...";
// for a cloud-scoped group without a pool, "Servers or a VM Group / Host Group
// are required ...") and on an update that removes the last server. Morpheus
// 9.0 accepted empty groups.
//
// The provider models membership with one hpe_morpheus_*_affinity_group_member
// resource per server, so on 9.1 the group must be seeded at create through
// servers, and the last member cannot be removed while the group
// stands. These checks surface both at plan time, before the API refuses them
// mid-apply. They follow versioncheck's fail-open policy: when the appliance
// version cannot be read nothing is refused, and the API's own error stands.
package affinityshape

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

// Enforced reports whether the appliance enforces the 9.1 rule shape. A nil
// version (lookup failed) is not enforced: fail open, as versioncheck does.
func Enforced(ctx context.Context, v *versioncheck.Version, lookupErr error) bool {
	if lookupErr != nil || v == nil {
		tflog.Warn(ctx, "affinity rule-shape gate: appliance version unknown; not enforcing", map[string]any{
			"error": fmt.Sprint(lookupErr),
		})

		return false
	}
	ok, err := versioncheck.Satisfies(v, constants.AffinityRuleShapeMinVersion)
	if err != nil {
		tflog.Warn(ctx, "affinity rule-shape gate: bad constraint; not enforcing", map[string]any{"error": err.Error()})

		return false
	}

	return ok
}

// DecideCreate is the pure decision for a group create (or replacement): on
// an enforcing appliance a configuration that leaves servers unset, or sets
// it empty, is refused with an error on the attribute -- Morpheus 9.1 would
// refuse the memberless group at apply. A wholly unknown set, or one with
// unknown elements, passes: the API decides at apply.
func DecideCreate(ctx context.Context, servers types.Set, v *versioncheck.Version, lookupErr error) diag.Diagnostics {
	var diags diag.Diagnostics

	if !Enforced(ctx, v, lookupErr) {
		return diags
	}
	// servers is Optional+Computed: unconfigured it is UNKNOWN in the plan,
	// not null, so the caller passes the CONFIG value, where unconfigured is
	// null. A wholly unknown configured set (e.g. a whole for_each output) and
	// any element, known or not, pass -- the recommended form is
	// servers = [hpe_morpheus_instance.x.compute_servers[0]], whose value is
	// unknown while the VM is still being created.
	if servers.IsUnknown() || (!servers.IsNull() && len(servers.Elements()) > 0) {
		return diags
	}

	diags.AddAttributeError(path.Root("servers"),
		"Affinity groups need a member on this appliance",
		fmt.Sprintf("Morpheus %s refuses to create an affinity group with no servers (an affinity rule needs either a "+
			"servers list or a VM Group and Host Group). Set servers to at least one compute server id -- for a VM "+
			"managed here, hpe_morpheus_instance.<name>.compute_servers[0] -- so the group is created with its first "+
			"members; this resource then owns the membership. Or create the group after the servers it will hold.",
			v.Original()),
	)

	return diags
}

// DecideLastMemberRemoval decides what a member resource's Delete does when
// the server is the group's only remaining member on an enforcing appliance:
// the API would refuse the emptied group, so the membership write is skipped
// (skip == true) and a warning says the server stays in the group until the
// group itself is destroyed. Terraform destroys a member before its group, so
// a plan-time refusal here would make destroying both together impossible;
// skipping keeps `terraform destroy` working and the appliance truthful.
func DecideLastMemberRemoval(
	ctx context.Context, groupServers []int64, serverID int64, v *versioncheck.Version, lookupErr error,
) (skip bool, diags diag.Diagnostics) {
	if !Enforced(ctx, v, lookupErr) {
		return false, diags
	}
	if len(groupServers) != 1 || groupServers[0] != serverID {
		return false, diags
	}

	diags.AddWarning(
		"Last server left in the affinity group",
		fmt.Sprintf("Morpheus %s refuses to leave an affinity group with no servers, so server %d was NOT removed "+
			"from the group and is no longer managed by Terraform: it stays as the group's only member and the group's "+
			"servers attribute keeps reporting it. This is what a destroy of the group and its members needs (the "+
			"group is destroyed next and takes the membership with it). If you meant to keep the group without this "+
			"server, add another member first and re-apply, or re-create the group with servers set to the "+
			"members you want.", v.Original(), serverID),
	)

	return true, diags
}

// RequireCreate is the I/O half of DecideCreate: it fetches the appliance
// version and decides. Call it from ModifyPlan when the resource is being
// created or replaced, passing the CONFIG value of servers (see DecideCreate).
func RequireCreate(ctx context.Context, client *sdk.APIClient, servers types.Set) diag.Diagnostics {
	v, err := versioncheck.Appliance(ctx, client)

	return DecideCreate(ctx, servers, v, err)
}

// SkipLastMemberRemoval is the I/O half of DecideLastMemberRemoval, for a
// member resource's Delete once it has read the group's current servers.
func SkipLastMemberRemoval(
	ctx context.Context, client *sdk.APIClient, groupServers []int64, serverID int64,
) (bool, diag.Diagnostics) {
	v, err := versioncheck.Appliance(ctx, client)

	return DecideLastMemberRemoval(ctx, groupServers, serverID, v, err)
}

// Servers returns the known int64 elements of an servers set.
func Servers(set types.Set) []int64 {
	return knownElements(set)
}

func knownElements(set types.Set) []int64 {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	out := make([]int64, 0, len(set.Elements()))
	for _, e := range set.Elements() {
		if v, ok := e.(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
			out = append(out, v.ValueInt64())
		}
	}

	return out
}

// SetOf builds an int64 set value from ids (test helper and request builder).
func SetOf(ids ...int64) types.Set {
	vals := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		vals = append(vals, types.Int64Value(id))
	}

	return types.SetValueMust(types.Int64Type, vals)
}
