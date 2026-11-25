// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	var plan, state ApplianceSetupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateSetupChanges(plan, state)...)
}

func validateSetupChanges(plan, state ApplianceSetupModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	for _, field := range []struct {
		name     string
		planned  types.String
		previous types.String
	}{
		{"appliance_url", plan.ApplianceURL, state.ApplianceURL},
		{"appliance_name", plan.ApplianceName, state.ApplianceName},
		{"account_name", plan.AccountName, state.AccountName},
		{"username", plan.Username, state.Username},
		{"email", plan.Email, state.Email},
		{"first_name", plan.FirstName, state.FirstName},
		{"last_name", plan.LastName, state.LastName},
		{"hub_mode", plan.HubMode, state.HubMode},
	} {
		if !field.previous.IsNull() && !field.planned.IsUnknown() && !field.planned.Equal(field.previous) {
			diagnostics.AddAttributeError(path.Root(field.name), "Setup inputs cannot be changed",
				"Appliance setup runs only once. Manage the existing tenant, administrator, or appliance settings "+
					"through their dedicated resources or the appliance UI; this resource cannot apply this change.")
		}
	}

	return diagnostics
}
