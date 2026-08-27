// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package versioncheck

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Require is a hard version gate: it refuses to let an operation proceed when
// the target appliance is older than a feature needs.
//
// It is the I/O half of the gate — it fetches the appliance version and hands
// the result, error and all, to Decide. Call it at the top of a CRUD method,
// once the API client exists:
//
//	resp.Diagnostics.Append(versioncheck.Require(
//	    ctx, client, "Cloud affinity groups", constants.AffinityGroupMinVersion,
//	)...)
//	if resp.Diagnostics.HasError() {
//	    return
//	}
//
// Deliberately NOT called from Configure. terraform-plugin-framework runs
// Configure on every RPC for the type, including ValidateResourceConfig and
// UpgradeResourceState, so gating there would put a network round trip behind
// `terraform validate` and behind state upgrades. Gating in the CRUD methods
// keeps the call on operations that were always going to talk to the API
// anyway, and costs exactly one extra request per operation.
//
// feature names the thing being gated, phrased as a plural noun so it reads
// correctly in the diagnostic ("Cloud affinity groups require ..."). constraint
// is a go-version constraint such as ">= 8.0.10".
func Require(
	ctx context.Context,
	client *sdk.APIClient,
	feature string,
	constraint string,
) diag.Diagnostics {
	v, err := Appliance(ctx, client)

	return Decide(ctx, feature, constraint, v, err)
}

// Decide is the pure decision half of Require, split out so the gate's
// behaviour can be unit tested without a live appliance. It returns an error
// diagnostic if, and only if, the appliance version is known and fails the
// constraint.
//
// FAIL OPEN when the version cannot be determined.
//
// The version comes from GET /api/health, which is permission-guarded: the API
// only serves it to callers holding the `admin-health` permission at read or
// full access. A user with every permission needed to manage the gated resource
// can still be refused the health endpoint, and a reverse proxy or a transient
// outage can hide it too. Failing closed would turn any of those into a total,
// unworkaroundable failure against a perfectly healthy 8.0.10+ appliance — the
// practitioner would have to go and get an unrelated permission granted before
// Terraform would run at all.
//
// Failing open is strictly no worse than the pre-gate behaviour: the operation
// proceeds and, if the appliance really is too old, the API returns the same
// confusing error it always did. The gate is an improvement in the common case
// and a no-op in the case it cannot judge. That also matches how the instance
// resource treats an unreadable version (skip, do not block), so the provider
// behaves consistently.
//
// The skip is logged rather than surfaced as a warning diagnostic: a warning
// would be emitted on every read of every gated resource and data source, which
// is noise a practitioner cannot act on.
func Decide(
	ctx context.Context,
	feature string,
	constraint string,
	v *Version,
	lookupErr error,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if lookupErr != nil || v == nil {
		reason := "appliance did not report a version"
		if lookupErr != nil {
			reason = lookupErr.Error()
		}

		tflog.Warn(ctx, "skipping appliance version gate: version unavailable", map[string]any{
			"feature":    feature,
			"constraint": constraint,
			"reason":     reason,
		})

		return diags
	}

	ok, err := Satisfies(v, constraint)
	if err != nil {
		// An unparseable constraint is a provider bug, not an appliance
		// problem. Blocking the practitioner for our own typo helps nobody, so
		// this takes the same fail-open path.
		tflog.Warn(ctx, "skipping appliance version gate: malformed constraint", map[string]any{
			"feature":    feature,
			"constraint": constraint,
			"reason":     err.Error(),
		})

		return diags
	}

	if ok {
		return diags
	}

	diags.AddError(
		fmt.Sprintf("%s require a newer Morpheus appliance", feature),
		fmt.Sprintf(
			"%s require a Morpheus appliance version %s. This appliance reports "+
				"version %s. Upgrade the appliance to a version satisfying %s, or "+
				"remove the affected configuration.",
			feature, constraint, v.Original(), constraint,
		),
	)

	return diags
}
