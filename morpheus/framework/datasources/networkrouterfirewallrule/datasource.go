// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouterfirewallrule implements a data source for network_router_firewall_rule
package networkrouterfirewallrule

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                                 = "read network router firewall rule data source"
	ErrorNoValidSearchTerms                 = `no valid search terms - an id or name is required`
	ErrorNoNetworkRouterFirewallRuleFound   = `no network router firewall rule found`
	ErrorMultipleNetworkRouterFirewallRules = `multiple network router firewall rules were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_firewall_rule"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterFirewallRuleDataSourceSchema(ctx)
}

func firewallRuleAsState(
	rule *sdk.GetNetworkRouterFirewallRule200ResponseRule,
	routerID int64,
) NetworkRouterFirewallRuleModel {
	state := NetworkRouterFirewallRuleModel{
		Id:              convert.Int64ToType(rule.Id),
		RouterId:        types.Int64Value(routerID),
		Name:            convert.StrToType(rule.Name),
		Enabled:         convert.BoolToType(rule.Enabled),
		Priority:        convert.Int64ToType(rule.Priority),
		Direction:       convert.StrToType(rule.Direction),
		Policy:          convert.StrToType(rule.Policy),
		GroupName:       convert.StrToType(rule.GroupName),
		RuleType:        convert.StrToType(rule.RuleType),
		SourceType:      convert.StrToType(rule.SourceType),
		DestinationType: convert.StrToType(rule.DestinationType),
		ApplicationType: convert.StrToType(rule.ApplicationType),
	}

	// Protocol — nullable
	if rule.Protocol.IsSet() {
		state.Protocol = convert.StrToType(rule.Protocol.Get())
	} else {
		state.Protocol = types.StringNull()
	}

	// Application — nullable
	if rule.Application.IsSet() {
		state.Application = convert.StrToType(rule.Application.Get())
	} else {
		state.Application = types.StringNull()
	}

	// Code — nullable
	if rule.Code.IsSet() {
		state.Code = convert.StrToType(rule.Code.Get())
	} else {
		state.Code = types.StringNull()
	}

	// PortRange — nullable
	if rule.PortRange.IsSet() {
		state.PortRange = convert.StrToType(rule.PortRange.Get())
	} else {
		state.PortRange = types.StringNull()
	}

	// SourcePortRange — nullable
	if rule.SourcePortRange.IsSet() {
		state.SourcePortRange = convert.StrToType(rule.SourcePortRange.Get())
	} else {
		state.SourcePortRange = types.StringNull()
	}

	// DestinationPortRange — nullable
	if rule.DestinationPortRange.IsSet() {
		state.DestinationPortRange = convert.StrToType(rule.DestinationPortRange.Get())
	} else {
		state.DestinationPortRange = types.StringNull()
	}

	// SourceGroup — nullable
	if rule.SourceGroup.IsSet() {
		state.SourceGroup = convert.StrToType(rule.SourceGroup.Get())
	} else {
		state.SourceGroup = types.StringNull()
	}

	// SourceTier — nullable
	if rule.SourceTier.IsSet() {
		state.SourceTier = convert.StrToType(rule.SourceTier.Get())
	} else {
		state.SourceTier = types.StringNull()
	}

	return state
}

func getFirewallRuleByID(
	ctx context.Context,
	id int64,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRule200ResponseRule, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterFirewallRule(
		ctx, id, routerID,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router firewall rule %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.Rule, nil
}

func getFirewallRuleByName(
	ctx context.Context,
	name string,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRule200ResponseRule, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRoutersFirewallRules(
		ctx, routerID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router firewall rules with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.Rules
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterFirewallRuleFound)
	}

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
		return nil, errors.New(ErrorNoNetworkRouterFirewallRuleFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouterFirewallRules)
	}

	return getFirewallRuleByID(ctx, matchedIDs[0], routerID, apiClient)
}

func getFirewallRule(
	ctx context.Context,
	config *NetworkRouterFirewallRuleModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterFirewallRule200ResponseRule, error) {
	routerID := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getFirewallRuleByID(ctx, config.Id.ValueInt64(), routerID, apiClient)
	} else if !config.Name.IsNull() {
		return getFirewallRuleByName(ctx, config.Name.ValueString(), routerID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterFirewallRuleModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	rule, err := getFirewallRule(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	routerID := config.RouterId.ValueInt64()
	state := firewallRuleAsState(rule, routerID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
