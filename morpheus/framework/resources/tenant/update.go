// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const updateOperation = "update tenant resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan TenantModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var prior TenantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The master tenant may be updated (name, description, currency, subdomain
	// and the billing fields), but Morpheus rejects disabling it with an opaque
	// 400 and silently ignores a role assigned to it. Refuse both here with a
	// clear, named diagnostic; the ignored role would otherwise be reported
	// back as if it had been applied, because the read preserves the planned
	// base_role_id when the API returns no role.
	if prior.Master.ValueBool() {
		if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() && !plan.Enabled.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				path.Root("enabled"), updateOperation, masterTenantDisableMsg)

			return
		}
		if !plan.BaseRoleId.IsNull() && !plan.BaseRoleId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("base_role_id"), updateOperation, masterTenantRoleMsg)

			return
		}
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := plan.Id.ValueInt64()

	account := sdk.UpdateTenantRequestAccount{
		Name: plan.Name.ValueStringPointer(),
	}

	// description is the one attribute a practitioner can clear by removing it
	// from the configuration (it is Optional, not Computed, precisely so the
	// old value is not carried forward), so it is always sent: a null plan
	// value becomes an explicit JSON null, which Morpheus binds as an empty
	// description. Omitting the key instead would leave the previous value in
	// place and the read-back would disagree with the plan.
	account.Description.Set(plan.Description.ValueStringPointer())
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		account.Active = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Subdomain.IsNull() && !plan.Subdomain.IsUnknown() {
		account.Subdomain.Set(plan.Subdomain.ValueStringPointer())
	}
	if !plan.Currency.IsNull() && !plan.Currency.IsUnknown() {
		account.Currency = plan.Currency.ValueStringPointer()
	}
	if !plan.BaseRoleId.IsNull() && !plan.BaseRoleId.IsUnknown() {
		account.Role = &sdk.UpdateTenantRequestAccountRole{
			Id: plan.BaseRoleId.ValueInt64Pointer(),
		}
	}
	if !plan.AccountNumber.IsNull() && !plan.AccountNumber.IsUnknown() {
		account.AccountNumber.Set(plan.AccountNumber.ValueStringPointer())
	}
	if !plan.AccountName.IsNull() && !plan.AccountName.IsUnknown() {
		account.AccountName.Set(plan.AccountName.ValueStringPointer())
	}
	if !plan.CustomerNumber.IsNull() && !plan.CustomerNumber.IsUnknown() {
		account.CustomerNumber.Set(plan.CustomerNumber.ValueStringPointer())
	}

	updateReq := sdk.UpdateTenantRequest{Account: account}

	result, hresp, err := client.TenantsAPI.UpdateTenant(ctx, id).
		UpdateTenantRequest(updateReq).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "tenant", plan.Name.ValueString(), err, hresp)

		return
	}

	// Appliances before 8.1.0 report some validation failures as HTTP 200 with
	// success:false; surface the envelope's message. The update model has no
	// typed success field, so it is read from the additional properties.
	if result != nil {
		if detail, failed := apiFailure(nil, result.AdditionalProperties); failed {
			resp.Diagnostics.AddError(
				updateOperation,
				fmt.Sprintf("Morpheus rejected the request to update tenant %d: %s", id, detail),
			)

			return
		}
	}

	state, notFound, diags := getTenantAsState(ctx, id, client, plan)
	if notFound {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("tenant %d was updated but GET returned 404", id),
		)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
