// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data NetworkDomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, found, diags := getNetworkDomainAsState(ctx, data.Id.ValueInt64(), client)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// The domain was deleted out of band (Morpheus returned 404); remove it
	// from state so the next plan recreates it instead of erroring.
	if !found {
		tflog.Info(ctx, fmt.Sprintf("network domain %d not found, removing from state", data.Id.ValueInt64()))
		resp.State.RemoveResource(ctx)

		return
	}

	// Preserve write-only / deprecated values from prior state: the password
	// is never read from the API, and these are not returned by the GET.
	state.DomainPasswordWoVersion = data.DomainPasswordWoVersion
	state.DomainPassword = data.DomainPassword
	state.AutoJoinDomain = data.AutoJoinDomain
	state.WriteOnlyRevision = data.WriteOnlyRevision

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, NetworkDomainResourceSchema(ctx), &resp.State)...,
	)
}

func getNetworkDomainAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (NetworkDomainModel, bool, diag.Diagnostics) {
	var state NetworkDomainModel
	var diags diag.Diagnostics

	domainResp, httpResp, err := client.NetworksAPI.
		GetNetworkDomain(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return state, false, diags
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading network domain",
			fmt.Sprintf("network domain %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	d := domainResp.NetworkDomain
	if d == nil {
		diags.AddError(
			"error reading network domain",
			fmt.Sprintf("network domain %d response did not contain a domain", id),
		)

		return state, false, diags
	}

	state.Id = convert.Int64ToType(d.Id)
	state.Name = convert.StrToType(d.Name)
	state.Active = convert.BoolToType(d.Active)
	state.Description = convert.StrToType(d.Description.Get())
	state.Visibility = convert.StrToType(d.Visibility)
	state.DomainController = convert.BoolToType(d.DomainController)
	state.PublicZone = convert.BoolToType(d.PublicZone)
	state.DomainUsername = convert.StrToType(d.DomainUsername.Get())

	// tenant_id comes from account.id.
	if d.Account != nil && d.Account.Id != nil {
		state.TenantId = convert.Int64ToType(d.Account.Id)
	} else {
		state.TenantId = types.Int64Null()
	}

	// domain_password_wo is write-only; never read from the API.
	// domain_password_wo_version and write_only_revision are preserved by
	// callers from plan/state.

	return state, true, diags
}
