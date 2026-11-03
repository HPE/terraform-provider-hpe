// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clearable

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Before the plan modifiers run, the framework marks every Computed attribute
// that is null in the configuration and has no default as unknown — "known
// after apply" — so the provider can report whatever the update makes of it.
// It does that only when the proposed plan already differs from the prior
// state. When the clear is the only change it does not: Terraform proposes
// the prior value for a removed Optional+Computed attribute, the proposed
// plan equals the state, and the pass is skipped. Marking the pending
// attributes then leaves every other computed attribute at its prior value,
// and one the update changes on the appliance — a status, a last-updated
// timestamp — fails the apply with "Provider produced inconsistent result".
//
// UnknownComputed runs the equivalent of that pass: it marks the computed
// attributes that the framework would have, then runs each marked
// attribute's plan modifiers so that UseStateForUnknown and the like keep
// their say, exactly as they do on any other update. It is idempotent —
// attributes already unknown are left alone, and the modifiers' outcome is
// the same the second time — so Plan calls it whenever it has something
// pending, without having to know whether the framework's pass ran.
//
// Like the framework's pass, it leaves the children of a set element alone
// in practice: the element is keyed by its value, and the configuration's
// element — its computed children null — does not match the plan's, so the
// configuration lookup fails. Where it does match, a child that carries plan
// modifiers is still left alone, because a set element has no stable path
// to run them by. A path the schema does not resolve (a block, the inside of
// an atomic attribute) is left alone rather than reported.
func UnknownComputed(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if resp.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return diags
	}

	type marked struct {
		tfPath *tftypes.AttributePath
		attr   any
	}

	var toModify []marked

	transformed, err := tftypes.Transform(resp.Plan.Raw, func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if len(p.Steps()) == 0 || !v.IsKnown() {
			return v, nil
		}

		a, err := resp.Plan.Schema.AttributeAtTerraformPath(ctx, p)
		if err != nil {
			// A block, the inside of an atomic or dynamic attribute: the
			// framework leaves these alone too.
			return v, nil
		}

		if !a.IsComputed() || hasDefault(a) {
			return v, nil
		}

		configValue, _, err := tftypes.WalkAttributePath(req.Config.Raw, p)
		if err != nil {
			return v, nil
		}

		cv, ok := configValue.(tftypes.Value)
		if !ok || !cv.IsNull() {
			return v, nil
		}

		if insideSetElement(p) && hasPlanModifiers(a) {
			return v, nil
		}

		toModify = append(toModify, marked{tfPath: p, attr: a})

		return tftypes.NewValue(v.Type(), tftypes.UnknownValue), nil
	})
	if err != nil {
		diags.AddError("clearable attributes", "marking computed attributes unknown: "+err.Error())

		return diags
	}

	resp.Plan.Raw = transformed

	for _, m := range toModify {
		if !hasPlanModifiers(m.attr) {
			continue
		}

		p, ok := toPath(m.tfPath)
		if !ok {
			continue
		}

		diags.Append(runPlanModifiers(ctx, p, m.attr, req, resp)...)
		if diags.HasError() {
			return diags
		}
	}

	return diags
}

// hasDefault reports whether the schema attribute carries a Default; the
// framework never marks a defaulted attribute unknown. Every attribute kind
// exposes it as a field named Default, so reflection saves sixteen cases.
func hasDefault(a any) bool {
	f := reflect.ValueOf(a).FieldByName("Default")

	return f.IsValid() && !f.IsNil()
}

// hasPlanModifiers reports whether the schema attribute carries any plan
// modifiers, by the same means.
func hasPlanModifiers(a any) bool {
	f := reflect.ValueOf(a).FieldByName("PlanModifiers")

	return f.IsValid() && f.Len() > 0
}

// insideSetElement reports whether the path passes through a set element.
func insideSetElement(p *tftypes.AttributePath) bool {
	for _, step := range p.Steps() {
		if _, ok := step.(tftypes.ElementKeyValue); ok {
			return true
		}
	}

	return false
}

// toPath converts a terraform-plugin-go path to a framework one. Set
// elements are keyed by their value and are not converted; the caller has
// already left those alone.
func toPath(tp *tftypes.AttributePath) (path.Path, bool) {
	p := path.Empty()

	for _, step := range tp.Steps() {
		switch s := step.(type) {
		case tftypes.AttributeName:
			p = p.AtName(string(s))
		case tftypes.ElementKeyInt:
			p = p.AtListIndex(int(s))
		case tftypes.ElementKeyString:
			p = p.AtMapKey(string(s))
		default:
			return path.Empty(), false
		}
	}

	return p, true
}

// runPlanModifiers runs the attribute's plan modifiers at p, with the plan
// value now unknown, and writes their result back — the framework's
// attribute pass for one attribute. Replacement requests and diagnostics go
// to the ModifyPlan response as they would from the framework.
//
//nolint:gocyclo,funlen,maintidx // one case per attribute kind; each is the same shape
func runPlanModifiers(
	ctx context.Context,
	p path.Path,
	a any,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) diag.Diagnostics {
	var diags diag.Diagnostics

	replace := func(r bool) {
		if r && !resp.RequiresReplace.Contains(p) {
			resp.RequiresReplace.Append(p)
		}
	}

	switch sa := a.(type) {
	case schema.StringAttribute:
		r := planmodifier.StringRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.StringResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyString(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.Int64Attribute:
		r := planmodifier.Int64Request{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.Int64Response{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyInt64(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.Int32Attribute:
		r := planmodifier.Int32Request{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.Int32Response{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyInt32(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.BoolAttribute:
		r := planmodifier.BoolRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.BoolResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyBool(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.Float64Attribute:
		r := planmodifier.Float64Request{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.Float64Response{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyFloat64(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.Float32Attribute:
		r := planmodifier.Float32Request{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.Float32Response{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyFloat32(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.NumberAttribute:
		r := planmodifier.NumberRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.NumberResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyNumber(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.ListAttribute, schema.ListNestedAttribute:
		var mods []planmodifier.List
		if l, ok := sa.(schema.ListAttribute); ok {
			mods = l.PlanModifiers
		} else if l, ok := sa.(schema.ListNestedAttribute); ok {
			mods = l.PlanModifiers
		}

		r := planmodifier.ListRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range mods {
			mr := planmodifier.ListResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyList(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.SetAttribute, schema.SetNestedAttribute:
		var mods []planmodifier.Set
		if s, ok := sa.(schema.SetAttribute); ok {
			mods = s.PlanModifiers
		} else if s, ok := sa.(schema.SetNestedAttribute); ok {
			mods = s.PlanModifiers
		}

		r := planmodifier.SetRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range mods {
			mr := planmodifier.SetResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifySet(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.MapAttribute, schema.MapNestedAttribute:
		var mods []planmodifier.Map
		if m, ok := sa.(schema.MapAttribute); ok {
			mods = m.PlanModifiers
		} else if m, ok := sa.(schema.MapNestedAttribute); ok {
			mods = m.PlanModifiers
		}

		r := planmodifier.MapRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range mods {
			mr := planmodifier.MapResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyMap(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.ObjectAttribute, schema.SingleNestedAttribute:
		var mods []planmodifier.Object
		if o, ok := sa.(schema.ObjectAttribute); ok {
			mods = o.PlanModifiers
		} else if o, ok := sa.(schema.SingleNestedAttribute); ok {
			mods = o.PlanModifiers
		}

		r := planmodifier.ObjectRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range mods {
			mr := planmodifier.ObjectResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyObject(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	case schema.DynamicAttribute:
		r := planmodifier.DynamicRequest{
			Path: p, PathExpression: p.Expression(), Private: resp.Private,
			Config: req.Config, Plan: resp.Plan, State: req.State,
		}
		diags.Append(get(ctx, p, req.Config, req.State, resp.Plan, &r.ConfigValue, &r.StateValue, &r.PlanValue)...)
		if diags.HasError() {
			return diags
		}

		for _, m := range sa.PlanModifiers {
			mr := planmodifier.DynamicResponse{PlanValue: r.PlanValue, Private: r.Private}
			m.PlanModifyDynamic(ctx, r, &mr)
			diags.Append(mr.Diagnostics...)
			replace(mr.RequiresReplace)
			r.PlanValue, r.Private = mr.PlanValue, mr.Private
		}

		diags.Append(resp.Plan.SetAttribute(ctx, p, r.PlanValue)...)
	}

	return diags
}

// get reads the configuration, state and plan values at p into the typed
// targets a plan modifier request wants.
func get(
	ctx context.Context,
	p path.Path,
	config tfsdk.Config,
	state tfsdk.State,
	plan tfsdk.Plan,
	configValue, stateValue, planValue any,
) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(config.GetAttribute(ctx, p, configValue)...)
	diags.Append(state.GetAttribute(ctx, p, stateValue)...)
	diags.Append(plan.GetAttribute(ctx, p, planValue)...)

	return diags
}
