// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouterfirewallrulegroup implements a data source for
// network router firewall rule groups.
package networkrouterfirewallrulegroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                                      = "read network router firewall rule group data source"
	ErrorNoValidSearchTerms                      = `no valid search terms - an id or name is required`
	ErrorNoNetworkRouterFirewallRuleGroupFound   = `no network router firewall rule group found`
	ErrorMultipleNetworkRouterFirewallRuleGroups = `multiple network router firewall rule groups were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_firewall_rule_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterFirewallRuleGroupDataSourceSchema(ctx)
	resp.Schema.Description = "Provides a network router firewall rule group data source."
	resp.Schema.MarkdownDescription = "Provides a network router firewall rule group data source."
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	rg, err := getFirewallRuleGroup(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	if rg == nil {
		resp.Diagnostics.AddError(summary, ErrorNoNetworkRouterFirewallRuleGroupFound)

		return
	}

	routerID := config.RouterId.ValueInt64()

	state, diags := ruleGroupAsState(ctx, rg, routerID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getFirewallRuleGroup(
	ctx context.Context,
	config *NetworkRouterFirewallRuleGroupModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup, error) {
	routerID := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getRuleGroupByID(ctx, config.Id.ValueInt64(), routerID, apiClient)
	} else if !config.Name.IsNull() {
		return getRuleGroupByName(ctx, config.Name.ValueString(), routerID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

func getRuleGroupByID(
	ctx context.Context,
	id int64,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup, error) {
	// SDK arg order: group id first, router id second.
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterFirewallRuleGroup(ctx, id, routerID).Execute()
	if hresp != nil && hresp.Body != nil {
		defer hresp.Body.Close()
	}

	if r == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router firewall rule group %d: %s",
			id, errfmt.ErrMsg(err, hresp),
		)
	}

	return r.RuleGroup, nil
}

func getRuleGroupByName(
	ctx context.Context,
	name string,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRouterFirewallRuleGroups(ctx, routerID).Execute()
	if hresp != nil && hresp.Body != nil {
		defer hresp.Body.Close()
	}

	if rs == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router firewall rule groups (router_id=%d): %s",
			routerID, errfmt.ErrMsg(err, hresp),
		)
	}

	items := rs.RuleGroups

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *items[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterFirewallRuleGroupFound)
	}

	if len(matchedIDs) > 1 {
		ids := make([]string, 0, len(matchedIDs))
		for _, id := range matchedIDs {
			ids = append(ids, fmt.Sprintf("%d", id))
		}

		return nil, fmt.Errorf(
			"%s with name %s. IDs: %s. Please specify an ID instead",
			ErrorMultipleNetworkRouterFirewallRuleGroups,
			name,
			strings.Join(ids, ", "),
		)
	}

	return getRuleGroupByID(ctx, matchedIDs[0], routerID, apiClient)
}

func ruleGroupAsState(
	ctx context.Context,
	rg *sdk.GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup,
	routerID int64,
) (NetworkRouterFirewallRuleGroupModel, diag.Diagnostics) {
	var allDiags diag.Diagnostics

	// Map tenants slice to a set of {id, name} objects.
	tenantObjs := make([]attr.Value, 0, len(rg.Tenants))

	for j := range rg.Tenants {
		t := rg.Tenants[j]

		tv, tDiags := NewTenantsValue(
			TenantsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(t.Id),
				"name": convert.StrToType(t.Name),
			},
		)
		allDiags.Append(tDiags...)

		tenantObjs = append(tenantObjs, tv)
	}

	if allDiags.HasError() {
		return NetworkRouterFirewallRuleGroupModel{}, allDiags
	}

	var tenantsSet types.Set

	if len(tenantObjs) == 0 {
		tenantsSet = types.SetNull(TenantsValue{}.Type(ctx))
	} else {
		var setDiags diag.Diagnostics

		tenantsSet, setDiags = types.SetValue(TenantsValue{}.Type(ctx), tenantObjs)
		allDiags.Append(setDiags...)
	}

	if allDiags.HasError() {
		return NetworkRouterFirewallRuleGroupModel{}, allDiags
	}

	return NetworkRouterFirewallRuleGroupModel{
		Description: convert.StrToType(rg.Description.Get()),
		ExternalId:  convert.StrToType(rg.ExternalId),
		GroupLayer:  convert.StrToType(rg.GroupLayer),
		Id:          convert.Int64ToType(rg.Id),
		Name:        convert.StrToType(rg.Name),
		Priority:    convert.Int64ToType(rg.Priority),
		RouterId:    types.Int64Value(routerID),
		Status:      convert.StrToType(rg.Status),
		Tenants:     tenantsSet,
		Visibility:  convert.StrToType(rg.Visibility),
	}, allDiags
}
