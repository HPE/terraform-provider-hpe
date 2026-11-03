// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Appliance-side drift detection.
//
// The provider's digest (see wodigest.go) answers "did the configured value
// change in HCL?". It is blind to a change made on the appliance (UI, API,
// another tool). Morpheus returns a deterministic hash of each stored secret
// beside its masked value; the provider never computes or interprets that
// hash, it only records the hash read back right after it last sent a value
// (Observed) and compares it for string equality on every later read
// (CheckDrift). An inequality means the stored secret changed out of band.
//
// APIHash names the attributes that take part, with an extractor per
// attribute; an attribute without one keeps today's config-side behaviour
// exactly.

// policyFor returns the drift policy for one attribute. An immutable attribute
// is always WarnPolicy: re-sending its value would be a replacement, which
// must never happen silently.
func (t Tracker) policyFor(attr string) DriftPolicy {
	if contains(t.Immutable, attr) {
		return WarnPolicy
	}

	return t.OnDrift
}

// Observe records, as the drift baseline, the appliance's hash of each tracked
// secret as read back right after the provider sent the value. It is called
// from Create and Update after the read-back. For each attribute in APIHash
// whose extractor reports the field present, it sets Observed[attr]=hash and
// Drift[attr]=false; an attribute whose field is not present is left untouched
// (nothing can be concluded). A nil resp is a no-op. The record is persisted.
//
// Import cannot seed Observed: a baseline needs a send, and import never
// sends. Detection therefore starts after the first apply that sets the value.
func (t Tracker) Observe(ctx context.Context, private PrivateState, resp any) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp == nil || len(t.APIHash) == 0 {
		return diags
	}

	stored, ok, ldiags := Load(ctx, private, t.Key())
	diags.Append(ldiags...)

	if diags.HasError() || !ok {
		return diags
	}

	observed := cloneStrMap(stored.Observed)
	if observed == nil {
		observed = map[string]string{}
	}

	drift := cloneBoolMap(stored.Drift)
	if drift == nil {
		drift = map[string]bool{}
	}

	changed := false

	for attr, extract := range t.APIHash {
		hash, present := extract(resp)
		if !present {
			continue
		}

		observed[attr] = hash
		drift[attr] = false
		changed = true
	}

	if !changed {
		return diags
	}

	stored.Observed = observed
	stored.Drift = drift

	diags.Append(Store(ctx, private, t.Key(), stored)...)

	return diags
}

// StoreObserved is Store followed by Observe, so a resource records both the
// configured-value digests and the appliance-hash baseline in one call after a
// Create or Update read-back. Store keeps its own signature for resources that
// do not do drift detection.
func (t Tracker) StoreObserved(
	ctx context.Context,
	private PrivateState,
	cfg tfsdk.Config,
	resp any,
) diag.Diagnostics {
	diags := t.Store(ctx, private, cfg)
	if diags.HasError() {
		return diags
	}

	diags.Append(t.Observe(ctx, private, resp)...)

	return diags
}

// CheckDrift compares, on a Read, each tracked secret's current appliance hash
// with its Observed baseline. CheckDrift is the single owner of the Drift flag:
// for an attribute with an Observed entry whose extractor reports the field
// present, it sets Drift[attr] = hash != Observed (true OR false); when the
// extractor reports the field not present it clears Drift[attr] (deletes the
// key) -- "no information" must not keep asserting a stale drift. An attribute
// with no baseline (never sent, or imported) is left alone. The record is
// persisted when anything changed. Returns the names now drifted, sorted, so
// the resource can log them.
func (t Tracker) CheckDrift(
	ctx context.Context,
	private PrivateState,
	resp any,
) (drifted []string, diags diag.Diagnostics) {
	if resp == nil || len(t.APIHash) == 0 {
		return nil, diags
	}

	stored, ok, ldiags := Load(ctx, private, t.Key())
	diags.Append(ldiags...)

	if diags.HasError() || !ok || len(stored.Observed) == 0 {
		return nil, diags
	}

	drift := cloneBoolMap(stored.Drift)
	if drift == nil {
		drift = map[string]bool{}
	}

	changed := false

	for attr, extract := range t.APIHash {
		baseline, haveBaseline := stored.Observed[attr]
		if !haveBaseline {
			continue
		}

		hash, present := extract(resp)
		if !present {
			// No information this read: clear any flag rather than leave a
			// stale assertion of drift. A Warn-policy flag that was never
			// cleared otherwise repeats its warning forever; a Resend flag
			// left true after a not-present read mis-plans a re-send.
			if _, had := drift[attr]; had {
				delete(drift, attr)
				changed = true
			}

			continue
		}

		now := hash != baseline
		if drift[attr] != now {
			changed = true
		}

		drift[attr] = now

		if now {
			drifted = append(drifted, attr)
		}
	}

	sort.Strings(drifted)

	if changed {
		stored.Drift = drift
		diags.Append(Store(ctx, private, t.Key(), stored)...)
	}

	return drifted, diags
}

// planDrift turns the drift recorded in stored into a plan decision and the
// diagnostics ModifyPlan reports. For each attribute whose stored secret
// changed out of band and which is still configured (present in current), a
// ResendPolicy plans a re-send (resend true) with an informational note that
// the configured value will be re-sent; a WarnPolicy -- which every immutable
// attribute is forced to, since a re-send would be a replacement -- warns only
// and plans nothing. An attribute that is no longer configured has nothing to
// restore and is ignored. It is pure: it reads stored and current and never
// touches private state.
func (t Tracker) planDrift(stored Digests, current Attrs) (resend bool, diags diag.Diagnostics) {
	for _, attr := range t.driftedConfigured(stored, current) {
		if t.policyFor(attr) == ResendPolicy {
			resend = true

			// The framework has no Info severity for plan diagnostics; a
			// Warning carries the informational wording.
			diags.AddWarning(
				"Write-only value changed on the appliance",
				attr+" was changed on the appliance after Terraform last set it; "+
					"the configured value will be re-sent.")

			continue
		}

		diags.AddWarning(
			"Write-only value changed on the appliance",
			attr+" was changed on the appliance after Terraform last set it. It "+
				"cannot be updated in place; set the attribute to the appliance's "+
				"value, or replace the resource. This warning repeats until the "+
				"configured value matches the appliance or the resource is "+
				"replaced.")
	}

	return resend, diags
}

// NullableStringHash is a HashExtractor building block for an appliance hash
// carried by an SDK NullableString (the shape generated for a string the API
// may omit or render null). It implements the HashExtractor absent-vs-null rule:
//
//   - !IsSet()                -> ("", false)  the API omitted the field: no
//     information, nothing can be concluded.
//   - IsSet(), Get() == nil   -> ("", true)   the field is present and null:
//     the appliance holds no secret, a real comparable value.
//   - IsSet(), Get() != nil   -> (*Get(), true)  the field's hash.
//
// The parameter is an interface, not the concrete SDK type, so this helper has
// no SDK dependency; the generated NullableString satisfies it.
func NullableStringHash(ns interface {
	IsSet() bool
	Get() *string
},
) (hash string, present bool) {
	if !ns.IsSet() {
		return "", false
	}

	if v := ns.Get(); v != nil {
		return *v, true
	}

	return "", true
}
