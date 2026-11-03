// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// domainPasswordForUpdate decides what to send as networkDomain.domainPassword
// on an update. The field is ALWAYS sent, because omitting it clears the stored
// password (see passwordMask):
//
//   - re-sending (ModifyPlan marked write_only_revision unknown) and the value
//     is configured: send the configured value;
//   - the deprecated plaintext domain_password is configured: send it;
//   - otherwise: send the mask, which the API reads as "leave it unchanged".
func domainPasswordForUpdate(resend bool, passwordWo, password types.String) *string {
	if resend && !passwordWo.IsNull() && !passwordWo.IsUnknown() {
		return passwordWo.ValueStringPointer()
	}

	// The deprecated plaintext domain_password: sent on every update while
	// configured, exactly as the SDKv2 resource did (the API would otherwise
	// clear it). It is never configured together with domain_password_wo.
	if !password.IsNull() && !password.IsUnknown() {
		return password.ValueStringPointer()
	}

	mask := passwordMask

	return &mask
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState, config NetworkDomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	domain := &sdk.UpdateNetworkDomainRequestNetworkDomain{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		domain.Name = plan.Name.ValueStringPointer()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		domain.Description = plan.Description.ValueStringPointer()
	}

	if !plan.DomainUsername.IsNull() && !plan.DomainUsername.IsUnknown() {
		domain.DomainUsername = plan.DomainUsername.ValueStringPointer()
	}

	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		domain.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.PublicZone.IsNull() && !plan.PublicZone.IsUnknown() {
		domain.PublicZone = plan.PublicZone.ValueBoolPointer()
	}

	if !plan.DomainController.IsNull() && !plan.DomainController.IsUnknown() {
		domain.DomainController = plan.DomainController.ValueBoolPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		domain.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.TenantId.IsNull() && !plan.TenantId.IsUnknown() {
		domain.Account = &sdk.UpdateNetworkDomainRequestNetworkDomainAccount{
			Id: plan.TenantId.ValueInt64Pointer(),
		}
	}

	// domain_password is always sent: omitting it clears the stored password.
	resend := wodigest.Resend(ctx, req.Plan)
	domain.DomainPassword = domainPasswordForUpdate(resend, config.DomainPasswordWo, plan.DomainPassword)

	updateReq := sdk.UpdateNetworkDomainRequest{NetworkDomain: domain}

	id := currentState.Id.ValueInt64()

	_, httpResp, err := client.NetworksAPI.
		UpdateNetworkDomain(ctx, id).
		UpdateNetworkDomainRequest(updateReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating network domain",
			"network domain "+plan.Name.ValueString()+" PUT failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	// The PUT succeeded, so whatever was sent is what the appliance now holds:
	// record the digests before the read-back, so a failed read cannot leave a
	// stale digest that would re-send the password on the next plan.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	state, found, diags := getNetworkDomainAsState(ctx, id, client)
	if !found && !diags.HasError() {
		diags.AddError("read network domain", "domain not found immediately after update")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve write-only / deprecated marker from plan; bump the revision only
	// when the value was re-sent; always record the digests.
	state.DomainPasswordWoVersion = plan.DomainPasswordWoVersion
	state.DomainPassword = plan.DomainPassword
	state.AutoJoinDomain = plan.AutoJoinDomain
	state.WriteOnlyRevision = currentState.WriteOnlyRevision
	if resend {
		state.WriteOnlyRevision = wodigest.NewRevision()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
