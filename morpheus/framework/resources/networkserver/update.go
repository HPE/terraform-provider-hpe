// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

func (r *networkServerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState NetworkServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// service_password_wo is write-only: read it from req.Config.
	var servicePasswordWo, serviceTokenWo types.String
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("service_password_wo"), &servicePasswordWo)...,
	)
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("service_token_wo"), &serviceTokenWo)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	// Attributes removed from the configuration since the last apply, to be
	// cleared.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resend := wodigest.Resend(ctx, req.Plan)

	serverReq, diags := buildUpdateRequest(ctx, &plan, &priorState, servicePasswordWo, serviceTokenWo, pending, resend)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkServer(ctx, id).
		UpdateNetworkServerRequest(sdk.UpdateNetworkServerRequest{
			NetworkServer: serverReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.NetworksAPI.GetNetworkServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_server", "", err, httpResp)

		return
	}

	if readResult.NetworkServer == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkServer is nil in the response")

		return
	}

	state, err := networkServerAsState(ctx, &plan, readResult.NetworkServer)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map network server", err.Error())

		return
	}

	// write_only_revision is not in the API. Keep the prior value unless a
	// write-only value was re-sent (ModifyPlan marked it unknown).
	state.WriteOnlyRevision = priorState.WriteOnlyRevision
	if resend {
		state.WriteOnlyRevision = wodigest.NewRevision()
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildUpdateRequest builds the anyOf update body. The server binds a fixed set
// of fields on update: name, serviceUsername, servicePort, serviceHost,
// serviceUrl, serviceMode, networkFilter, tenantMatch, zoneId, visibility;
// servicePassword and serviceToken only when their key is present; plus config,
// credential and tenants. The SDN variant types only name, visibility, the two
// config booleans and tenants (controller_address is RequiresReplace); zoneId is
// sent alongside, since the server binds it for every type.
func buildUpdateRequest(
	ctx context.Context,
	plan, priorState *NetworkServerModel,
	servicePasswordWo, serviceTokenWo types.String,
	pending clearable.Pending,
	resend bool,
) (*sdk.UpdateNetworkServerRequestNetworkServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !plan.ConfigJuniperSdn.IsNull() && !plan.ConfigJuniperSdn.IsUnknown() {
		return &sdk.UpdateNetworkServerRequestNetworkServer{
			HPEMorpheusSDNNetworkServerUpdate: buildSdnUpdate(plan, pending),
		}, diags
	}

	nsx, ndiags := buildGenericUpdate(ctx, plan, priorState, servicePasswordWo, serviceTokenWo, pending, resend)
	diags.Append(ndiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &sdk.UpdateNetworkServerRequestNetworkServer{
		NSXNetworkServerUpdate: nsx,
	}, diags
}

func buildSdnUpdate(plan *NetworkServerModel, pending clearable.Pending) *sdk.HPEMorpheusSDNNetworkServerUpdate {
	cfg := plan.ConfigJuniperSdn

	upd := &sdk.HPEMorpheusSDNNetworkServerUpdate{
		Name:   plan.Name.ValueStringPointer(),
		Config: &sdk.HPEMorpheusSDNNetworkServerUpdateConfig{},
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		upd.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !cfg.ShowIntrospect.IsNull() && !cfg.ShowIntrospect.IsUnknown() {
		upd.Config.ShowIntrospect = cfg.ShowIntrospect.ValueBoolPointer()
	}
	if !cfg.EnableSubTenants.IsNull() && !cfg.EnableSubTenants.IsUnknown() {
		upd.Config.EnableSubTenants = cfg.EnableSubTenants.ValueBoolPointer()
	}

	upd.Tenants = sdnUpdateTenants(plan.TenantIds)

	// A removed tenant_ids is sent as an explicit empty list: the typed field
	// is omitempty, and an omitted key leaves the tenants unchanged.
	if pending.Clears("tenant_ids") {
		upd.AdditionalProperties = map[string]interface{}{"tenants": []any{}}
	}

	// The appliance binds zoneId for every type (NetworkService.updateNetworkServer),
	// so cloud_id is sent — and cleared with an explicit null — rather than
	// left to drift.
	clearable.NullableInt64(pending, "cloud_id", plan.CloudId, &upd.ZoneId)

	return upd
}

func buildGenericUpdate(
	ctx context.Context,
	plan, priorState *NetworkServerModel,
	servicePasswordWo, serviceTokenWo types.String,
	pending clearable.Pending,
	resend bool,
) (*sdk.NSXNetworkServerUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics

	nsx := &sdk.NSXNetworkServerUpdate{
		Name: plan.Name.ValueStringPointer(),
	}

	clearable.NullableString(pending, "service_url", plan.ServiceUrl, &nsx.ServiceUrl)
	clearable.NullableString(pending, "service_username", plan.ServiceUsername, &nsx.ServiceUsername)

	// servicePassword is sent only when the write-only digest says it changed.
	if resend && !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		nsx.ServicePassword = *sdk.NewNullableString(servicePasswordWo.ValueStringPointer())
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		nsx.Visibility = plan.Visibility.ValueStringPointer()
	}

	nsx.Tenants = nsxUpdateTenants(plan.TenantIds)
	if pending.Clears("tenant_ids") {
		nsx.AdditionalProperties = map[string]interface{}{"tenants": []any{}} // explicit empty: omitted means unchanged
	}

	clearable.NullableString(pending, "service_host", plan.ServiceHost, &nsx.ServiceHost)
	clearable.NullableString(pending, "service_mode", plan.ServiceMode, &nsx.ServiceMode)
	clearable.NullableString(pending, "network_filter", plan.NetworkFilter, &nsx.NetworkFilter)
	clearable.NullableString(pending, "tenant_match", plan.TenantMatch, &nsx.TenantMatch)
	clearable.NullableInt64(pending, "service_port", plan.ServicePort, &nsx.ServicePort)
	clearable.NullableInt64(pending, "cloud_id", plan.CloudId, &nsx.ZoneId)

	// serviceToken, like the password, is sent only when the write-only digest
	// says it changed.
	if resend {
		setNullableString(&nsx.ServiceToken, serviceTokenWo)
	}

	// Credential switching: send the id when set. When it is now null but was
	// set before, send {type: local} so the server drops the stored link (it
	// preserves the existing link when credential is omitted).
	switch {
	case !plan.CredentialId.IsNull() && !plan.CredentialId.IsUnknown():
		nsx.Credential = &sdk.NSXNetworkServerUpdateCredential{Id: plan.CredentialId.ValueInt64Pointer()}
	case !priorState.CredentialId.IsNull() && !priorState.CredentialId.IsUnknown():
		nsx.Credential = &sdk.NSXNetworkServerUpdateCredential{Type: sdk.PtrString("local")}
	}

	configMap, cdiags := configProperties(ctx, plan.Config)
	diags.Append(cdiags...)
	if diags.HasError() {
		return nil, diags
	}
	if configMap != nil {
		nsx.Config = &sdk.NSXNetworkServerUpdateConfig{AdditionalProperties: configMap}
	}

	return nsx, diags
}

func nsxUpdateTenants(ids types.Set) []sdk.NSXNetworkServerUpdateTenantsInner {
	values := tenantInts(ids)
	if values == nil {
		return nil
	}

	out := make([]sdk.NSXNetworkServerUpdateTenantsInner, 0, len(values))
	for _, id := range values {
		out = append(out, sdk.NSXNetworkServerUpdateTenantsInner{Id: idPtr(id)})
	}

	return out
}

func sdnUpdateTenants(ids types.Set) []sdk.HPEMorpheusSDNNetworkServerUpdateTenantsInner {
	values := tenantInts(ids)
	if values == nil {
		return nil
	}

	out := make([]sdk.HPEMorpheusSDNNetworkServerUpdateTenantsInner, 0, len(values))
	for _, id := range values {
		out = append(out, sdk.HPEMorpheusSDNNetworkServerUpdateTenantsInner{Id: idPtr(id)})
	}

	return out
}
