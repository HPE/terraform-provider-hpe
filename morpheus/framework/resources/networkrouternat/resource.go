// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_nat"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkRouterNatResourceSchema(ctx)
}

// Create

const createOperation = "create network router nat resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterNatModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only attributes (action, firewall, service) are nullified in the
	// plan by the framework. Read them from the raw config instead.
	var config NetworkRouterNatModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(createOperation, "failed to create client: "+err.Error())

		return
	}

	routerID := plan.RouterId.ValueInt64()

	// action, firewall, and service are write-only — sourced from config.
	// The Morpheus OptionType declares firewall as required with a UI-only
	// default of MATCH_INTERNAL_ADDRESS; the provider must supply it when
	// the practitioner omits it.
	natConfig := buildCreateNatConfig(&config)

	nat := sdk.CreateNetworkRouterNatRequestNetworkRouterNAT{
		Name:   plan.Name.ValueString(),
		Config: natConfig,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		nat.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		nat.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.SourceNetwork.IsNull() && !plan.SourceNetwork.IsUnknown() {
		nat.SourceNetwork = plan.SourceNetwork.ValueStringPointer()
	}
	if !plan.DestinationNetwork.IsNull() && !plan.DestinationNetwork.IsUnknown() {
		nat.DestinationNetwork = plan.DestinationNetwork.ValueStringPointer()
	}
	if !plan.TranslatedNetwork.IsNull() && !plan.TranslatedNetwork.IsUnknown() {
		nat.TranslatedNetwork = plan.TranslatedNetwork.ValueStringPointer()
	}
	if !plan.TranslatedPorts.IsNull() && !plan.TranslatedPorts.IsUnknown() {
		nat.TranslatedPorts = plan.TranslatedPorts.ValueStringPointer()
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		nat.Priority = plan.Priority.ValueInt64Pointer()
	}
	// protocol is deprecated (superseded by service) but retained for backward
	// compatibility; still sent so existing configurations keep working.
	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		nat.Protocol = plan.Protocol.ValueStringPointer()
	}

	createReq := sdk.CreateNetworkRouterNatRequest{
		NetworkRouterNAT: &nat,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterNat(ctx, routerID).
		CreateNetworkRouterNatRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d NAT POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if !result.Id.IsSet() || result.Id.Get() == nil {
		resp.Diagnostics.AddError("API returned nil", "ID is nil in the response")

		return
	}

	id := *result.Id.Get()
	plan.Id = types.Int64Value(id)

	taintState := func() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_router_nat",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getNatAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("NAT %d was created but could not be read", id),
		)
		taintState()

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		taintState()
	}
}

// Read

const readOperation = "read network router nat resource"

func getNatAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterNatModel,
) (NetworkRouterNatModel, diag.Diagnostics) {
	var state NetworkRouterNatModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterNat(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("NAT %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, diags
	}

	nat := resp.NetworkRouterNAT
	if nat == nil {
		diags.AddError("API returned nil", "NetworkRouterNAT is nil in the response")

		return state, diags
	}

	if nat.Id != nil {
		state.Id = types.Int64Value(int64(*nat.Id))
	}

	state.RouterId = plan.RouterId

	state.ExternalId = convert.StrToType(nat.ExternalId)

	state.Name = convert.StrToType(nat.Name)

	// action, firewall, and service are write-only — always null in state.
	state.Action = types.StringNull()
	state.Firewall = types.StringNull()
	state.Service = types.StringNull()

	// Carry forward the version companions from the plan/prior state.
	state.ActionVersion = plan.ActionVersion
	state.FirewallVersion = plan.FirewallVersion
	state.ServiceVersion = plan.ServiceVersion

	state.Description = convert.StrToType(nat.Description)

	state.Enabled = convert.BoolToType(nat.Enabled)

	state.SourceNetwork = convert.StrToType(nat.SourceNetwork)

	state.DestinationNetwork = convert.StrToType(nat.DestinationNetwork.Get())

	state.TranslatedNetwork = convert.StrToType(nat.TranslatedNetwork)

	state.TranslatedPorts = convert.StrToType(nat.TranslatedPorts.Get())

	if nat.Priority != nil {
		state.Priority = types.Int64Value(int64(*nat.Priority))
	} else {
		state.Priority = types.Int64Null()
	}

	// protocol is deprecated (superseded by service) and the API no longer
	// persists it, so it is omitted from the response. Fall back to the plan
	// value when the API omits it (matching action/firewall/service) so the
	// configured value round-trips and does not produce an inconsistent result
	// after apply.
	if p := nat.Protocol.Get(); p != nil {
		state.Protocol = types.StringValue(*p)
	} else if !plan.Protocol.IsUnknown() {
		state.Protocol = plan.Protocol
	} else {
		state.Protocol = types.StringNull()
	}

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterNatModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, diags := getNatAsState(ctx, plan.Id.ValueInt64(), plan.RouterId.ValueInt64(), client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkRouterNatModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only attributes (action, firewall, service) are nullified in the
	// plan by the framework. Read them from the raw config instead.
	var config NetworkRouterNatModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("update network router nat resource", "failed to create client: "+err.Error())

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	natConfig := buildUpdateNatConfig(&config)

	nat := sdk.UpdateNetworkRouterNatRequestNetworkRouterNAT{
		Name:   plan.Name.ValueStringPointer(),
		Config: natConfig,
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		nat.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		nat.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.SourceNetwork.IsNull() && !plan.SourceNetwork.IsUnknown() {
		nat.SourceNetwork = plan.SourceNetwork.ValueStringPointer()
	}
	if !plan.DestinationNetwork.IsNull() && !plan.DestinationNetwork.IsUnknown() {
		nat.DestinationNetwork = plan.DestinationNetwork.ValueStringPointer()
	}
	if !plan.TranslatedNetwork.IsNull() && !plan.TranslatedNetwork.IsUnknown() {
		nat.TranslatedNetwork = plan.TranslatedNetwork.ValueStringPointer()
	}
	if !plan.TranslatedPorts.IsNull() && !plan.TranslatedPorts.IsUnknown() {
		nat.TranslatedPorts = plan.TranslatedPorts.ValueStringPointer()
	}
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		nat.Priority = plan.Priority.ValueInt64Pointer()
	}
	// protocol is deprecated (superseded by service) but retained for backward
	// compatibility; still sent so existing configurations keep working.
	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		nat.Protocol = plan.Protocol.ValueStringPointer()
	}

	updateReq := sdk.UpdateNetworkRouterNatRequest{
		NetworkRouterNAT: &nat,
	}

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterNat(ctx, id, routerID).
		UpdateNetworkRouterNatRequest(updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network router nat resource",
			fmt.Sprintf("NAT %d PUT failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, diags := getNatAsState(ctx, id, routerID, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkRouterNatModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete network router nat resource", "failed to create client: "+err.Error())

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting NAT %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterNat(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("NAT %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router nat resource",
			fmt.Sprintf("NAT %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)
	}
}

// Import

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network router nat resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.nat_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router nat resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router nat resource",
			"provided nat_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("router_id"), routerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// defaultFirewallMatch is the firewall value applied on create when the
// practitioner supplies none.
//
// firewall is write-only, and write-only attributes cannot carry a schema
// default, so the default has to be applied here. The Morpheus OptionType is
// required, and applyOptionTypes does not honour the seed's defaultValue -
// that is a UI form default only - so omitting it would send null to a
// required field.
const defaultFirewallMatch = "MATCH_INTERNAL_ADDRESS"

// buildCreateNatConfig assembles the NAT config sent on create.
//
// It reads from the *config*, never the plan: action, firewall and service are
// write-only, so the framework nullifies them in the plan. Sourcing them from
// the plan compiles cleanly and silently sends nulls.
func buildCreateNatConfig(
	config *NetworkRouterNatModel,
) sdk.CreateNetworkRouterNatRequestNetworkRouterNATConfig {
	natConfig := sdk.CreateNetworkRouterNatRequestNetworkRouterNATConfig{
		Action: config.Action.ValueString(),
	}

	if !config.Firewall.IsNull() && !config.Firewall.IsUnknown() {
		natConfig.Firewall = config.Firewall.ValueStringPointer()
	} else {
		firewall := defaultFirewallMatch
		natConfig.Firewall = &firewall
	}

	if !config.Service.IsNull() && !config.Service.IsUnknown() {
		natConfig.Service = config.Service.ValueStringPointer()
	}

	return natConfig
}

// buildUpdateNatConfig assembles the NAT config sent on update.
//
// Unlike create it applies no default. The fields are omitempty pointers and
// the controller merges the payload over the current config, so leaving a key
// out preserves whatever the router already has. Sending a default here would
// overwrite a real NSX-T value with a guess.
func buildUpdateNatConfig(
	config *NetworkRouterNatModel,
) *sdk.UpdateNetworkRouterNatRequestNetworkRouterNATConfig {
	natConfig := &sdk.UpdateNetworkRouterNatRequestNetworkRouterNATConfig{
		Action: config.Action.ValueStringPointer(),
	}

	if !config.Firewall.IsNull() && !config.Firewall.IsUnknown() {
		natConfig.Firewall = config.Firewall.ValueStringPointer()
	}

	if !config.Service.IsNull() && !config.Service.IsUnknown() {
		natConfig.Service = config.Service.ValueStringPointer()
	}

	return natConfig
}
