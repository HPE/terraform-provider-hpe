// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"errors"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The plugin framework marks every computed attribute that is null in the
// configuration as unknown whenever the planned state differs from the prior
// state (MarkComputedNilsAsUnknown in server_planresourcechange.go). That is
// deliberate — the framework's own comment notes that later plan modifier passes
// are expected to put known values back — but it means a change to one attribute
// makes every unrelated computed attribute show as "(known after apply)".
//
// On this resource that produces a large, unreviewable diff for a small edit:
// changing service_plan_options.max_memory also churns connection_info, labels,
// and every computed field of network_interfaces and volumes.
//
// Restoring values has to be done narrowly. Update re-reads the instance from
// the API and writes what it returns into state, so a planned value that the
// read path then overrides produces "Provider produced inconsistent result
// after apply" — a hard failure, worse than a noisy plan. Two rules keep this
// safe:
//
//  1. Only attributes whose value is stable and API-assigned are restored:
//     interface and volume identity. Attributes such as volume storage_profile
//     are excluded, because the post-apply read prefers the API value whenever
//     the planned value is null or unknown (see instance_read.go), so pinning
//     them from prior state can contradict what apply returns.
//  2. A value is only ever filled in where the plan says unknown, and only from
//     a prior-state value that is itself known and non-null. A null is never
//     pinned, so restoring cannot turn "unknown, resolve later" into an
//     assertion the read path will contradict.
//
// Each group is additionally gated on its own triggers: if the practitioner
// changed the collection, nothing in it is restored, because the API may assign
// new identities and a positional restore could pin the previous element's
// values onto a different one.

// errUnexpectedPathValue is returned when walking an attribute path yields
// something other than a value, which should not happen for a well-formed plan
// or state.
var errUnexpectedPathValue = errors.New("unexpected type at attribute path")

// restoreRule describes attributes that can be taken from prior state, and the
// attributes whose modification would invalidate them.
type restoreRule struct {
	// attribute is the root attribute this rule applies to.
	attribute string

	// fields, when non-empty, restricts the rule to those nested attributes of
	// each element of the collection. When empty the root attribute value itself
	// is restored.
	fields []string

	// triggers are the attributes that, if changed by the practitioner, mean the
	// restored values can no longer be trusted. An attribute is normally its own
	// trigger.
	triggers []string
}

// restoreRules lists what may be recovered from prior state.
//
// Only identity attributes are listed for the nested collections. Attributes
// that the post-apply read sources from the API — volume storage_profile and
// controller_mount_point, interface ip_address — are deliberately absent.
var restoreRules = []restoreRule{
	// Interface identity is stable unless the interfaces are reconfigured. Note
	// that from 8.1.2 a network change is applied in place rather than forcing
	// replacement, and the API may recreate interfaces, so network_interfaces
	// gates itself.
	{
		attribute: "network_interfaces",
		fields:    []string{"id", "name", "primary_interface"},
		triggers:  []string{"network_interfaces"},
	},

	// Volume identity is stable unless the volumes are reconfigured.
	{
		attribute: "volumes",
		fields:    []string{"id"},
		triggers:  []string{"volumes"},
	},

	// Labels are a flat collection with no nested attributes, and only change
	// when the practitioner changes them.
	{attribute: "labels", triggers: []string{"labels"}},

	// Addresses can change whenever the instance is reconfigured or restarted,
	// so anything driving an infrastructure operation invalidates them.
	{
		attribute: "connection_info",
		triggers: []string{
			"service_plan_options",
			"network_interfaces",
			"volumes",
			"config",
			"plan_id",
			"layout_id",
		},
	},
}

// restoreUnchangedComputedAttributes puts prior state values back into the plan
// for computed attributes whose triggering configuration has not been modified,
// so that an unrelated edit does not show them as "(known after apply)".
func restoreUnchangedComputedAttributes(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Only applies to updates: on create there is no prior state to restore
	// from, and on destroy there is no plan to modify.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	applicable := make([]restoreRule, 0, len(restoreRules))

	for _, rule := range restoreRules {
		if anyAttributeChanged(req.Plan.Raw, req.State.Raw, rule.triggers) {
			continue
		}

		applicable = append(applicable, rule)
	}

	if len(applicable) == 0 {
		return
	}

	restored, err := restoreFromState(resp.Plan.Raw, req.State.Raw, applicable)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error restoring planned values",
			"An unexpected error occurred while restoring unchanged computed "+
				"attributes into the plan. This is always a problem with the "+
				"provider. Please report the following to the provider "+
				"developer:\n\n"+err.Error(),
		)

		return
	}

	resp.Plan.Raw = restored
}

// restoreFromState fills unknown values in plan from prior state, for the paths
// covered by the supplied rules.
//
// Only unknown planned values are filled, and only from state values that are
// known, non-null and of the same type. Anything else is left as planned, so a
// value the read path will supply is never pre-empted.
func restoreFromState(
	plan, state tftypes.Value,
	rules []restoreRule,
) (tftypes.Value, error) {
	return tftypes.Transform(
		plan,
		func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if v.IsKnown() || !matchesAnyRestoreRule(p, rules) {
				return v, nil
			}

			stateValue, err := valueAtPath(state, p)
			if err != nil {
				// No counterpart in prior state — for example a newly added
				// element. Leave it unknown.
				return v, nil
			}

			if !stateValue.IsKnown() || stateValue.IsNull() {
				return v, nil
			}

			if !stateValue.Type().Equal(v.Type()) {
				return v, nil
			}

			return stateValue, nil
		},
	)
}

// matchesAnyRestoreRule reports whether path is covered by one of the rules.
func matchesAnyRestoreRule(p *tftypes.AttributePath, rules []restoreRule) bool {
	return slices.ContainsFunc(rules, func(r restoreRule) bool {
		return matchesRestoreRule(p, r)
	})
}

// matchesRestoreRule reports whether path is covered by a single rule.
//
// For a rule without fields the path must be the root attribute itself. For a
// rule with fields the path must address one of those fields within an element
// of the collection, e.g. volumes[0].id.
func matchesRestoreRule(p *tftypes.AttributePath, r restoreRule) bool {
	steps := p.Steps()
	if len(steps) == 0 {
		return false
	}

	root, ok := steps[0].(tftypes.AttributeName)
	if !ok || string(root) != r.attribute {
		return false
	}

	if len(r.fields) == 0 {
		return len(steps) == 1
	}

	// collection[index].field
	if len(steps) != 3 {
		return false
	}

	if _, ok := steps[1].(tftypes.ElementKeyInt); !ok {
		return false
	}

	leaf, ok := steps[2].(tftypes.AttributeName)
	if !ok {
		return false
	}

	return slices.Contains(r.fields, string(leaf))
}

// anyAttributeChanged reports whether the practitioner changed any of the named
// attributes.
func anyAttributeChanged(plan, state tftypes.Value, attributes []string) bool {
	for _, name := range attributes {
		if attributeChanged(plan, state, name) {
			return true
		}
	}

	return false
}

// attributeChanged reports whether the named attribute differs between the plan
// and prior state, ignoring values the framework marked unknown.
//
// Unknown values carry no intent: they are what the framework substituted for
// computed attributes, not something the practitioner asked for. They are
// therefore filled from prior state before comparing, so that only values the
// practitioner actually set are considered.
func attributeChanged(plan, state tftypes.Value, attribute string) bool {
	p := tftypes.NewAttributePath().WithAttributeName(attribute)

	planValue, err := valueAtPath(plan, p)
	if err != nil {
		return true
	}

	stateValue, err := valueAtPath(state, p)
	if err != nil {
		return true
	}

	filled, err := fillUnknownsFromState(planValue, stateValue)
	if err != nil {
		return true
	}

	return !filled.Equal(stateValue)
}

// fillUnknownsFromState replaces every unknown value in plan with the value at
// the same path in state. Values with no counterpart in state — for example a
// newly added list element — are left unknown, so the comparison still reports a
// difference.
func fillUnknownsFromState(plan, state tftypes.Value) (tftypes.Value, error) {
	return tftypes.Transform(
		plan,
		func(p *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
			if v.IsKnown() {
				return v, nil
			}

			stateValue, err := valueAtPath(state, p)
			if err != nil {
				return v, nil
			}

			if !stateValue.Type().Equal(v.Type()) {
				return v, nil
			}

			return stateValue, nil
		},
	)
}

// valueAtPath returns the value at path within val.
func valueAtPath(val tftypes.Value, p *tftypes.AttributePath) (tftypes.Value, error) {
	// A path relative to the value itself resolves to the value.
	if len(p.Steps()) == 0 {
		return val, nil
	}

	raw, _, err := tftypes.WalkAttributePath(val, p)
	if err != nil {
		return tftypes.Value{}, err
	}

	value, ok := raw.(tftypes.Value)
	if !ok {
		return tftypes.Value{}, errUnexpectedPathValue
	}

	return value, nil
}
