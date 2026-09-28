// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clearable lets a resource clear an optional attribute on the
// appliance when the user removes it from the configuration.
//
// Terraform cannot express that on its own for an attribute that is both
// Optional and Computed: a null configuration with a non-null state means
// either "the user removed it" or "the user never set it and the server
// computed it", and the provider is given the configuration, the prior state
// and the private state — not the previous configuration. So the resource
// records, in its private state at every apply, which of its clearable
// attributes the configuration set. At the next plan an attribute that was
// configured then, is null now and is non-null in state is planned as a
// clear; one that was never configured is left to the server, as today.
//
// For an attribute that is Optional but not Computed the framework already
// plans the null; what was missing was the Update sending it. The same
// helpers serve both.
//
// A resource declares its clearable attributes with a Mode each, and records
// the reason for every optional attribute that is not clearable; Check holds
// the two lists against the schema so an attribute cannot be left
// unclassified. ModifyPlan calls Plan, Create and Update call Record after a
// successful apply, and Update decides what to send with the Pending helpers.
package clearable

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Mode says what a cleared attribute reads back as.
type Mode int

const (
	// Clear plans the attribute as null: the server forgets the value and
	// returns none. The common case.
	Clear Mode = iota
	// Release plans the attribute as unknown: the server forgets the user's
	// value and computes its own, which becomes the new state — a gateway
	// derived from the CIDR, say. For attributes the server re-derives when
	// told to forget.
	Release
)

// Attribute is one clearable attribute of a resource, by its schema name.
type Attribute struct {
	Name string
	Mode Mode
}

// Attributes are a resource's clearable attributes.
type Attributes []Attribute

// NotClearable records, by name, why an optional attribute is not clearable:
// a boolean with no empty state, a value the API never returns, one of a
// pair where the other is required, and so on.
type NotClearable map[string]string

// PrivateKey is the private-state key under which the record is kept.
const PrivateKey = "clearable.configured"

// recordVersion changes when the record's form changes, so a record written by
// an older provider is treated as absent rather than misread.
const recordVersion = 1

// record is what is kept in private state: the clearable attributes whose
// configuration value was non-null at the last apply.
type record struct {
	Version int      `json:"v"`
	Attrs   []string `json:"attrs"`
}

// privateState is the subset of the framework's private-state type the
// package needs; the request and response Private fields satisfy it.
type privateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// configured returns the recorded set, and false when there is no record —
// a resource created before this mechanism, or imported, or whose record was
// written by an incompatible version.
func configured(ctx context.Context, private privateState) (map[string]bool, bool, diag.Diagnostics) {
	if private == nil {
		return nil, false, nil
	}

	raw, diags := private.GetKey(ctx, PrivateKey)
	if diags.HasError() || len(raw) == 0 {
		return nil, false, diags
	}

	var r record
	if err := json.Unmarshal(raw, &r); err != nil || r.Version != recordVersion {
		return nil, false, diags
	}

	set := make(map[string]bool, len(r.Attrs))
	for _, name := range r.Attrs {
		set[name] = true
	}

	return set, true, diags
}

// Record writes to private state which of the clearable attributes have a
// non-null value in config. Call it from Create and Update once the apply has
// succeeded, with the request's configuration — never from ModifyPlan, where
// a plan that is not applied would leave a record of something never applied.
//
// In Update the request's and the response's private state are the same
// object, so Compute must read the old record before Record overwrites it.
func Record(ctx context.Context, config tfsdk.Config, private privateState, attrs Attributes) diag.Diagnostics {
	var diags diag.Diagnostics

	if private == nil {
		diags.AddError("clearable attributes", "private state is unavailable; the configured attributes could not be recorded")

		return diags
	}

	names := make([]string, 0, len(attrs))

	for _, a := range attrs {
		var v attr.Value

		diags.Append(config.GetAttribute(ctx, path.Root(a.Name), &v)...)
		if diags.HasError() {
			return diags
		}

		if v != nil && !v.IsNull() {
			names = append(names, a.Name)
		}
	}

	sort.Strings(names)

	raw, err := json.Marshal(record{Version: recordVersion, Attrs: names})
	if err != nil {
		diags.AddError("clearable attributes", "encoding the record: "+err.Error())

		return diags
	}

	diags.Append(private.SetKey(ctx, PrivateKey, raw)...)

	return diags
}

// Pending are the attributes being cleared in the current plan or apply: the
// clearable attributes that are null in the configuration and non-null in
// state — and, for an attribute that is also Computed, recorded as configured
// at the last apply. An attribute that is Optional but not Computed needs no
// record: its plan is its configuration, so a null there with a value in
// state is always the user removing it, and Terraform will hold the apply to
// that plan.
type Pending map[string]Attribute

// Compute works out the pending set from the configuration, the state and the
// private state. Plan calls it for ModifyPlan; Update calls it directly with
// its own request.
func Compute(
	ctx context.Context,
	config tfsdk.Config,
	state tfsdk.State,
	private privateState,
	attrs Attributes,
) (Pending, diag.Diagnostics) {
	var diags diag.Diagnostics

	pending := Pending{}

	if state.Raw.IsNull() {
		return pending, diags
	}

	recorded, _, ldiags := configured(ctx, private)
	diags.Append(ldiags...)
	if diags.HasError() {
		return pending, diags
	}

	schemaAttrs := state.Schema.GetAttributes()

	for _, a := range attrs {
		sa, ok := schemaAttrs[a.Name]
		if !ok {
			diags.AddError("clearable attributes", fmt.Sprintf("%q is declared clearable but is not in the schema", a.Name))

			return pending, diags
		}

		if sa.IsComputed() && !recorded[a.Name] {
			continue
		}

		var c, s attr.Value

		diags.Append(config.GetAttribute(ctx, path.Root(a.Name), &c)...)
		diags.Append(state.GetAttribute(ctx, path.Root(a.Name), &s)...)
		if diags.HasError() {
			return pending, diags
		}

		if c == nil || !c.IsNull() {
			continue
		}

		if s == nil || s.IsNull() || s.IsUnknown() {
			continue
		}

		pending[a.Name] = a
	}

	return pending, diags
}

// Plan marks the pending attributes in the plan — null for Clear, unknown for
// Release — so that Terraform shows the change and runs Update. It returns
// the pending set so the resource can decide whether a clear requires
// replacement. Call it from ModifyPlan; a destroy (null plan) and a create
// (null state) have nothing to clear and return an empty set.
func Plan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	attrs Attributes,
) (Pending, diag.Diagnostics) {
	if req.Plan.Raw.IsNull() {
		return Pending{}, nil
	}

	pending, diags := Compute(ctx, req.Config, req.State, req.Private, attrs)
	if diags.HasError() || len(pending) == 0 {
		return pending, diags
	}

	diags.Append(Mark(ctx, &resp.Plan, pending)...)

	return pending, diags
}

// Mark sets each pending attribute in plan to null (Clear) or unknown
// (Release). Plan does this; it is separate so it can be tested without the
// framework's private-state type.
func Mark(ctx context.Context, plan *tfsdk.Plan, pending Pending) diag.Diagnostics {
	var diags diag.Diagnostics

	schemaAttrs := plan.Schema.GetAttributes()

	for _, name := range pending.Names() {
		a, ok := schemaAttrs[name]
		if !ok {
			diags.AddError("clearable attributes", fmt.Sprintf("%q is declared clearable but is not in the schema", name))

			return diags
		}

		typ := a.GetType()
		tfType := typ.TerraformType(ctx)

		var raw tftypes.Value
		if pending[name].Mode == Release {
			raw = tftypes.NewValue(tfType, tftypes.UnknownValue)
		} else {
			raw = tftypes.NewValue(tfType, nil)
		}

		value, err := typ.ValueFromTerraform(ctx, raw)
		if err != nil {
			diags.AddError("clearable attributes", fmt.Sprintf("%s: building the planned value: %s", name, err.Error()))

			return diags
		}

		diags.Append(plan.SetAttribute(ctx, path.Root(name), value)...)
		if diags.HasError() {
			return diags
		}
	}

	return diags
}

// Names returns the pending attribute names, sorted.
func (p Pending) Names() []string {
	names := make([]string, 0, len(p))
	for name := range p {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Clears reports whether the named attribute is being cleared.
func (p Pending) Clears(name string) bool {
	_, ok := p[name]

	return ok
}

// String returns what Update sends for a plain string field of the request,
// which the SDK omits when nil: a pointer to the empty string when the
// attribute is being cleared (Morpheus stores a cleared string as null), a
// pointer to the planned value when it is known, and nil — send nothing —
// otherwise.
func (p Pending) String(name string, plan types.String) *string {
	if p.Clears(name) {
		empty := ""

		return &empty
	}

	if plan.IsNull() || plan.IsUnknown() {
		return nil
	}

	v := plan.ValueString()

	return &v
}

// Nullable sets a nullable request field: to null when the attribute is being
// cleared, to value when known is true, and not at all otherwise. It reports
// whether it set anything. The field is any of the SDK's Nullable types, all
// of which have a Set method taking a pointer.
func Nullable[T any, N interface{ Set(*T) }](p Pending, name string, known bool, value *T, target N) bool {
	switch {
	case p.Clears(name):
		target.Set(nil)

		return true
	case known:
		target.Set(value)

		return true
	default:
		return false
	}
}

// NullableString is Nullable for a string attribute.
func NullableString[N interface{ Set(*string) }](p Pending, name string, plan types.String, target N) bool {
	known := !plan.IsNull() && !plan.IsUnknown()

	var v *string
	if known {
		s := plan.ValueString()
		v = &s
	}

	return Nullable(p, name, known, v, target)
}

// NullableInt64 is Nullable for an integer attribute.
func NullableInt64[N interface{ Set(*int64) }](p Pending, name string, plan types.Int64, target N) bool {
	known := !plan.IsNull() && !plan.IsUnknown()

	var v *int64
	if known {
		i := plan.ValueInt64()
		v = &i
	}

	return Nullable(p, name, known, v, target)
}

// Collection is a planned list or set value: types.List and types.Set both
// satisfy it.
type Collection interface {
	attr.Value
	ElementsAs(ctx context.Context, target any, allowUnhandled bool) diag.Diagnostics
}

// Strings returns what Update sends for a list field from a collection of
// strings: an empty, non-nil slice when the attribute is being cleared, the
// planned elements when known, and false — send nothing — otherwise.
func (p Pending) Strings(ctx context.Context, name string, plan Collection) ([]string, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Clears(name) {
		return []string{}, true, diags
	}

	if plan.IsNull() || plan.IsUnknown() {
		return nil, false, diags
	}

	var out []string

	diags.Append(plan.ElementsAs(ctx, &out, false)...)
	if diags.HasError() {
		return nil, false, diags
	}

	if out == nil {
		out = []string{}
	}

	return out, true, diags
}

// Int64s is Strings for a collection of integers.
func (p Pending) Int64s(ctx context.Context, name string, plan Collection) ([]int64, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if p.Clears(name) {
		return []int64{}, true, diags
	}

	if plan.IsNull() || plan.IsUnknown() {
		return nil, false, diags
	}

	var out []int64

	diags.Append(plan.ElementsAs(ctx, &out, false)...)
	if diags.HasError() {
		return nil, false, diags
	}

	if out == nil {
		out = []int64{}
	}

	return out, true, diags
}

// Computed returns the names of the clearable attributes that are Computed as
// well as Optional. A resource with any of them must record what was
// configured (Record, from Create and Update) and plan the clears itself
// (Plan, from ModifyPlan), since Terraform proposes the prior state for a
// computed attribute the configuration leaves out; a resource with none needs
// neither, because the plan of an Optional-only attribute is its
// configuration. A resource's unit test asserts which case it is in, so the
// wiring cannot drift from the classification.
func Computed(schemaAttrs map[string]schema.Attribute, clearable Attributes) []string {
	var names []string

	for _, a := range clearable {
		if s, ok := schemaAttrs[a.Name]; ok && s.IsComputed() {
			names = append(names, a.Name)
		}
	}

	sort.Strings(names)

	return names
}

// Check holds a resource's classification against its schema and returns
// every problem: an optional attribute (Optional, with or without Computed,
// not write-only, not required) that is neither clearable nor recorded as not
// clearable; a name in both lists; a name in neither the schema; a clearable
// attribute that the user cannot remove (required, computed-only or
// write-only); Release on an attribute that is not computed. A unit test in
// the resource calls it, so the classification stays complete as the schema
// changes.
func Check(schemaAttrs map[string]schema.Attribute, clearable Attributes, not NotClearable) []string {
	var problems []string

	declared := make(map[string]bool, len(clearable))

	for _, a := range clearable {
		if declared[a.Name] {
			problems = append(problems, fmt.Sprintf("%s: declared clearable twice", a.Name))
		}

		declared[a.Name] = true

		if _, ok := not[a.Name]; ok {
			problems = append(problems, fmt.Sprintf("%s: both clearable and not clearable", a.Name))
		}

		s, ok := schemaAttrs[a.Name]
		if !ok {
			problems = append(problems, fmt.Sprintf("%s: declared clearable but not in the schema", a.Name))

			continue
		}

		switch {
		case s.IsWriteOnly():
			problems = append(problems, fmt.Sprintf("%s: write-only attributes are not clearable (they are never in state)", a.Name))
		case s.IsRequired():
			problems = append(problems, fmt.Sprintf("%s: a required attribute cannot be removed", a.Name))
		case !s.IsOptional():
			problems = append(problems, fmt.Sprintf("%s: a computed-only attribute is never configured", a.Name))
		case a.Mode == Release && !s.IsComputed():
			problems = append(problems, fmt.Sprintf("%s: only a computed attribute can be planned unknown (Release)", a.Name))
		}
	}

	for name, reason := range not {
		if _, ok := schemaAttrs[name]; !ok {
			problems = append(problems, fmt.Sprintf("%s: recorded as not clearable but not in the schema", name))
		}

		if reason == "" {
			problems = append(problems, fmt.Sprintf("%s: not clearable without a reason", name))
		}
	}

	for name, s := range schemaAttrs {
		if !s.IsOptional() || s.IsRequired() || s.IsWriteOnly() {
			continue
		}

		if _, ok := not[name]; ok || declared[name] {
			continue
		}

		problems = append(problems, fmt.Sprintf("%s: optional attribute is neither clearable nor recorded as not clearable", name))
	}

	sort.Strings(problems)

	return problems
}
