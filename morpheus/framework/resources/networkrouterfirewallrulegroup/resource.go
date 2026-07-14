// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_firewall_rule_group"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkRouterFirewallRuleGroupResourceSchema(ctx)
}

// Create

const createOperation = "create network router firewall rule group resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterFirewallRuleGroupModel

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

	ruleGroup := sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroup{
		Name:         plan.Name.ValueString(),
		ExternalType: plan.ExternalType.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ruleGroup.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		ruleGroup.Priority = plan.Priority.ValueInt64Pointer()
	}

	if !plan.GroupLayer.IsNull() && !plan.GroupLayer.IsUnknown() {
		ruleGroup.GroupLayer = plan.GroupLayer.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ruleGroup.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64

		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tenants := make([]sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroupTenantsInner, 0, len(ids))
		for i := range ids {
			id := ids[i]
			tenants = append(tenants, sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroupTenantsInner{Id: &id})
		}

		ruleGroup.Tenants = tenants
	}

	createReq := sdk.CreateNetworkRouterFirewallRuleGroupRequest{
		RuleGroup: &ruleGroup,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterFirewallRuleGroup(ctx, routerID).
		CreateNetworkRouterFirewallRuleGroupRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d firewall rule group POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
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
			ResourceType: "network_router_firewall_rule_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, plan)
	if notFound {
		// Unexpected: resource was just created but GET returned 404.
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("firewall rule group %d was created but GET returned 404", id),
		)
		taintState()

		return
	}

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("firewall rule group %d was created but could not be read", id),
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

const readOperation = "read network router firewall rule group resource"

// getRuleGroupAsState fetches the firewall rule group by ID and maps it into a
// NetworkRouterFirewallRuleGroupModel. The plan/prior-state model is required
// so that write-only fields (external_type, visibility, tenant_ids) can be
// preserved — none of these are returned by the single GET endpoint.
//
// Returns (state, notFound, diags). notFound is true when the API returned 404;
// the caller should remove the resource from state in that case.
func getRuleGroupAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	prior NetworkRouterFirewallRuleGroupModel,
) (NetworkRouterFirewallRuleGroupModel, bool, diag.Diagnostics) {
	var state NetworkRouterFirewallRuleGroupModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterFirewallRuleGroup(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return state, true, diags
		}

		diags.AddError(
			readOperation,
			fmt.Sprintf("firewall rule group %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, false, diags
	}

	rg := resp.RuleGroup
	if rg == nil {
		diags.AddError("API returned nil", "RuleGroup is nil in the response")

		return state, false, diags
	}

	state.Id = convert.Int64ToType(rg.Id)
	state.RouterId = prior.RouterId // path param — always preserve from prior/plan

	state.Name = convert.StrToType(rg.Name)
	state.Priority = convert.Int64ToType(rg.Priority)
	state.GroupLayer = convert.StrToType(rg.GroupLayer)
	state.ExternalId = convert.StrToType(rg.ExternalId)
	state.Status = convert.StrToType(rg.Status)

	if rg.Description.IsSet() {
		state.Description = convert.StrToType(rg.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	// Write-only fields: the single GET response does not include these.
	// Preserve the plan/prior-state values so Terraform does not see spurious
	// diffs. On import these will be null — users must re-apply to set them.
	state.ExternalType = prior.ExternalType
	state.Visibility = prior.Visibility
	state.TenantIds = prior.TenantIds

	return state, false, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var prior NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, notFound, diags := getRuleGroupAsState(ctx, prior.Id.ValueInt64(), prior.RouterId.ValueInt64(), client, prior)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}

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
	var plan NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("update network router firewall rule group resource", "failed to create client: "+err.Error())

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	// The SDK's UpdateNetworkRouterFirewallRuleGroupRequest.RuleGroup is typed
	// as map[string]interface{} — build it manually.
	ruleGroup := map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"externalType": plan.ExternalType.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ruleGroup["description"] = plan.Description.ValueString()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		ruleGroup["priority"] = plan.Priority.ValueInt64()
	}

	if !plan.GroupLayer.IsNull() && !plan.GroupLayer.IsUnknown() {
		ruleGroup["groupLayer"] = plan.GroupLayer.ValueString()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ruleGroup["visibility"] = plan.Visibility.ValueString()
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64

		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tenants := make([]map[string]interface{}, 0, len(ids))
		for _, tid := range ids {
			tenants = append(tenants, map[string]interface{}{"id": tid})
		}

		ruleGroup["tenants"] = tenants
	}

	updateReq := sdk.UpdateNetworkRouterFirewallRuleGroupRequest{
		RuleGroup: ruleGroup,
	}

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterFirewallRuleGroup(ctx, id, routerID).
		UpdateNetworkRouterFirewallRuleGroupRequest(updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update network router firewall rule group resource",
			fmt.Sprintf("firewall rule group %d PUT failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, plan)
	if notFound {
		// Unexpected: resource existed during PUT but GET returned 404.
		resp.Diagnostics.AddError(
			"update network router firewall rule group resource",
			fmt.Sprintf("firewall rule group %d was updated but GET returned 404", id),
		)

		return
	}

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
	var state NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete network router firewall rule group resource", "failed to create client: "+err.Error())

		return
	}

	id := state.Id.ValueInt64()
	routerID := state.RouterId.ValueInt64()

	tflog.Debug(ctx, fmt.Sprintf("Deleting firewall rule group %d on router %d", id, routerID))

	_, hresp, err := client.NetworksAPI.
		DeleteNetworkRouterFirewallRuleGroup(ctx, id, routerID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Firewall rule group %d already deleted (404)", id))

			return
		}

		resp.Diagnostics.AddError(
			"delete network router firewall rule group resource",
			fmt.Sprintf("firewall rule group %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp)),
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
			"import network router firewall rule group resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.group_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"provided group_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	// Seed a minimal prior state so getRuleGroupAsState can preserve write-only
	// fields. On import these fields will be null — users must re-apply.
	prior := NetworkRouterFirewallRuleGroupModel{
		RouterId:     types.Int64Value(routerID),
		ExternalType: types.StringNull(),
		Visibility:   types.StringNull(),
	}

	var setDiags diag.Diagnostics
	prior.TenantIds, setDiags = types.SetValue(types.Int64Type, []attr.Value{})
	resp.Diagnostics.Append(setDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, prior)
	if notFound {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			fmt.Sprintf("firewall rule group %d not found on router %d", id, routerID),
		)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
