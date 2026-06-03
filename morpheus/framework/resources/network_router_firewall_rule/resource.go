// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_router_firewall_rule

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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
	resp.TypeName = strings.Join(
		[]string{
			req.ProviderTypeName,
			constants.SubProviderName,
			"network_router_firewall_rule",
		},
		"_",
	)
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
		enabled := plan.Enabled.ValueBool()
		rule.Enabled = &enabled
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		priority := plan.Priority.ValueInt64()
		rule.Priority = &priority
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

	id := *result.Id.Get()
	plan.Id = types.Int64Value(id)

	state, pdiags := getFirewallRuleAsState(ctx, id, routerID, client, plan)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	if rule != nil && rule.Id != nil {
		state.Id = types.Int64Value(*rule.Id)
	}

	state.RouterId = plan.RouterId

	if rule != nil && rule.Name != nil {
		state.Name = types.StringValue(*rule.Name)
	}

	if rule != nil && rule.Enabled != nil {
		state.Enabled = types.BoolValue(*rule.Enabled)
	}

	if rule != nil && rule.Priority != nil {
		state.Priority = types.Int64Value(*rule.Priority)
	} else {
		state.Priority = types.Int64Null()
	}

	if rule != nil && rule.Direction != nil {
		state.Direction = types.StringValue(*rule.Direction)
	} else {
		state.Direction = types.StringNull()
	}

	if rule != nil && rule.Policy != nil {
		state.Policy = types.StringValue(*rule.Policy)
	} else {
		state.Policy = types.StringNull()
	}

	if rule != nil && rule.Protocol.IsSet() {
		state.Protocol = types.StringValue(*rule.Protocol.Get())
	} else {
		state.Protocol = types.StringNull()
	}

	if rule != nil && rule.PortRange.IsSet() {
		state.PortRange = types.StringValue(*rule.PortRange.Get())
	} else {
		state.PortRange = types.StringNull()
	}

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
		enabled := plan.Enabled.ValueBool()
		rule.Enabled = &enabled
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		priority := plan.Priority.ValueInt64()
		rule.Priority = &priority
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
