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
	upsert []sdk.AddEnvVariablesRequestEnvsInner
	remove []string
}

func (d evarsDiff) empty() bool {
	return len(d.upsert) == 0 && len(d.remove) == 0
}

// evarsByName flattens an evars set into name -> value. An unknown or null set
// contributes nothing; the caller decides what "nothing" means.
func evarsByName(ctx context.Context, set basetypes.SetValue) (map[string]string, diag.Diagnostics) {
	out := map[string]string{}
	if set.IsNull() || set.IsUnknown() {
		return out, nil
	}

	var elems []EvarsValue
	diags := set.ElementsAs(ctx, &elems, false)
	if diags.HasError() {
		return nil, diags
	}

	for _, e := range elems {
		out[e.Name.ValueString()] = e.Value.ValueString()
	}

	return out, nil
}

// diffEvars computes the changes needed to move from the current (state)
// variables to the planned ones.
//
// A planned variable that is absent from state, or present with a different
// value, is upserted. A state variable absent from the plan is removed --
// including when the plan omits `evars` entirely, which is how a practitioner
// removes every variable. Output is sorted by name so the calls are
// deterministic and the diff is easy to assert in tests.
func diffEvars(planned, current map[string]string) evarsDiff {
	var d evarsDiff

	for name, value := range planned {
		if cur, ok := current[name]; !ok || cur != value {
			d.upsert = append(d.upsert, sdk.AddEnvVariablesRequestEnvsInner{Name: name, Value: value})
		}
	}
	for name := range current {
		if _, ok := planned[name]; !ok {
			d.remove = append(d.remove, name)
		}
	}

	sort.Slice(d.upsert, func(i, j int) bool { return d.upsert[i].Name < d.upsert[j].Name })
	sort.Strings(d.remove)

	return d
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

	if len(changes.upsert) > 0 {
		_, hresp, err := client.InstancesAPI.AddEnvVariables(ctx, id).
			AddEnvVariablesRequest(sdk.AddEnvVariablesRequest{Envs: changes.upsert}).
			Execute()
		if !envsWriteSucceeded(hresp) {
			diags.AddError(
				"update instance environment variables",
				fmt.Sprintf("adding or updating %d environment variable(s) on instance %d: %s",
					len(changes.upsert), id, errfmt.ErrMsg(err, hresp)),
			)

			return diags
		}
	}

	for _, name := range changes.remove {
		_, hresp, err := client.InstancesAPI.DeleteEnvVariable(ctx, id, name).Execute()
		if !envsWriteSucceeded(hresp) {
			diags.AddError(
				"update instance environment variables",
				fmt.Sprintf("removing environment variable %q from instance %d: %s",
					name, id, errfmt.ErrMsg(err, hresp)),
			)

			return diags
		}
	}

	return diags
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
