// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, config NetworkDomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	domain := &sdk.CreateNetworkDomainRequestNetworkDomain{
		Name: plan.Name.ValueStringPointer(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		domain.Description = plan.Description.ValueStringPointer()
	}

	if !plan.DomainUsername.IsNull() && !plan.DomainUsername.IsUnknown() {
		domain.DomainUsername = plan.DomainUsername.ValueStringPointer()
	}

	// Optional+Computed bools carry a schema default, so plan always has a
	// value; send it.
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

	// tenant_id maps to account.id, and is only honoured for the master
	// tenant; a sub-tenant owner is ignored by the API and forced to private.
	if !plan.TenantId.IsNull() && !plan.TenantId.IsUnknown() {
		domain.Account = &sdk.CreateNetworkDomainRequestNetworkDomainAccount{
			Id: plan.TenantId.ValueInt64Pointer(),
		}
	}

	// On create, send the password only when it is configured; omit otherwise.
	// The deprecated plaintext domain_password is honoured until its removal
	// (ConflictsWith keeps the two exclusive).
	switch {
	case !config.DomainPasswordWo.IsNull() && !config.DomainPasswordWo.IsUnknown():
		domain.DomainPassword = config.DomainPasswordWo.ValueStringPointer()
	case !plan.DomainPassword.IsNull() && !plan.DomainPassword.IsUnknown():
		domain.DomainPassword = plan.DomainPassword.ValueStringPointer()
	}

	createReq := sdk.CreateNetworkDomainRequest{NetworkDomain: domain}

	createResp, httpResp, err := client.NetworksAPI.
		CreateNetworkDomain(ctx).
		CreateNetworkDomainRequest(createReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating network domain",
			"network domain "+plan.Name.ValueString()+" POST failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	created := createResp.NetworkDomain
	if created == nil || created.Id == nil {
		resp.Diagnostics.AddError(
			"error creating network domain",
			"network domain "+plan.Name.ValueString()+": id is nil",
		)

		return
	}

	id := *created.Id

	state, found, diags := getNetworkDomainAsState(ctx, id, client)
	if !found && !diags.HasError() {
		diags.AddError("read network domain", "domain not found immediately after create")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_domain",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// Preserve write-only / deprecated marker from plan; record the digests.
	state.DomainPasswordWoVersion = plan.DomainPasswordWoVersion
	state.DomainPassword = plan.DomainPassword
	state.AutoJoinDomain = plan.AutoJoinDomain
	state.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
