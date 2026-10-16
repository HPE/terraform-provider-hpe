package networkpoolserver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

var (
	_ resource.Resource                   = &networkPoolServerResource{}
	_ resource.ResourceWithConfigure      = &networkPoolServerResource{}
	_ resource.ResourceWithImportState    = &networkPoolServerResource{}
	_ resource.ResourceWithModifyPlan     = &networkPoolServerResource{}
	_ resource.ResourceWithValidateConfig = &networkPoolServerResource{}
)

type networkPoolServerResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkPoolServerResource{}
}

func (r *networkPoolServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_pool_server"
}

func (r *networkPoolServerResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkPoolServerResourceSchema(ctx)
}

func (r *networkPoolServerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkPoolServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// service_password_wo is a write-only attribute: its value is only present in
	// the configuration, never in the plan or state. Read it from req.Config.
	var servicePasswordWo types.String
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("service_password_wo"), &servicePasswordWo)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	// The SDK uses a oneOf union for this request. We use InfobloxNetworkPoolServer as the
	// concrete type because it is the superset of all pool server types (Infoblox, Bluecat,
	// phpIPAM, SolarWinds, EfficientIP/SOLIDserver, ...). All common fields are accepted
	// regardless of type. The API selects the actual type from the "type" code in the request
	// body (NetworkPoolServerType.findByCode), which is driven by either type_code (sent
	// directly) or type_id (resolved to its code below).
	typeCode, typeDiags := resolveNetworkPoolServerTypeCode(ctx, client, plan.TypeCode, plan.TypeId)
	resp.Diagnostics.Append(typeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	infoblox := &sdk.InfobloxNetworkPoolServer{}
	infoblox.Type = typeCode
	infoblox.Name = plan.Name.ValueString()
	if !plan.ServiceUrl.IsNull() {
		infoblox.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	}
	if !plan.ServiceUsername.IsNull() {
		infoblox.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		infoblox.ServicePassword = *sdk.NewNullableString(servicePasswordWo.ValueStringPointer())
	}
	if !plan.IgnoreSsl.IsNull() {
		infoblox.IgnoreSsl = plan.IgnoreSsl.ValueBoolPointer()
	}
	if !plan.Enabled.IsNull() {
		infoblox.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.NetworkFilter.IsNull() {
		infoblox.NetworkFilter = *sdk.NewNullableString(plan.NetworkFilter.ValueStringPointer())
	}
	if !plan.ZoneFilter.IsNull() {
		infoblox.ZoneFilter = *sdk.NewNullableString(plan.ZoneFilter.ValueStringPointer())
	}
	if !plan.TenantMatch.IsNull() {
		infoblox.TenantMatch = *sdk.NewNullableString(plan.TenantMatch.ValueStringPointer())
	}
	if !plan.ServiceMode.IsNull() {
		infoblox.ServiceMode = plan.ServiceMode.ValueStringPointer()
	}
	if !plan.ServiceThrottleRate.IsNull() {
		rate := plan.ServiceThrottleRate.ValueInt64()
		infoblox.ServiceThrottleRate = *sdk.NewNullableInt64(&rate)
	}

	// inventory_existing is a fieldContext=config option (config.inventoryExisting) shared by
	// all pool server types. It is stored in the integration's generic config map.
	if !plan.InventoryExisting.IsNull() && !plan.InventoryExisting.IsUnknown() {
		infoblox.Config = &sdk.InfobloxNetworkPoolServerConfig{
			InventoryExisting: convert.BoolTypeToStringPointerOnOff(plan.InventoryExisting),
		}
	}

	// Credential: use credential_id for stored credentials
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		cred := &sdk.InfobloxNetworkPoolServerCredential{}
		cred.Type = &idStr
		infoblox.Credential = cred
	}

	// Convert the superset to the typed variant for the type (see request.go).
	serverReq := createRequest(typeCode, infoblox, appIDValue(plan.AppId))

	result, httpResp, err := client.NetworksAPI.CreateNetworkPoolServer(ctx).
		CreateNetworkPoolServerRequest(sdk.CreateNetworkPoolServerRequest{
			NetworkPoolServer: &serverReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.NetworkPoolServer == nil || result.NetworkPoolServer.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "NetworkPoolServer ID is nil in the create response")

		return
	}

	id := *result.NetworkPoolServer.Id

	readResult, httpResp, err := client.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool_server", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_pool_server",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.NetworkPoolServer == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPoolServer is nil in the response")

		return
	}
	mapReadResponseToModel(ctx, &plan, readResult.NetworkPoolServer)

	plan.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkPoolServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool_server", "", err, httpResp)

		return
	}

	server := result.NetworkPoolServer
	if server == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPoolServer is nil in the response")

		return
	}
	mapReadResponseToModel(ctx, &state, server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, NetworkPoolServerResourceSchema(ctx), &resp.State)...,
	)
}

func (r *networkPoolServerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkPoolServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// service_password_wo is a write-only attribute: read it from req.Config.
	var servicePasswordWo types.String
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, path.Root("service_password_wo"), &servicePasswordWo)...,
	)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	infobloxUpdate := &sdk.InfobloxNetworkPoolServerUpdate{}
	infobloxUpdate.Name = plan.Name.ValueStringPointer()
	clearable.NullableString(pending, "service_url", plan.ServiceUrl, &infobloxUpdate.ServiceUrl)
	clearable.NullableString(pending, "service_username", plan.ServiceUsername, &infobloxUpdate.ServiceUsername)
	// Re-sent when ModifyPlan marked write_only_revision unknown (its digest
	// changed, or the deprecated marker did). See write_only.go.
	if wodigest.Resend(ctx, req.Plan) && !servicePasswordWo.IsNull() && !servicePasswordWo.IsUnknown() {
		infobloxUpdate.ServicePassword = *sdk.NewNullableString(servicePasswordWo.ValueStringPointer())
	}
	if !plan.IgnoreSsl.IsNull() {
		infobloxUpdate.IgnoreSsl = plan.IgnoreSsl.ValueBoolPointer()
	}
	if !plan.Enabled.IsNull() {
		infobloxUpdate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	clearable.NullableString(pending, "network_filter", plan.NetworkFilter, &infobloxUpdate.NetworkFilter)
	clearable.NullableString(pending, "zone_filter", plan.ZoneFilter, &infobloxUpdate.ZoneFilter)
	clearable.NullableString(pending, "tenant_match", plan.TenantMatch, &infobloxUpdate.TenantMatch)
	if !plan.ServiceMode.IsNull() {
		infobloxUpdate.ServiceMode = plan.ServiceMode.ValueStringPointer()
	}
	if !plan.ServiceThrottleRate.IsNull() {
		rate := plan.ServiceThrottleRate.ValueInt64()
		infobloxUpdate.ServiceThrottleRate = *sdk.NewNullableInt64(&rate)
	}

	// inventory_existing (config.inventoryExisting) is shared by all pool server types.
	if !plan.InventoryExisting.IsNull() && !plan.InventoryExisting.IsUnknown() {
		infobloxUpdate.Config = &sdk.InfobloxNetworkPoolServerUpdateConfig{
			InventoryExisting: convert.BoolTypeToStringPointerOnOff(plan.InventoryExisting),
		}
	}

	// Credential: use credential_id for stored credentials
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		cred := &sdk.InfobloxNetworkPoolServerUpdateCredential{}
		cred.Type = &idStr
		infobloxUpdate.Credential = cred
	}

	// type_code replaces on change, so the state's code is the server's type.
	serverReq := updateRequest(plan.TypeCode.ValueString(), infobloxUpdate, appIDValue(plan.AppId))

	_, httpResp, err := client.NetworksAPI.UpdateNetworkPoolServer(ctx, id).
		UpdateNetworkPoolServerRequest(sdk.UpdateNetworkPoolServerRequest{
			NetworkPoolServer: &serverReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_pool_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Re-read to get current state
	readResult, httpResp, err := client.NetworksAPI.GetNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_pool_server", "", err, httpResp)

		return
	}

	server := readResult.NetworkPoolServer
	if server == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkPoolServer is nil in the response")

		return
	}
	mapReadResponseToModel(ctx, &plan, server)

	// write_only_revision is not in the API. Keep the prior value unless a
	// write-only value was re-sent (ModifyPlan marked it unknown); and always
	// record the digests of the configured values. See write_only.go.
	var priorState NetworkPoolServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	plan.WriteOnlyRevision = priorState.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkPoolServerResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkPoolServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.DeleteNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_pool_server", "", err, httpResp)

		return
	}
}

func (r *networkPoolServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// "<id>" or "<id>#service_password_wo=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", rawID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapReadResponseToModel(
	ctx context.Context,
	model *NetworkPoolServerModel,
	server *sdk.GetNetworkPoolServer200ResponseNetworkPoolServer,
) {
	if server.Id != nil {
		model.Id = types.Int64Value(*server.Id)
	}
	if server.Name != nil {
		model.Name = types.StringValue(*server.Name)
	}
	// type_id and type_code are mutually exclusive on input, but the API returns both
	// the numeric id and the stable code for the resolved type, so reflect both actual
	// values in state. ConflictsWith is config-only (it never inspects state), so having
	// both populated in state does not trigger a validation error.
	if t := server.Type; t != nil {
		model.TypeId = convert.Int64ToType(t.Id)
		model.TypeCode = convert.StrToType(t.Code)
	}
	if server.ServiceUrl.IsSet() && server.ServiceUrl.Get() != nil {
		model.ServiceUrl = types.StringValue(*server.ServiceUrl.Get())
	} else {
		model.ServiceUrl = types.StringNull()
	}
	if server.IgnoreSsl.IsSet() && server.IgnoreSsl.Get() != nil {
		model.IgnoreSsl = types.BoolValue(*server.IgnoreSsl.Get())
	}
	if server.Enabled != nil {
		model.Enabled = types.BoolValue(*server.Enabled)
	}
	if server.Status != nil {
		model.Status = types.StringValue(*server.Status)
	} else {
		model.Status = types.StringNull()
	}
	if server.NetworkFilter.IsSet() && server.NetworkFilter.Get() != nil {
		model.NetworkFilter = types.StringValue(*server.NetworkFilter.Get())
	} else {
		model.NetworkFilter = types.StringNull()
	}
	if server.ZoneFilter.IsSet() && server.ZoneFilter.Get() != nil {
		model.ZoneFilter = types.StringValue(*server.ZoneFilter.Get())
	} else {
		model.ZoneFilter = types.StringNull()
	}
	if server.TenantMatch.IsSet() && server.TenantMatch.Get() != nil {
		model.TenantMatch = types.StringValue(*server.TenantMatch.Get())
	} else {
		model.TenantMatch = types.StringNull()
	}
	if server.ServiceMode.IsSet() && server.ServiceMode.Get() != nil {
		model.ServiceMode = types.StringValue(*server.ServiceMode.Get())
	} else {
		model.ServiceMode = types.StringNull()
	}
	if server.ServiceThrottleRate.IsSet() && server.ServiceThrottleRate.Get() != nil {
		model.ServiceThrottleRate = types.Int64Value(*server.ServiceThrottleRate.Get())
	} else {
		model.ServiceThrottleRate = types.Int64Null()
	}
	// Write-only field (service_password_wo): preserve plan/state value
	// service_username: not returned by the API, preserve plan/state value

	// Credential: extract ID from response if it's a stored credential
	if server.Credential != nil && server.Credential.Type != nil && *server.Credential.Type != "local" {
		if id := server.Credential.Id.Get(); id != nil {
			model.CredentialId = types.Int64Value(*id)
		}
	}

	// inventory_existing: read from the generic config map (config object varies by
	// pool server type). The stored representation is environment-dependent, so coerce
	// common truthy forms. When the key is absent, preserve the existing plan/state value
	// to avoid spurious drift (Morpheus often omits unchecked checkboxes).
	if server.Config != nil {
		if v, ok := server.Config["inventoryExisting"]; ok {
			switch val := v.(type) {
			case bool:
				model.InventoryExisting = types.BoolValue(val)
			case string:
				model.InventoryExisting = types.BoolValue(convert.StringToBool(ctx, val).ValueBool())
			}
		}
		// app_id (phpIPAM config.appId): mirror what the appliance holds, so a
		// server-side change surfaces as drift. The appliance refuses an empty
		// value on create and update, so an empty echo carries no information
		// and the plan/state value is kept.
		if isPhpIPAM(server) {
			if str, ok := server.Config["appId"].(string); ok && str != "" {
				model.AppId = types.StringValue(str)
			}
		}
	}
}

// resolveNetworkPoolServerTypeCode determines the pool server "type" code to send to the
// API. type_code and type_id are mutually exclusive: when type_code is set it is used
// directly; when type_id is set its code is looked up via the pool server types endpoint.
// The Morpheus API resolves the pool server type by code (NetworkPoolServerType.findByCode).
func resolveNetworkPoolServerTypeCode(
	ctx context.Context,
	client *sdk.APIClient,
	typeCode types.String,
	typeID types.Int64,
) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !typeCode.IsNull() && !typeCode.IsUnknown() && typeCode.ValueString() != "" {
		return typeCode.ValueString(), diags
	}

	if !typeID.IsNull() && !typeID.IsUnknown() {
		id := typeID.ValueInt64()
		result, httpResp, err := client.NetworksAPI.GetNetworkPoolServerType(ctx, id).Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			diags.AddError(
				"Unable to resolve network pool server type",
				fmt.Sprintf("Could not look up network pool server type with ID %d: %s", id, err),
			)

			return "", diags
		}
		if result.NetworkPoolServerType == nil || result.NetworkPoolServerType.Code == nil {
			diags.AddError(
				"Unable to resolve network pool server type",
				fmt.Sprintf("Network pool server type %d did not return a type code", id),
			)

			return "", diags
		}

		return *result.NetworkPoolServerType.Code, diags
	}

	diags.AddError(
		"Missing network pool server type",
		"Either type_code or type_id must be set to create a network pool server.",
	)

	return "", diags
}

// appIDValue returns the planned app_id, or nil when it is not set.
func appIDValue(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return v.ValueStringPointer()
}

// isPhpIPAM reports whether the read pool server is of the phpIPAM type.
func isPhpIPAM(server *sdk.GetNetworkPoolServer200ResponseNetworkPoolServer) bool {
	return server.Type != nil && server.Type.Code != nil && *server.Type.Code == typeCodePhpIPAM
}
