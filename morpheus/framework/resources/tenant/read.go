// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

const readOperation = "read tenant resource"

// getTenantAsState fetches the tenant by ID and maps it into a TenantModel.
//
// The prior plan/state model is required so that remove_resources -- a synthetic
// delete-time-only flag that is never sent on create/update and never returned
// by the API -- can be preserved from prior/plan.
//
// parent_id is computed_optional: when omitted the API assigns the
// authenticated user's tenant as the parent, so it is populated from the read
// response's parent.id (with UseStateForUnknown keeping it stable across plans).
//
// Returns (state, notFound, diags). notFound is true when the API returned 404;
// the caller should remove the resource from state (Read) or error (Create).
func getTenantAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	prior TenantModel,
) (TenantModel, bool, diag.Diagnostics) {
	var state TenantModel
	var diags diag.Diagnostics

	resp, hresp, err := client.TenantsAPI.GetTenant(ctx, id).Execute()
	if errfmt.IsNotFound(hresp) {
		return state, true, diags
	}
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		diags.AddError(
			readOperation,
			fmt.Sprintf("tenant %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, false, diags
	}

	if resp == nil || resp.Account == nil {
		diags.AddError("API returned nil", "Account is nil in the tenant response")

		return state, false, diags
	}

	a := resp.Account

	state.Id = convert.Int64ToType(a.Id)
	state.Name = convert.StrToType(a.Name)
	state.Description = convert.StrToType(a.Description.Get())
	state.Enabled = convert.BoolToType(a.Active)
	state.Subdomain = convert.StrToType(a.Subdomain)
	state.Currency = convert.StrToType(a.Currency)
	state.AccountNumber = convert.StrToType(a.AccountNumber.Get())
	state.AccountName = convert.StrToType(a.AccountName)
	state.CustomerNumber = convert.StrToType(a.CustomerNumber.Get())
	state.ExternalId = convert.StrToType(a.ExternalId.Get())
	state.Master = convert.BoolToType(a.Master)
	state.DateCreated = convert.TimeToType(a.DateCreated)
	state.LastUpdated = convert.TimeToType(a.LastUpdated)

	if a.Role != nil {
		state.BaseRoleId = convert.Int64ToType(a.Role.Id)
		state.BaseRoleName = convert.StrToType(a.Role.Authority)
	} else {
		// No role is returned for the master tenant, which has no base role.
		// Preserve the prior value where it is known, but never carry an
		// unknown into state: on an update plan a computed_optional attribute
		// with no state value is unknown, and unknowns are invalid after apply.
		state.BaseRoleId = prior.BaseRoleId
		if state.BaseRoleId.IsUnknown() {
			state.BaseRoleId = convert.Int64ToType(nil)
		}
		state.BaseRoleName = convert.StrToType(nil)
	}

	if a.Parent != nil {
		state.ParentId = convert.Int64ToType(a.Parent.Id)
		state.ParentName = convert.StrToType(a.Parent.Name)
		state.ParentSubdomain = convert.StrToType(a.Parent.Subdomain.Get())
	} else {
		// The master tenant has no parent.
		state.ParentId = convert.Int64ToType(nil)
		state.ParentName = convert.StrToType(nil)
		state.ParentSubdomain = convert.StrToType(nil)
	}

	if a.Stats != nil {
		state.InstanceCount = convert.Int64ToType(a.Stats.InstanceCount)
		state.UserCount = convert.Int64ToType(a.Stats.UserCount)
	} else {
		state.InstanceCount = convert.Int64ToType(nil)
		state.UserCount = convert.Int64ToType(nil)
	}

	// Not returned by the GET endpoint - preserve from prior/plan.
	state.RemoveResources = prior.RemoveResources

	return state, false, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var prior TenantModel

	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	state, notFound, diags := getTenantAsState(ctx, prior.Id.ValueInt64(), client, prior)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, TenantResourceSchema(ctx), &resp.State)...,
	)
}
