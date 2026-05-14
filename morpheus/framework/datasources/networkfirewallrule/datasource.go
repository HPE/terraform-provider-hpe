// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkfirewallrule implements a data source for network_firewall_rule
package networkfirewallrule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                           = "read network firewall rule data source"
	ErrorNoValidSearchTerms           = `no valid search terms - an id or name is required`
	ErrorRunningPreApply              = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkFirewallRuleFound   = `no network firewall rule found`
	ErrorMultipleNetworkFirewallRules = `multiple network firewall rules were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_firewall_rule"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkFirewallRuleDataSourceSchema(ctx)
}

func firewallRuleAsState(
	ctx context.Context,
	rule *sdk.GetNetworkFirewallRule200ResponseRule,
	serverId int64,
) (NetworkFirewallRuleModel, error) {
	var convErr error

	applications, diags := convert.ToSetType(ctx, rule.Applications,
		func(app sdk.GetNetworkFirewallRule200ResponseRuleApplicationsInner) ApplicationsValue {
			v, vDiags := NewApplicationsValue(
				ApplicationsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.StrToType(app.Id),
					"name": convert.StrToType(app.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating applications value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return NetworkFirewallRuleModel{}, fmt.Errorf("error creating applications set")
	}

	destinations, diags := convert.ToSetType(ctx, rule.Destinations,
		func(d sdk.GetNetworkFirewallRule200ResponseRuleDestinationsInner) DestinationsValue {
			v, vDiags := NewDestinationsValue(
				DestinationsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.StrToType(d.Id),
					"name": convert.StrToType(d.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating destinations value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return NetworkFirewallRuleModel{}, fmt.Errorf("error creating destinations set")
	}

	sources, diags := convert.ToSetType(ctx, rule.Sources,
		func(s sdk.GetNetworkFirewallRule200ResponseRuleSourcesInner) SourcesValue {
			v, vDiags := NewSourcesValue(
				SourcesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.StrToType(s.Id),
					"name": convert.StrToType(s.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating sources value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return NetworkFirewallRuleModel{}, fmt.Errorf("error creating sources set")
	}

	scopes, diags := convert.ToSetType(ctx, rule.Scopes,
		func(s sdk.GetNetworkFirewallRule200ResponseRuleScopesInner) ScopesValue {
			v, vDiags := NewScopesValue(
				ScopesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.StrToType(s.Id),
					"name": convert.StrToType(s.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating scopes value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return NetworkFirewallRuleModel{}, fmt.Errorf("error creating scopes set")
	}

	profiles, diags := convert.ToSetType(ctx, rule.Profiles,
		func(p sdk.GetNetworkFirewallRule200ResponseRuleProfilesInner) ProfilesValue {
			v, vDiags := NewProfilesValue(
				ProfilesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":   convert.StrToType(p.Id),
					"name": convert.StrToType(p.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating profiles value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return NetworkFirewallRuleModel{}, fmt.Errorf("error creating profiles set")
	}

	appliedTargets := types.DynamicNull()
	if len(rule.AppliedTargets) > 0 {
		raw, err := json.Marshal(rule.AppliedTargets)
		if err != nil {
			return NetworkFirewallRuleModel{}, fmt.Errorf("error marshalling applied targets: %w", err)
		}

		appliedTargets = types.DynamicValue(types.StringValue(string(raw)))
	}

	ruleGroup := NewRuleGroupValueNull()
	if rule.RuleGroup.IsSet() {
		rg := rule.RuleGroup.Get()
		v, rgDiags := NewRuleGroupValue(
			RuleGroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(rg.Id),
				"name": convert.StrToType(rg.Name),
			},
		)
		if rgDiags.HasError() {
			return NetworkFirewallRuleModel{}, fmt.Errorf("error creating rule group value")
		}

		ruleGroup = v
	}

	config := types.DynamicNull()
	if rule.Config != nil {
		v, err := convert.MapToDynamic(ctx, rule.Config)
		if err != nil {
			return NetworkFirewallRuleModel{}, fmt.Errorf("error creating config value: %w", err)
		}

		config = v
	}

	return NetworkFirewallRuleModel{
		Applications:    applications,
		AppliedTargets:  appliedTargets,
		Config:          config,
		DestinationType: convert.StrToType(rule.DestinationType),
		Destinations:    destinations,
		Direction:       convert.StrToType(rule.Direction),
		Enabled:         convert.BoolToType(rule.Enabled),
		GroupName:       convert.StrToType(rule.GroupName),
		Id:              convert.Int64ToType(rule.Id),
		Name:            convert.StrToType(rule.Name),
		Policy:          convert.StrToType(rule.Policy),
		Priority:        convert.Int64ToType(rule.Priority),
		Profiles:        profiles,
		RuleGroup:       ruleGroup,
		Scopes:          scopes,
		ServerId:        types.Int64Value(serverId),
		SourceType:      convert.StrToType(rule.SourceType),
		Sources:         sources,
	}, nil
}

func getFirewallRuleByID(
	ctx context.Context,
	id int64,
	serverId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkFirewallRule200ResponseRule, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkFirewallRule(
		ctx, id, serverId,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network firewall rule %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	rule := r.GetRule()

	return &rule, nil
}

// listRuleSummary is a lightweight struct for unmarshalling the list endpoint
// response, which returns rules as interface{}.
type listRuleSummary struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func getFirewallRuleByName(
	ctx context.Context,
	name string,
	serverId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkFirewallRule200ResponseRule, error) {
	// Phrase is used because the API does not expose an exact Name filter.
	// The subsequent loop performs exact-match filtering on the results.
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkFirewallRules(
		ctx, serverId,
	).Phrase(name).Max(10000).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network firewall rules: %s",
			providererrors.ErrMsg(err, hresp),
		)
	}

	raw := rs.GetRules()

	rawJSON, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("error marshalling rules list: %w", err)
	}

	var items []listRuleSummary
	if err := json.Unmarshal(rawJSON, &items); err != nil {
		return nil, fmt.Errorf("error unmarshalling rules list: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("%s matching name %q", ErrorNoNetworkFirewallRuleFound, name)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}

		if items[i].Id != nil {
			matchedIDs = append(matchedIDs, *items[i].Id)
		}
	}

	if len(matchedIDs) == 0 {
		return nil, fmt.Errorf("%s matching name %q", ErrorNoNetworkFirewallRuleFound, name)
	} else if len(matchedIDs) > 1 {
		return nil, fmt.Errorf("%s matching name %q", ErrorMultipleNetworkFirewallRules, name)
	}

	return getFirewallRuleByID(ctx, matchedIDs[0], serverId, apiClient)
}

func getFirewallRule(
	ctx context.Context,
	config *NetworkFirewallRuleModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkFirewallRule200ResponseRule, error) {
	serverId := config.ServerId.ValueInt64()

	if !config.Id.IsNull() {
		return getFirewallRuleByID(ctx, config.Id.ValueInt64(), serverId, apiClient)
	} else if !config.Name.IsNull() {
		return getFirewallRuleByName(ctx, config.Name.ValueString(), serverId, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkFirewallRuleModel

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

	serverId := config.ServerId.ValueInt64()
	state, err := firewallRuleAsState(ctx, rule, serverId)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
