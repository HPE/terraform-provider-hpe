package network_pool_server

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &networkPoolServerResource{}
	_ resource.ResourceWithConfigure   = &networkPoolServerResource{}
	_ resource.ResourceWithImportState = &networkPoolServerResource{}
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_pool_server"
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

	// The SDK uses a oneOf union for this request. We use InfobloxNetworkPoolServer as the
	// concrete type because it is the superset of all pool server types (Infoblox, Bluecat,
	// phpIPAM, SolarWinds). The API resolves the actual type from type_id, not from the
	// "type" field in the request body. All common fields are accepted regardless of type.
	infoblox := sdk.NewInfobloxNetworkPoolServerWithDefaults()
	infoblox.Type = "infoblox"
	infoblox.Name = plan.Name.ValueString()
	if !plan.ServiceUrl.IsNull() {
		infoblox.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	}
	if !plan.ServiceUsername.IsNull() {
		infoblox.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !plan.ServicePasswordWo.IsNull() && !plan.ServicePasswordWo.IsUnknown() {
		infoblox.ServicePassword = *sdk.NewNullableString(plan.ServicePasswordWo.ValueStringPointer())
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

	// Credential: use credential_id for stored credentials
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		cred := sdk.NewInfobloxNetworkPoolServerCredentialWithDefaults()
		cred.Type = &idStr
		infoblox.Credential = cred
	}

	serverReq := sdk.CreateNetworkPoolServerRequestNetworkPoolServer{
		InfobloxNetworkPoolServer: infoblox,
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkPoolServer(ctx).
		CreateNetworkPoolServerRequest(sdk.CreateNetworkPoolServerRequest{
			NetworkPoolServer: &serverReq,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_pool_server", plan.Name.ValueString(), err, httpResp)

		return
	}

	createServer := result.GetNetworkPoolServer()
	id := (&createServer).GetId()

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

	server := readResult.GetNetworkPoolServer()
	mapReadResponseToModel(&plan, &server)

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

	server := result.GetNetworkPoolServer()
	mapReadResponseToModel(&state, &server)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	id := plan.Id.ValueInt64()

	infobloxUpdate := sdk.NewInfobloxNetworkPoolServerUpdateWithDefaults()
	infobloxUpdate.Name = plan.Name.ValueStringPointer()
	if !plan.ServiceUrl.IsNull() {
		infobloxUpdate.ServiceUrl = *sdk.NewNullableString(plan.ServiceUrl.ValueStringPointer())
	}
	if !plan.ServiceUsername.IsNull() {
		infobloxUpdate.ServiceUsername = *sdk.NewNullableString(plan.ServiceUsername.ValueStringPointer())
	}
	if !plan.ServicePasswordWo.IsNull() && !plan.ServicePasswordWo.IsUnknown() {
		infobloxUpdate.ServicePassword = *sdk.NewNullableString(plan.ServicePasswordWo.ValueStringPointer())
	}
	if !plan.IgnoreSsl.IsNull() {
		infobloxUpdate.IgnoreSsl = plan.IgnoreSsl.ValueBoolPointer()
	}
	if !plan.Enabled.IsNull() {
		infobloxUpdate.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.NetworkFilter.IsNull() {
		infobloxUpdate.NetworkFilter = *sdk.NewNullableString(plan.NetworkFilter.ValueStringPointer())
	}
	if !plan.ZoneFilter.IsNull() {
		infobloxUpdate.ZoneFilter = *sdk.NewNullableString(plan.ZoneFilter.ValueStringPointer())
	}
	if !plan.TenantMatch.IsNull() {
		infobloxUpdate.TenantMatch = *sdk.NewNullableString(plan.TenantMatch.ValueStringPointer())
	}
	if !plan.ServiceMode.IsNull() {
		infobloxUpdate.ServiceMode = plan.ServiceMode.ValueStringPointer()
	}
	if !plan.ServiceThrottleRate.IsNull() {
		rate := plan.ServiceThrottleRate.ValueInt64()
		infobloxUpdate.ServiceThrottleRate = *sdk.NewNullableInt64(&rate)
	}

	// Credential: use credential_id for stored credentials
	if !plan.CredentialId.IsNull() {
		idStr := strconv.FormatInt(plan.CredentialId.ValueInt64(), 10)
		cred := sdk.NewInfobloxNetworkPoolServerUpdateCredentialWithDefaults()
		cred.Type = &idStr
		infobloxUpdate.Credential = cred
	}

	serverReq := sdk.InfobloxNetworkPoolServerUpdateAsUpdateNetworkPoolServerRequestNetworkPoolServer(infobloxUpdate)

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

	server := readResult.GetNetworkPoolServer()
	mapReadResponseToModel(&plan, &server)

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

	_, httpResp, err := client.NetworksAPI.DeleteNetworkPoolServer(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_pool_server", "", err, httpResp)

		return
	}
}

func (r *networkPoolServerResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(
	model *NetworkPoolServerModel,
	server *sdk.CreateNetworkPoolServer200ResponseAllOfNetworkPoolServer,
) {
	if server.Id != nil {
		model.Id = types.Int64Value(*server.Id)
	}
	if server.Name != nil {
		model.Name = types.StringValue(*server.Name)
	}
	if t := server.Type; t != nil && t.Id != nil {
		model.TypeId = types.Int64Value(*t.Id)
	}
	if server.ServiceUrl.IsSet() && server.ServiceUrl.Get() != nil {
		model.ServiceUrl = types.StringValue(*server.ServiceUrl.Get())
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

	// Credential: extract ID from response if it's a stored credential
	if server.Credential != nil && server.Credential.Type != nil && *server.Credential.Type != "local" {
		if id := server.Credential.Id.Get(); id != nil {
			model.CredentialId = types.Int64Value(*id)
		}
	}
}

func mapReadResponseToModel(
	model *NetworkPoolServerModel,
	server *sdk.GetNetworkPoolServer200ResponseNetworkPoolServer,
) {
	if server.Id != nil {
		model.Id = types.Int64Value(*server.Id)
	}
	if server.Name != nil {
		model.Name = types.StringValue(*server.Name)
	}
	if t := server.Type; t != nil && t.Id != nil {
		model.TypeId = types.Int64Value(*t.Id)
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
}
