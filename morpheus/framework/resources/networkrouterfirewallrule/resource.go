// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule

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
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_firewall_rule"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkRouterFirewallRuleResourceSchema(ctx)
}

// Create

const createOperation = "create network router firewall rule resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterFirewallRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(createOperation, "failed to create client: "+err.Error())

		return
	}

	routerID := plan.RouterId.ValueInt64()

	rule := sdk.CreateNetworkRouterFirewallRuleRequestRule{}
	rule.Name = plan.Name.ValueString()

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		rule.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		rule.Priority = plan.Priority.ValueInt64Pointer()
	}

	// policy, direction, protocol and port_range are top-level rule fields;
	// parent_id is nested under config.parentId and is required by NSX-T (its
	// absence causes a server-side nil pointer). Guard the optional fields so an
	// unknown value is omitted rather than sent as "".
	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		rule.Policy = plan.Policy.ValueStringPointer()
	}
	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		rule.Direction = plan.Direction.ValueStringPointer()
	}
	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		rule.Protocol = plan.Protocol.ValueStringPointer()
	}
	if !plan.PortRange.IsNull() && !plan.PortRange.IsUnknown() {
		rule.PortRange = plan.PortRange.ValueStringPointer()
	}
	if !plan.SourceType.IsNull() && !plan.SourceType.IsUnknown() {
		rule.SourceType = plan.SourceType.ValueStringPointer()
	}
	if !plan.DestinationType.IsNull() && !plan.DestinationType.IsUnknown() {
		rule.DestinationType = plan.DestinationType.ValueStringPointer()
	}
	if !plan.Application.IsNull() && !plan.Application.IsUnknown() {
		rule.Application = plan.Application.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	rule.Config = &sdk.CreateNetworkRouterFirewallRuleRequestRuleConfig{
		ParentId: plan.ParentId.ValueStringPointer(),
	}

	createReq := sdk.CreateNetworkRouterFirewallRuleRequest{
		Rule: &rule,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterFirewallRule(ctx, routerID).
		CreateNetworkRouterFirewallRuleRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d firewall rule POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
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
			ResourceType: "network_router_firewall_rule",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, pdiags := getFirewallRuleAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("firewall rule %d was created but could not be read", id),
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

const readOperation = "read network router firewall rule resource"

func getFirewallRuleAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	plan NetworkRouterFirewallRuleModel,
) (NetworkRouterFirewallRuleModel, diag.Diagnostics) {
	var state NetworkRouterFirewallRuleModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterFirewallRule(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			readOperation,
			fmt.Sprintf("firewall rule %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, diags
	}

	rule := resp.Rule
	if rule == nil {
		diags.AddError("API returned nil", "Rule is nil in the response")

		return state, diags
	}

	state.Id = convert.Int64ToType(rule.Id)

	state.RouterId = plan.RouterId

	state.Name = convert.StrToType(rule.Name)

	state.Enabled = convert.BoolToType(rule.Enabled)

	state.Priority = convert.Int64ToType(rule.Priority)

	state.Direction = convert.StrToType(rule.Direction)
	state.Policy = convert.StrToType(rule.Policy)
	state.Protocol = convert.StrToType(rule.Protocol.Get())
	state.PortRange = convert.StrToType(rule.PortRange.Get())
	state.SourceType = convert.StrToType(rule.SourceType)
	state.DestinationType = convert.StrToType(rule.DestinationType)
	state.Application = convert.StrToType(rule.Application.Get())

	// description is a write-only input: the API accepts it on create/update but
	// does not return it, so preserve the configured/prior value.
	state.Description = plan.Description

	// parent_id is a create-time (RequiresReplace) input that the response does
	// not echo back cleanly, so preserve the configured value.
	state.ParentId = plan.ParentId

	return state, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan NetworkRouterFirewallRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, diags := getFirewallRuleAsState(ctx, plan.Id.ValueInt64(), plan.RouterId.ValueInt64(), client, plan)
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
	var plan NetworkRouterFirewallRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("update network router firewall rule resource", "failed to create client: "+err.Error())

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	rule := sdk.UpdateNetworkRouterFirewallRuleRequestRule{}
	rule.Name = plan.Name.ValueString()

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		rule.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		rule.Priority = plan.Priority.ValueInt64Pointer()
	}

	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		rule.Policy = plan.Policy.ValueStringPointer()
	}
	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		rule.Direction = plan.Direction.ValueStringPointer()
	}
	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		rule.Protocol = plan.Protocol.ValueStringPointer()
	}
	if !plan.PortRange.IsNull() && !plan.PortRange.IsUnknown() {
		rule.PortRange = plan.PortRange.ValueStringPointer()
	}
	if !plan.SourceType.IsNull() && !plan.SourceType.IsUnknown() {
		rule.SourceType = plan.SourceType.ValueStringPointer()
	}
	if !plan.DestinationType.IsNull() && !plan.DestinationType.IsUnknown() {
		rule.DestinationType = plan.DestinationType.ValueStringPointer()
	}
	if !plan.Application.IsNull() && !plan.Application.IsUnknown() {
		rule.Application = plan.Application.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	rule.Config = &sdk.UpdateNetworkRouterFirewallRuleRequestRuleConfig{
		ParentId: plan.ParentId.ValueStringPointer(),
	}

	updateReq := sdk.UpdateNetworkRouterFirewallRuleRequest{
		Rule: &rule,
	}

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterFirewallRule(ctx, id, routerID).
		UpdateNetworkRouterFirewallRuleRequest(updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network router firewall rule resource",
			fmt.Sprintf("firewall rule %d PUT failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, diags := getFirewallRuleAsState(ctx, id, routerID, client, plan)
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
	var state NetworkRouterFirewallRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete network router firewall rule resource", "failed to create client: "+err.Error())

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting firewall rule %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterFirewallRule(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Firewall rule %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router firewall rule resource",
			fmt.Sprintf("firewall rule %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp)),
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
			"import network router firewall rule resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.rule_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule resource",
			"provided rule_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("router_id"), routerID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
