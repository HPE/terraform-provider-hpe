// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// evarsDiff is the set of environment-variable changes needed to move an
// instance from its current variables to the planned ones.
//
// The instance update endpoint does not carry environment variables at all, so
// they are applied through the dedicated per-instance environment-variable
// operations instead: one add-or-update call carrying every variable whose
// value differs or is new (the API matches on name and replaces the value),
// and one delete per variable that is no longer planned. MORPH-612.
type evarsDiff struct {
	// upsert: variables that are new, or whose value alone changed. The
	// add-or-update call handles both (it matches by name and replaces the value).
	upsert []sdk.AddEnvVariablesRequestEnvsInner
	// recreate: existing variables whose export or masked changed. The
	// add-or-update call updates ONLY the value of an existing variable and
	// ignores the options, so these must be deleted and added again.
	recreate []sdk.AddEnvVariablesRequestEnvsInner
	// remove: variables no longer planned.
	remove []string
}

func (d evarsDiff) empty() bool {
	return len(d.upsert) == 0 && len(d.recreate) == 0 && len(d.remove) == 0
}

// evarSettings is everything settable on a variable besides its name. Export
// and masked carry the appliance defaults (true / false) when the practitioner
// leaves them unset, which is also what the schema defaults them to.
type evarSettings struct {
	value  string
	export bool
	masked bool
}

// evarsByName flattens an evars set into name -> settings. An unknown or null
// set contributes nothing; the caller decides what "nothing" means.
func evarsByName(ctx context.Context, set basetypes.SetValue) (map[string]evarSettings, diag.Diagnostics) {
	out := map[string]evarSettings{}
	if set.IsNull() || set.IsUnknown() {
		return out, nil
	}

	var elems []EvarsValue
	diags := set.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil, diags
	}

	for _, e := range elems {
		out[e.Name.ValueString()] = evarSettings{
			value:  e.Value.ValueString(),
			export: boolOrDefault(e.Export, true),
			masked: boolOrDefault(e.Masked, false),
		}
	}

	return out, nil
}

// boolOrDefault returns the bool's value, or def when it is null or unknown.
func boolOrDefault(b basetypes.BoolValue, def bool) bool {
	if b.IsNull() || b.IsUnknown() {
		return def
	}

	return b.ValueBool()
}

// diffEvars computes the changes needed to move from the current (state)
// variables to the planned ones.
//
// A planned variable that is absent from state, or present with a different
// value, export or masked setting, is upserted (the add-or-update call replaces
// all three by name). A state variable absent from the plan is removed --
// including when the plan omits `evars` entirely, which is how a practitioner
// removes every variable. Output is sorted by name so the calls are
// deterministic and the diff is easy to assert in tests.
func diffEvars(planned, current map[string]evarSettings) evarsDiff {
	var d evarsDiff

	for name, want := range planned {
		cur, exists := current[name]
		if exists && cur == want {
			continue
		}
		entry := evarRequest(name, want)
		switch {
		case exists && (cur.export != want.export || cur.masked != want.masked):
			// The options on an existing variable cannot be changed in place.
			d.recreate = append(d.recreate, entry)
		default:
			// New variable, or an existing one whose value alone changed.
			d.upsert = append(d.upsert, entry)
		}
	}
	for name := range current {
		if _, ok := planned[name]; !ok {
			d.remove = append(d.remove, name)
		}
	}

	byName := func(list []sdk.AddEnvVariablesRequestEnvsInner) func(i, j int) bool {
		return func(i, j int) bool { return list[i].Name < list[j].Name }
	}
	sort.Slice(d.upsert, byName(d.upsert))
	sort.Slice(d.recreate, byName(d.recreate))
	sort.Strings(d.remove)

	return d
}

// evarRequest builds the add-or-update entry for one variable, always carrying
// export and masked explicitly so the API never falls back to its own defaults
// for a variable the practitioner is setting.
func evarRequest(name string, s evarSettings) sdk.AddEnvVariablesRequestEnvsInner {
	export, masked := s.export, s.masked

	return sdk.AddEnvVariablesRequestEnvsInner{
		Name:   name,
		Value:  s.value,
		Export: &export,
		Masked: &masked,
	}
}

// applyEvarsUpdate brings the instance's user-defined environment variables in
// line with the plan. It is a no-op when the plan is unknown (nothing to apply
// yet) or when nothing differs.
func applyEvarsUpdate(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	plan, state InstanceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if plan.Evars.IsUnknown() {
		return diags
	}

	planned, d := evarsByName(ctx, plan.Evars)
	diags.Append(d...)
	current, d := evarsByName(ctx, state.Evars)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	changes := diffEvars(planned, current)
	if changes.empty() {
		return diags
	}

	// Variables whose options changed must be removed first: the add-or-update
	// call would otherwise only replace their value.
	for _, entry := range changes.recreate {
		if err := deleteEvar(ctx, client, id, entry.Name); err != nil {
			diags.AddError("update instance environment variables", err.Error())

			return diags
		}
	}

	if toAdd := append(changes.upsert, changes.recreate...); len(toAdd) > 0 {
		if err := addEvars(ctx, client, id, toAdd); err != nil {
			diags.AddError("update instance environment variables", err.Error())

			return diags
		}
	}

	for _, name := range changes.remove {
		if err := deleteEvar(ctx, client, id, name); err != nil {
			diags.AddError("update instance environment variables", err.Error())

			return diags
		}
	}

	return diags
}

// addEvars issues one add-or-update call for the given variables.
func addEvars(ctx context.Context, client *sdk.APIClient, id int64, envs []sdk.AddEnvVariablesRequestEnvsInner) error {
	_, hresp, err := client.InstancesAPI.AddEnvVariables(ctx, id).
		AddEnvVariablesRequest(sdk.AddEnvVariablesRequest{Envs: envs}).
		Execute()
	if !envsWriteSucceeded(hresp) {
		return fmt.Errorf("adding or updating %d environment variable(s) on instance %d: %s",
			len(envs), id, errfmt.ErrMsg(err, hresp))
	}

	return nil
}

// deleteEvar removes one variable by name.
func deleteEvar(ctx context.Context, client *sdk.APIClient, id int64, name string) error {
	result, hresp, err := client.InstancesAPI.DeleteEnvVariable(ctx, id, name).Execute()
	if !envsWriteSucceeded(hresp) {
		return fmt.Errorf("removing environment variable %q from instance %d: %s",
			name, id, errfmt.ErrMsg(err, hresp))
	}

	// A 2xx is the signal (see envsWriteSucceeded), but when the echo did
	// decode, a body saying success:false is a refusal (MORPH-17987). A decode
	// error leaves result nil, which CheckDeleted treats as not refused.
	if err := errfmt.CheckDeleted(result, nil, hresp); err != nil {
		return fmt.Errorf("removing environment variable %q from instance %d: %w", name, id, err)
	}

	return nil
}

// envsWriteSucceeded reports whether an environment-variable write was accepted.
//
// The HTTP status is the signal, not the SDK's returned error. The envs
// endpoints echo the instance's full environment in their response, including
// read-only variables whose `value` is polymorphic (string or number), which
// the generated model cannot decode -- so on a perfectly successful write the
// SDK still returns a decode error. The read path (getInstanceEnvVars) makes
// the same judgement for the same reason. The write itself has been applied
// by the time a 2xx is returned; the response body is only an echo.
func envsWriteSucceeded(hresp *http.Response) bool {
	return hresp != nil && hresp.StatusCode >= 200 && hresp.StatusCode < 300
}

// applyEvarOptionsAfterCreate applies export/masked for any planned variable
// whose settings differ from the appliance defaults, immediately after create.
//
// The create request carries each variable's name and value only -- the API
// discards export and masked on that path -- so a variable created with, say,
// masked = true would otherwise land unmasked. And because the variable now
// already exists, a plain add-or-update would only touch its value: the API
// applies the options only when it creates a variable. So each affected
// variable is deleted and added again with its options. Variables at the
// defaults are left alone: create already produced exactly that. MORPH-17657.
func applyEvarOptionsAfterCreate(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	plan InstanceModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	planned, d := evarsByName(ctx, plan.Evars)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	var recreate []sdk.AddEnvVariablesRequestEnvsInner
	for name, want := range planned {
		if want.export && !want.masked {
			continue // appliance defaults; the create already produced this
		}
		recreate = append(recreate, evarRequest(name, want))
	}
	if len(recreate) == 0 {
		return diags
	}
	sort.Slice(recreate, func(i, j int) bool { return recreate[i].Name < recreate[j].Name })

	for _, entry := range recreate {
		if err := deleteEvar(ctx, client, id, entry.Name); err != nil {
			diags.AddError("apply instance environment variable options", err.Error())

			return diags
		}
	}
	if err := addEvars(ctx, client, id, recreate); err != nil {
		diags.AddError("apply instance environment variable options", err.Error())
	}

	return diags
}
