// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *networkServerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
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

	// service_password_wo is write-only: present only in the configuration,
	// never in the plan or state. Read it from req.Config.
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

	serverReq, diags := buildCreateRequest(ctx, &plan, servicePasswordWo, serviceTokenWo)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkServer(ctx).
		CreateNetworkServerRequest(sdk.CreateNetworkServerRequest{
			NetworkServer: serverReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.NetworkServer == nil || result.NetworkServer.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "NetworkServer ID is nil in the create response")

		return
	}

	id := *result.NetworkServer.Id

	readResult, httpResp, err := client.NetworksAPI.GetNetworkServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_server", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.NetworkServer == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkServer is nil in the response")
		// The create succeeded; keep the id so the next apply can finish or
		// destroy it instead of orphaning it on the appliance.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	state, err := networkServerAsState(ctx, &plan, readResult.NetworkServer)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map network server", err.Error())
		// The server exists; keep its id rather than orphan it.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	state.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildCreateRequest builds the anyOf create body. The SDN type
// (morpheus-junipersdn-network) uses the HPEMorpheusSDNNetworkServer variant;
// every other type uses NSXNetworkServer as the generic variant, which carries
// every field the server binds in
// the typed fields. The Morpheus server binds a fixed set of fields on
// create: name, serviceUsername, servicePassword, servicePort, serviceHost,
// serviceUrl, serviceMode, networkFilter, tenantMatch, zoneId, serviceToken,
// visibility, type, config, credential and tenants. enabled and description
// are read-only and never sent.
func buildCreateRequest(
	ctx context.Context,
	plan *NetworkServerModel,
	servicePasswordWo, serviceTokenWo types.String,
) (*sdk.CreateNetworkServerRequestNetworkServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	typeCode := plan.TypeCode.ValueString()
	planHasSdn := !plan.ConfigJuniperSdn.IsNull() && !plan.ConfigJuniperSdn.IsUnknown()

	if planHasSdn {
		if typeCode != "" && typeCode != typeCodeJuniperSdn {
			diags.AddError(
				"Conflicting network server type",
				fmt.Sprintf(
					"config_juniper_sdn is only valid for the %q type, but type_code is %q.",
					typeCodeJuniperSdn, typeCode,
				),
			)

			return nil, diags
		}

		sdn := buildSdnCreate(plan, servicePasswordWo)

		return &sdk.CreateNetworkServerRequestNetworkServer{
			HPEMorpheusSDNNetworkServer: sdn,
		}, diags
	}

	if typeCode == typeCodeJuniperSdn {
		diags.AddError(
			"Missing SDN configuration",
			fmt.Sprintf("type_code %q requires a config_juniper_sdn block.", typeCodeJuniperSdn),
		)

		return nil, diags
	}

	nsx, ndiags := buildGenericCreate(ctx, plan, servicePasswordWo, serviceTokenWo)
	diags.Append(ndiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &sdk.CreateNetworkServerRequestNetworkServer{
		NSXNetworkServer: nsx,
	}, diags
}

func buildSdnCreate(plan *NetworkServerModel, _ types.String) *sdk.HPEMorpheusSDNNetworkServer {
	cfg := plan.ConfigJuniperSdn

	sdn := &sdk.HPEMorpheusSDNNetworkServer{
		Type: typeCodeJuniperSdn,
		Name: plan.Name.ValueString(),
		Config: sdk.HPEMorpheusSDNNetworkServerConfig{
			ControllerAddress: cfg.ControllerAddress.ValueString(),
		},
	}

	if !cfg.ShowIntrospect.IsNull() && !cfg.ShowIntrospect.IsUnknown() {
		sdn.Config.ShowIntrospect = cfg.ShowIntrospect.ValueBoolPointer()
	}
	if !cfg.EnableSubTenants.IsNull() && !cfg.EnableSubTenants.IsUnknown() {
		sdn.Config.EnableSubTenants = cfg.EnableSubTenants.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		sdn.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		sdn.ZoneId = plan.CloudId.ValueInt64Pointer()
	}

	sdn.Tenants = sdnTenants(plan.TenantIds)

	return sdn
}

// buildGenericCreate builds the NSXNetworkServer variant used for every
// non-SDN type. The struct carries every field the server binds; the
// free-form config map goes into the typed Config's additional properties,
// keyed as the type's option types name them.
func buildGenericCreate(
	ctx context.Context,
	plan *NetworkServerModel,
	servicePasswordWo, serviceTokenWo types.String,
) (*sdk.NSXNetworkServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	nsx := &sdk.NSXNetworkServer{
		Type: plan.TypeCode.ValueString(),
		Name: plan.Name.ValueString(),
	}

	// serviceUrl is a required (non-omitempty) NullableString on the struct;
	// send whatever the user set, or an explicit null when unset.
	if !plan.ServiceUrl.IsNull() && !plan.ServiceUrl.IsUnknown() {
		nsx.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	} else {
		nsx.ServiceUrl = *sdk.NewNullableString(nil)
	}
	if !plan.ServiceUsername.IsNull() && !plan.ServiceUsername.IsUnknown() {
		nsx.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		nsx.ServicePassword = *sdk.NewNullableString(servicePasswordWo.ValueStringPointer())
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		nsx.Visibility = plan.Visibility.ValueStringPointer()
	}
	// zoneId is required by the API for every non-SDN server type (the model
	// serialises it unconditionally, so an unset cloud would be sent as 0).
	if plan.CloudId.IsNull() || plan.CloudId.IsUnknown() {
		diags.AddError(
			"cloud_id is required",
			fmt.Sprintf("network server type %q must be attached to a cloud; set cloud_id.", plan.TypeCode.ValueString()),
		)

		return nil, diags
	}
	nsx.ZoneId = plan.CloudId.ValueInt64()

	// credential_id selects a stored credential. The server translates the id
	// carried in {id} to {type: id}. When it is unset the server defaults to
	// local credentials (serviceUsername/servicePassword).
	if !plan.CredentialId.IsNull() && !plan.CredentialId.IsUnknown() {
		nsx.Credential = &sdk.NSXNetworkServerCredential{
			Id: plan.CredentialId.ValueInt64Pointer(),
		}
	}

	nsx.Tenants = nsxTenants(plan.TenantIds)

	setNullableString(&nsx.ServiceHost, plan.ServiceHost)
	setNullableString(&nsx.ServiceMode, plan.ServiceMode)
	setNullableString(&nsx.NetworkFilter, plan.NetworkFilter)
	setNullableString(&nsx.TenantMatch, plan.TenantMatch)
	setNullableString(&nsx.ServiceToken, serviceTokenWo)
	if !plan.ServicePort.IsNull() && !plan.ServicePort.IsUnknown() {
		nsx.ServicePort.Set(plan.ServicePort.ValueInt64Pointer())
	}

	// The free-form config map holds type-specific option values keyed as the
	// type's option types name them.
	configMap, cdiags := configProperties(ctx, plan.Config)
	diags.Append(cdiags...)
	if diags.HasError() {
		return nil, diags
	}
	if configMap != nil {
		nsx.Config = &sdk.NSXNetworkServerConfig{AdditionalProperties: configMap}
	}

	return nsx, diags
}

// setNullableString sets a nullable request field from a known, non-null
// value and leaves it unsent otherwise.
func setNullableString(target *sdk.NullableString, v types.String) {
	if !v.IsNull() && !v.IsUnknown() {
		target.Set(v.ValueStringPointer())
	}
}

// configProperties converts the free-form config attribute to the map the
// request's config object carries; nil when the attribute is unset.
func configProperties(ctx context.Context, config types.Dynamic) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if config.IsNull() || config.IsUnknown() || config.IsUnderlyingValueNull() || config.IsUnderlyingValueUnknown() {
		return nil, diags
	}

	raw, err := convert.ValueToAny(ctx, config.UnderlyingValue())
	if err != nil {
		diags.AddError("Invalid config", fmt.Sprintf("failed to convert config: %s", err))

		return nil, diags
	}

	m, ok := raw.(map[string]interface{})
	if !ok {
		diags.AddError("Invalid config", fmt.Sprintf("config must be an object of option values, got %T", raw))

		return nil, diags
	}

	return m, diags
}

func nsxTenants(ids types.Set) []sdk.NSXNetworkServerTenantsInner {
	values := tenantInts(ids)
	if values == nil {
		return nil
	}

	out := make([]sdk.NSXNetworkServerTenantsInner, 0, len(values))
	for _, id := range values {
		out = append(out, sdk.NSXNetworkServerTenantsInner{Id: idPtr(id)})
	}

	return out
}

func sdnTenants(ids types.Set) []sdk.HPEMorpheusSDNNetworkServerTenantsInner {
	values := tenantInts(ids)
	if values == nil {
		return nil
	}

	out := make([]sdk.HPEMorpheusSDNNetworkServerTenantsInner, 0, len(values))
	for _, id := range values {
		out = append(out, sdk.HPEMorpheusSDNNetworkServerTenantsInner{Id: idPtr(id)})
	}

	return out
}

func tenantInts(ids types.Set) []int64 {
	if ids.IsNull() || ids.IsUnknown() {
		return nil
	}

	var out []int64
	for _, e := range ids.Elements() {
		if v, ok := e.(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
			out = append(out, v.ValueInt64())
		}
	}

	return out
}

func idPtr(v int64) *int64 { return &v }
