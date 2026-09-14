// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

const createOperation = "create tenant resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan TenantModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// Appliance version gate for parent_id (tenant hierarchy, 8.1.0+).
	// ModifyPlan applies the same gate, but it fails open when the version
	// cannot be read at plan time; repeating it here is what keeps an older
	// appliance from silently ignoring parentAccount and leaving state at odds
	// with the plan. Only configurations that set parent_id pay for the lookup.
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() {
		resp.Diagnostics.Append(versioncheck.RequireAttribute(
			ctx, client, path.Root("parent_id"), parentIdFeature, constants.TenantParentMinVersion,
		)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	account := sdk.AddTenantRequestAccount{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		account.Description.Set(plan.Description.ValueStringPointer())
	}
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
		account.Role = &sdk.AddTenantRequestAccountRole{
			Id: plan.BaseRoleId.ValueInt64Pointer(),
		}
	}
	// parentAccount is bound on create only and is only honoured for master
	// tenant callers. Sending it forces the new tenant under the nominated
	// parent instead of the caller's own tenant.
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() {
		account.ParentAccount = &sdk.AddTenantRequestAccountParentAccount{
			Id: plan.ParentId.ValueInt64Pointer(),
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

	createReq := sdk.AddTenantRequest{Account: account}

	result, hresp, err := client.TenantsAPI.AddTenant(ctx).
		AddTenantRequest(createReq).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "tenant", plan.Name.ValueString(), err, hresp)

		return
	}

	// Appliances before 8.1.0 report some validation failures as HTTP 200 with
	// success:false and no account; surface the envelope's message rather than
	// falling through to "Account ID is nil".
	if result != nil {
		if detail, failed := apiFailure(result.Success, result.AdditionalProperties); failed {
			resp.Diagnostics.AddError(
				createOperation,
				fmt.Sprintf("Morpheus rejected the request to create tenant %q: %s",
					plan.Name.ValueString(), detail),
			)

			return
		}
	}

	if result == nil || result.Account == nil || result.Account.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "Account ID is nil in the create response")

		return
	}

	id := *result.Account.Id

	state, notFound, diags := getTenantAsState(ctx, id, client, plan)
	if notFound {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("tenant %d was created but GET returned 404", id),
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "tenant",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "tenant",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// A nominated parent that the appliance did not bind means one of two
	// things: the appliance predates the tenant hierarchy (before 8.1.0) and
	// the version gate above could not read its version, or the caller is not
	// the master tenant, for which Morpheus ignores parentAccount. Either way
	// the tenant now exists under a different parent than planned. Terraform
	// would refuse the result as inconsistent anyway; say why, and taint the
	// tenant so the next apply replaces it once the configuration has been
	// corrected.
	if !plan.ParentId.IsNull() && !plan.ParentId.IsUnknown() &&
		(state.ParentId.IsNull() || state.ParentId.ValueInt64() != plan.ParentId.ValueInt64()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("parent_id"),
			createOperation,
			fmt.Sprintf(
				"tenant %d was created, but Morpheus did not place it under parent tenant %d "+
					"(the appliance reports parent %s). Either this appliance predates "+
					"Morpheus 8.1.0, which introduced the tenant hierarchy, or the provider "+
					"is not authenticated as the master tenant; in both cases parentAccount "+
					"is ignored. Remove parent_id, or upgrade the appliance / authenticate "+
					"as the master tenant, and apply again.",
				id, plan.ParentId.ValueInt64(), describeParent(state.ParentId),
			),
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "tenant",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// describeParent renders a read parent_id for a diagnostic: the id, or "none"
// when the tenant has no parent.
func describeParent(parentID types.Int64) string {
	if parentID.IsNull() || parentID.IsUnknown() {
		return "none"
	}

	return fmt.Sprintf("%d", parentID.ValueInt64())
}
