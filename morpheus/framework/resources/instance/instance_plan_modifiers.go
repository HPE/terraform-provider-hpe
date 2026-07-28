// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// errUnexpectedPathValue is returned when walking an attribute path yields
// something other than a value, which should not happen for a well-formed plan
// or state.
var errUnexpectedPathValue = errors.New("unexpected type at attribute path")

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
// A blanket UseStateForUnknown is not safe here. Update re-reads the instance
// from the API and writes whatever it returns into state, so pinning an
// attribute the API legitimately changes during the update would produce
// "Provider produced inconsistent result after apply" — a hard failure, worse
// than a noisy plan. Two cases matter in practice:
//
//   - connection_info holds the instance's addresses. A resize that cannot be
//     performed hot stops and restarts the instance, and the address can change
//     as a direct result.
//   - Interface and volume identity can change when those collections are
//     themselves reconfigured.
//
// So each group is restored from prior state only when nothing that could
// affect it has changed, and the whole collection is restored at once. Restoring
// per element would be positional, and would pin values from the wrong element
// if interfaces or volumes were added, removed or reordered.

// restoreRule describes one group of computed attributes that can be taken from
// prior state, and the attributes whose modification would invalidate it.
type restoreRule struct {
	// attribute is the attribute restored from prior state.
	attribute string

	// triggers are the attributes that, if changed by the practitioner, mean
	// the restored value can no longer be trusted. An attribute is normally
	// its own trigger.
	triggers []string
}

// restoreRules lists the computed attributes that can be recovered from prior
// state, together with what invalidates them.
var restoreRules = []restoreRule{
	// Interface identity (id, name, primary_interface, ...) is stable unless the
	// interfaces themselves are reconfigured. Note that on appliances from 8.1.2
	// a network change is applied in place rather than forcing replacement, and
	// the API may recreate interfaces, so network_interfaces must gate itself.
	{attribute: "network_interfaces", triggers: []string{"network_interfaces"}},

	// Volume identity is stable unless the volumes themselves are reconfigured.
	{attribute: "volumes", triggers: []string{"volumes"}},

	// Labels only change when the practitioner changes them.
	{attribute: "labels", triggers: []string{"labels"}},

	// Addresses can change whenever the instance is reconfigured or restarted,
	// so anything that drives an infrastructure operation invalidates them.
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

	for _, rule := range restoreRules {
		if anyAttributeChanged(req.Plan.Raw, req.State.Raw, rule.triggers) {
			continue
		}

		restoreAttributeFromState(ctx, req, resp, rule.attribute)
	}
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

// restoreAttributeFromState copies one attribute from prior state into the plan.
func restoreAttributeFromState(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	attribute string,
) {
	p := path.Root(attribute)

	// These attributes are all collections; read them as their concrete types so
	// the framework performs the conversion.
	switch attribute {
	case "labels":
		var value types.Set
		if diags := req.State.GetAttribute(ctx, p, &value); diags.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, p, value)...)
	default:
		var value types.List
		if diags := req.State.GetAttribute(ctx, p, &value); diags.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, p, value)...)
	}
}
