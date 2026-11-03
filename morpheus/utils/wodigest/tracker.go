// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RevisionAttribute is the computed attribute every resource with tracked
// write-only attributes has: the time of the last apply that sent a write-only
// value. ModifyPlan marks it unknown to plan an update.
const RevisionAttribute = "write_only_revision"

// Tracker is the resource-side glue for a resource whose write-only
// attributes are independent optional string values (passwords, keys): it
// decides in ModifyPlan, records in Create/Update/ImportState. A resource
// declares one as a package-level value:
//
//	var wo = wodigest.Tracker{
//	    Resource: "backup_host",
//	    Attributes: map[string]string{"ssh_password_wo": "ssh_password_wo_version"},
//	}
//
// and calls wo.ModifyPlan, wo.Store, wo.Import and wo.Resend from the
// corresponding methods. Attribute names are top-level or dotted paths into
// single nested attributes.
type Tracker struct {
	// Resource names the resource in messages and in the private-state key.
	Resource string
	// Attributes maps each write-only attribute to its deprecated *_version
	// marker ("" when it has none). A change to a marker is honoured as a
	// change for the whole 2.2 line.
	Attributes map[string]string
	// Immutable lists the attributes whose change requires replacement.
	Immutable []string
	// APIHash maps an attribute to an extractor of the appliance's hash of the
	// stored secret from the resource's read response. Only attributes listed
	// here take part in appliance-side drift detection; every other attribute
	// keeps today's config-side behaviour. The hash is compared for equality
	// across reads only; its algorithm is never assumed.
	APIHash map[string]HashExtractor
	// OnDrift is what a detected out-of-band change plans. An attribute in
	// Immutable always behaves as Warn regardless of this, since re-sending an
	// immutable value would be a replacement.
	OnDrift DriftPolicy
}

// HashExtractor returns the appliance's hash for one attribute from the
// resource's read response. present=false means the response has no such field
// at all (older appliance, a different object shape): nothing can be concluded.
// present=true with hash=="" means the field is there and null or empty: the
// appliance holds no secret -- a real, comparable value.
//
// Binding rule for implementers -- distinguish "the API omitted the field"
// from "the API rendered the field as null":
//
//   - field omitted by the API (no information)      -> ("", false)
//   - field rendered null/empty (no secret stored)   -> ("", true)
//   - field rendered with a hash                      -> (hash, true)
//
// An extractor that cannot tell omitted from null must prefer present=false:
// asserting "no secret" when the field was merely omitted would record a false
// drift. NullableStringHash implements the rule for a nullable field; a plain
// *string cannot tell omitted from null, so an extractor over one must read a
// nil pointer as present=false.
type HashExtractor func(resp any) (hash string, present bool)

// DriftPolicy is what ModifyPlan does when an attribute's stored secret was
// changed on the appliance. The value names carry a Policy suffix because the
// package already exports a Resend function (whether Update should re-send the
// values), which must keep its name for its many call sites.
type DriftPolicy int

const (
	// ResendPolicy plans one update that re-sends the configured value. The
	// default (zero value).
	ResendPolicy DriftPolicy = iota
	// WarnPolicy reports the drift as a warning and plans nothing.
	WarnPolicy
)

func cloneStrMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}

	return out
}

func cloneBoolMap(m map[string]bool) map[string]bool {
	if m == nil {
		return nil
	}

	out := make(map[string]bool, len(m))
	for k, v := range m {
		out[k] = v
	}

	return out
}

// Key is the private-state key.
func (t Tracker) Key() string { return t.Resource + ".wo_digests" }

// Names lists the tracked attributes, for ParseImportID.
func (t Tracker) Names() []string {
	names := make([]string, 0, len(t.Attributes))
	for k := range t.Attributes {
		names = append(names, k)
	}

	return names
}

// getter is what tfsdk.Config, tfsdk.Plan and tfsdk.State have in common.
type getter interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

func attrPath(name string) path.Path {
	parts := strings.Split(name, ".")
	p := path.Root(parts[0])

	for _, s := range parts[1:] {
		p = p.AtName(s)
	}

	return p
}

// configured returns the per-attribute digests of the tracked attributes as
// set in cfg (unset ones absent). A string is hashed as-is; any other value
// (a number, a dynamic map, a nested object) through its Terraform JSON form,
// which is what an importer reproduces with jsonencode().
func (t Tracker) configured(ctx context.Context, cfg getter) (Attrs, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrs := Attrs{}

	for name := range t.Attributes {
		var v attr.Value

		d := cfg.GetAttribute(ctx, attrPath(name), &v)
		diags.Append(d...)

		if d.HasError() || v == nil || v.IsNull() || v.IsUnknown() {
			continue
		}

		h, err := ValueOf(ctx, v)
		if err != nil {
			diags.AddError(t.Resource+": write-only digests", name+": "+err.Error())

			continue
		}

		attrs[name] = h
	}

	return attrs, diags
}

// markerChanged reports whether any marker differs between plan and state.
func (t Tracker) markerChanged(ctx context.Context, plan, state getter) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	for _, marker := range t.Attributes {
		if marker == "" {
			continue
		}

		var p, s attr.Value

		diags.Append(plan.GetAttribute(ctx, attrPath(marker), &p)...)
		diags.Append(state.GetAttribute(ctx, attrPath(marker), &s)...)

		if diags.HasError() {
			return false, diags
		}

		if !p.Equal(s) {
			return true, diags
		}
	}

	return false, diags
}

// ModifyPlan decides what to do about the write-only values and plans it:
// nothing, an in-place update (write_only_revision unknown, plus any extra
// computed attributes the caller names, which the re-send also changes), or a
// replacement. It does nothing on create or destroy.
func (t Tracker) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	alsoUnknown ...path.Path,
) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	current, diags := t.configured(ctx, req.Config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	stored, ok, diags := Load(ctx, req.Private, t.Key())
	resp.Diagnostics.Append(diags...)

	marker, diags := t.markerChanged(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	decision := DecideConfigured(stored, ok, current, t.Immutable)

	tflog.Debug(ctx, t.Resource+": write-only change detection", map[string]any{
		"digests_stored": ok, "sentinel": stored.Sentinel,
		"decision": decision.String(), "marker_changed": marker,
	})

	if ok && stored.Sentinel == ImportSentinel && !marker && len(current) > 0 {
		// The sentinel is replaced by real digests on the next apply of any
		// kind, so the warning is only for plans that would not apply.
		if !req.Plan.Raw.Equal(req.State.Raw) {
			return
		}

		names := t.Names()
		sort.Strings(names)

		how := "the configured values are sent on the next update of any other attribute" + t.markerHint()
		if len(t.Immutable) > 0 {
			how = "the configured values can only take effect through a replacement of the resource; " +
				"any update of another attribute records them as applied"
		}

		resp.Diagnostics.AddWarning(
			"Write-only values of an imported "+t.Resource+" are taken as already applied",
			"This resource was imported without the digests of its write-only values, so a change to "+
				strings.Join(names, ", ")+" is not detected until the resource is next applied. To have "+
				"changes detected from the start, import with the values' sha256 digests in the import ID: "+
				"\"<id>#"+names[0]+"=<sha256 hex>\". Otherwise "+how+".")

		return
	}

	// Appliance-side drift (CheckDrift recorded it on the last Read). planDrift
	// decides, per attribute, whether to re-send the configured value and emits
	// the diagnostics; ModifyPlan never writes private state, so the recorded
	// Drift flags are left as they are and the next successful apply re-baselines
	// them through StoreObserved.
	driftResend, driftDiags := t.planDrift(stored, current)
	resp.Diagnostics.Append(driftDiags...)

	if decision == Unchanged && !marker && !driftResend {
		return
	}

	// A change to an immutable value -- or to its marker, which for the 2.2
	// line still means "re-send", and re-sending an immutable value is a
	// replacement -- plans a replacement of the resource.
	if decision == Replace || (marker && t.immutableMarkerChanged(ctx, req.Plan, req.State)) {
		for _, name := range t.Immutable {
			resp.RequiresReplace = append(resp.RequiresReplace, attrPath(name))
		}

		return
	}

	// A Seed on a resource whose only write-only values are immutable has
	// nothing to send: record the digests through the update hook as usual
	// (Update stores them) -- Create-only values are simply not re-sent.

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(RevisionAttribute), types.Int64Unknown())...)

	for _, p := range alsoUnknown {
		var v attr.Value

		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, p, &v)...)

		if resp.Diagnostics.HasError() {
			return
		}

		u, err := unknownOf(ctx, v)
		if err != nil {
			resp.Diagnostics.AddError(t.Resource+": write-only digests", p.String()+": "+err.Error())

			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, p, u)...)
	}
}

// driftedConfigured returns, sorted, the tracked attributes that CheckDrift
// recorded as drifted and that are still configured (present in current). An
// attribute with no extractor, no recorded drift, or no configured value is
// excluded: there is nothing to restore. Read-only; it does not persist.
func (t Tracker) driftedConfigured(stored Digests, current Attrs) []string {
	if len(stored.Drift) == 0 || len(t.APIHash) == 0 {
		return nil
	}

	var names []string

	for attr := range t.APIHash {
		if !stored.Drift[attr] {
			continue
		}

		if _, configured := current[attr]; !configured {
			continue
		}

		names = append(names, attr)
	}

	sort.Strings(names)

	return names
}

// immutableMarkerChanged reports whether the marker of an immutable attribute
// was set or changed to a new value between state and plan. A marker that is
// being removed from the configuration (planned null) does not count: removing
// the deprecated attribute must never replace the resource.
func (t Tracker) immutableMarkerChanged(ctx context.Context, plan, state getter) bool {
	for _, name := range t.Immutable {
		marker := t.Attributes[name]
		if marker == "" {
			continue
		}

		var p, s attr.Value

		if plan.GetAttribute(ctx, attrPath(marker), &p).HasError() || state.GetAttribute(ctx, attrPath(marker), &s).HasError() {
			continue
		}

		if !p.Equal(s) && !p.IsNull() {
			return true
		}
	}

	return false
}

func (t Tracker) markerHint() string {
	for _, marker := range t.Attributes {
		if marker != "" {
			return ", or set " + marker
		}
	}

	return ""
}

// Resend reports whether Update should send the write-only values: ModifyPlan
// marked write_only_revision unknown. It learns about appliance-side drift
// through the same signal -- ModifyPlan marks the revision unknown for a
// ResendPolicy attribute whose stored secret drifted -- so this needs no
// change and no new channel for the drift case.
func Resend(ctx context.Context, plan tfsdk.Plan) bool {
	var v types.Int64
	if plan.GetAttribute(ctx, path.Root(RevisionAttribute), &v).HasError() {
		return false
	}

	return v.IsUnknown()
}

// Store records the digests of the configured write-only values, merged over
// what was stored before (an unset attribute is not managed, so its last
// digest is kept). Called by Create and, unconditionally, by Update, so that
// neither an import sentinel nor a missing digest survives an apply.
func (t Tracker) Store(ctx context.Context, private PrivateState, cfg tfsdk.Config) diag.Diagnostics {
	current, diags := t.configured(ctx, cfg)
	if diags.HasError() {
		return diags
	}

	stored, ok, ldiags := Load(ctx, private, t.Key())
	diags.Append(ldiags...)

	var observed map[string]string

	var drift map[string]bool

	if ok {
		// Carry the drift baseline forward across every Store: a Store that
		// does not re-send a secret must not discard what Observe/CheckDrift
		// recorded, whatever kind of record was stored before (including an
		// import sentinel, whose per-read Drift is still meaningful).
		// StoreObserved re-baselines the affected attributes itself.
		observed = cloneStrMap(stored.Observed)
		drift = cloneBoolMap(stored.Drift)
	}

	if ok && stored.Sentinel == "" {
		merged := Attrs{}
		for k, v := range stored.Attrs {
			merged[k] = v
		}

		for k, v := range current {
			merged[k] = v
		}

		current = merged
	}

	d := Compute(current, t.Immutable)
	d.Observed = observed
	d.Drift = drift

	diags.Append(Store(ctx, private, t.Key(), d)...)

	return diags
}

// Import splits the import ID into the resource's own ID and the digests it
// carries, records those (or the import sentinel), and gives the imported
// resource a write_only_revision. The caller parses the returned ID as it
// always did.
func (t Tracker) Import(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) (id string, ok bool) {
	id, attrs, err := ParseImportID(req.ID, t.Names())
	if err != nil {
		resp.Diagnostics.AddError("import "+t.Resource, err.Error())

		return "", false
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(RevisionAttribute), NewRevision())...)

	if attrs == nil {
		resp.Diagnostics.Append(StoreImportSentinel(ctx, resp.Private, t.Key())...)
	} else {
		resp.Diagnostics.Append(Store(ctx, resp.Private, t.Key(), Compute(attrs, t.Immutable))...)
	}

	return id, !resp.Diagnostics.HasError()
}

// unknownOf returns the unknown value of v's type. Every framework value type
// carries its own type, so the unknown is built from it; a value whose type
// cannot be recovered is a programming error, reported so a test catches it
// rather than a plan/apply inconsistency at runtime.
func unknownOf(ctx context.Context, v attr.Value) (attr.Value, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot mark a nil value unknown")
	}

	typ := v.Type(ctx)

	tfType := typ.TerraformType(ctx)

	unknown, err := typ.ValueFromTerraform(ctx, tftypes.NewValue(tfType, tftypes.UnknownValue))
	if err != nil {
		return nil, fmt.Errorf("cannot build the unknown value of %s: %w", typ, err)
	}

	return unknown, nil
}
