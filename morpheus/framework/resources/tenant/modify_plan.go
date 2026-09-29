// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

// ModifyPlan performs plan-time validation the generated schema cannot express:
//   - currency membership against the appliance's live currency list
//     (best-effort; requires an API client)
//   - the appliance version gate on parent_id (fail-open; requires an API
//     client)
//   - the parent_id / base_role_id relationship (best-effort; requires an API
//     client)
//
// The subdomain format and not-entirely-numeric rules are both enforced by
// schema-level regex_matches validators. Nothing to validate on destroy
// (Plan.Raw is null).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan TenantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Currency.IsNull() && !plan.Currency.IsUnknown() {
		r.validateCurrency(ctx, plan.Currency.ValueString(), &resp.Diagnostics)
	}

	// When a parent tenant is explicitly nominated, validate it.
	//
	// First the appliance version: parent_id is only honoured from 8.1.0,
	// which introduced the tenant hierarchy; older appliances ignore it, which
	// would leave state disagreeing with the plan and force a replacement on
	// every run. Refuse the plan outright on a known-too-old appliance, and
	// skip the remaining parent_id checks, which presuppose the feature exists.
	//
	// Then the caller: parent_id is only honoured for the master tenant, and
	// the assignable-roles endpoint differs by caller, so resolve the caller's
	// master status once and drive both parent_id validators from it. A
	// failure to determine it is surfaced as a warning rather than blocking
	// the plan.
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() {
		if r.parentIdUnsupported(ctx, &resp.Diagnostics) {
			return
		}

		client, err := r.NewClient(ctx)

		var isMaster bool
		if err == nil {
			isMaster, err = tenancy.CallerIsMaster(ctx, client)
		}

		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("parent_id"),
				"Could not validate parent_id at plan time",
				"Unable to determine whether the authenticated tenant is the master "+
					"tenant, so parent_id and base_role_id were not validated; "+
					"Morpheus will enforce them on apply. Error: "+err.Error(),
			)
		} else {
			parentID := plan.ParentId.ValueInt64()
			validateParentIdRequiresMaster(parentID, isMaster, &resp.Diagnostics)
			r.validateAvailableParentTenantRoles(
				ctx, parentID, plan.BaseRoleId, isMaster, &resp.Diagnostics)
		}
	}
}

// parentIdUnsupported applies the appliance version gate for parent_id and
// reports whether it refused. It appends the attribute error to diags when the
// appliance is known to be older than constants.TenantParentMinVersion, and
// fails open — returning false with no diagnostic — when the version cannot be
// determined or no API client is available (offline `terraform validate`).
// Create applies the same gate, so a plan that passed only because the version
// was unreadable is still caught before the API is called.
func (r *Resource) parentIdUnsupported(ctx context.Context, diags *diag.Diagnostics) bool {
	client, err := r.NewClient(ctx)
	if err != nil {
		return false
	}

	gate := versioncheck.RequireAttribute(
		ctx, client, path.Root("parent_id"), parentIdFeature, constants.TenantParentMinVersion,
	)
	diags.Append(gate...)

	return gate.HasError()
}

// validateCurrency checks the planned currency against the appliance's live
// "currencies" option source. Best-effort: if the client is unavailable (e.g.
// offline `terraform validate`) or the lookup fails, the plan is not blocked --
// Morpheus enforces the currency on apply.
func (r *Resource) validateCurrency(
	ctx context.Context,
	currency string,
	diags *diag.Diagnostics,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		return
	}

	resp, hresp, err := client.OptionsAPI.GetOptionSourceData(ctx, "currencies").Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		diags.AddAttributeWarning(
			path.Root("currency"),
			"Could not validate currency at plan time",
			"Unable to fetch the appliance currency list to validate currency; "+
				"Morpheus will enforce it on apply. Error: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if resp == nil {
		return
	}

	valid := make(map[string]struct{}, len(resp.Data))
	for _, item := range resp.Data {
		if v, ok := item["value"].(string); ok && v != "" {
			valid[v] = struct{}{}
		}
	}

	if len(valid) == 0 {
		return
	}

	if _, ok := valid[currency]; !ok {
		codes := make([]string, 0, len(valid))
		for c := range valid {
			codes = append(codes, c)
		}
		sort.Strings(codes)

		diags.AddAttributeError(
			path.Root("currency"),
			"Invalid currency",
			fmt.Sprintf(
				"currency %q is not enabled on the target appliance. Enabled currencies: %s.",
				currency, strings.Join(codes, ", ")),
		)
	}
}
