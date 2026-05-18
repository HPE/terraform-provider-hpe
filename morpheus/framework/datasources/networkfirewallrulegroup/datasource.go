// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkfirewallrulegroup implements a data source for
// network firewall rule groups.
package networkfirewallrulegroup

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

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary            = "read network firewall rule group data source"
	ErrorNoValidSearch = "no valid search terms - an id or name is required"
	ErrorNotFound      = "no network firewall rule group found"
	ErrorMultipleFound = "multiple network firewall rule groups found"
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_firewall_rule_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkFirewallRuleGroupDataSourceSchema(ctx)
	resp.Schema.Description = "Provides a network firewall rule group data source."
	resp.Schema.MarkdownDescription = "Provides a network firewall rule group data source."
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkFirewallRuleGroupModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	ruleGroup, err := getRuleGroup(ctx, &config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state, stDiags := ruleGroupAsState(ctx, ruleGroup, config.NetworkIntegrationId.ValueInt64())
	resp.Diagnostics.Append(stDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getRuleGroup(
	ctx context.Context,
	config *NetworkFirewallRuleGroupModel,
	client *sdk.APIClient,
) (*sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroup, error) {
	integrationID := config.NetworkIntegrationId.ValueInt64()

	if !config.Id.IsNull() {
		return getRuleGroupByID(ctx, config.Id.ValueInt64(), integrationID, client)
	} else if !config.Name.IsNull() {
		return getRuleGroupByName(ctx, config.Name.ValueString(), integrationID, client)
	}

	return nil, errors.New(ErrorNoValidSearch)
}

func getRuleGroupByID(
	ctx context.Context,
	id int64, integrationID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroup, error) {
	r, hresp, err := client.NetworksAPI.GetNetworkFirewallRuleGroup(ctx, id, integrationID).Execute()
	if hresp != nil && hresp.Body != nil {
		defer hresp.Body.Close()
	}

	if r == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for firewall rule group %d: %s",
			id, errfmt.ErrMsg(err, hresp),
		)
	}

	rg := r.GetRuleGroup()

	return &rg, nil
}

func getRuleGroupByName(
	ctx context.Context,
	name string,
	integrationID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroup, error) {
	r, hresp, err := client.NetworksAPI.GetNetworkFirewallRuleGroups(ctx, integrationID).Execute()
	if hresp != nil && hresp.Body != nil {
		defer hresp.Body.Close()
	}

	if r == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for firewall rule groups on integration %d: %s",
			integrationID, errfmt.ErrMsg(err, hresp),
		)
	}

	ruleGroups := r.GetRuleGroups()

	var matchedIDs []int64

	for _, rg := range ruleGroups {
		if rg.GetName() != name {
			continue
		}

		matchedIDs = append(matchedIDs, rg.GetId())
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNotFound)
	}

	if len(matchedIDs) > 1 {
		ids := make([]string, 0, len(matchedIDs))
		for _, id := range matchedIDs {
			ids = append(ids, fmt.Sprintf("%d", id))
		}

		return nil, fmt.Errorf(
			"%s with name %s. IDs: %s. Please specify an ID instead",
			ErrorMultipleFound,
			name,
			strings.Join(ids, ", "),
		)
	}

	return getRuleGroupByID(ctx, matchedIDs[0], integrationID, client)
}

func ruleGroupAsState(
	ctx context.Context,
	rg *sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroup,
	integrationID int64,
) (NetworkFirewallRuleGroupModel, diag.Diagnostics) {
	var allDiags diag.Diagnostics

	rules, diags := convert.ToSetType(ctx, rg.Rules,
		func(
			r sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner,
		) RulesValue {
			return mapRule(ctx, r, &allDiags)
		})
	allDiags.Append(diags...)

	if allDiags.HasError() {
		return NetworkFirewallRuleGroupModel{}, allDiags
	}

	return NetworkFirewallRuleGroupModel{
		Description:          convert.StrToType(rg.Description.Get()),
		GroupLayer:           convert.StrToType(rg.GroupLayer),
		Id:                   convert.Int64ToType(rg.Id),
		Name:                 convert.StrToType(rg.Name),
		Priority:             convert.Int64ToType(rg.Priority),
		Rules:                rules,
		NetworkIntegrationId: types.Int64Value(integrationID),
	}, allDiags
}

func mapRule(
	ctx context.Context,
	r sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner,
	diags *diag.Diagnostics,
) RulesValue {
	applications := mapIDNameSet(ctx, r.Applications, diags,
		func(a sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerApplicationsInner,
		) (string, string) {
			return a.GetId(), a.GetName()
		},
		NewApplicationsValue,
	)

	appliedTargets := mapEmptyObjectSet[AppliedTargetsValue](ctx, len(r.AppliedTargets), diags)

	configObj, cfgDiags := types.ObjectValue(
		ConfigValue{}.AttributeTypes(ctx),
		map[string]attr.Value{},
	)
	diags.Append(cfgDiags...)

	destinations := mapIDNameSet(ctx, r.Destinations, diags,
		func(d sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerDestinationsInner,
		) (string, string) {
			return d.GetId(), d.GetName()
		},
		NewDestinationsValue,
	)

	profiles := mapIDNameSet(ctx, r.Profiles, diags,
		func(p sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerProfilesInner,
		) (string, string) {
			return p.GetId(), p.GetName()
		},
		NewProfilesValue,
	)

	ruleGroupVal := mapRuleGroupInner(ctx, r.RuleGroup, diags)

	scopes := mapIDNameSet(ctx, r.Scopes, diags,
		func(s sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerScopesInner,
		) (string, string) {
			return s.GetId(), s.GetName()
		},
		NewScopesValue,
	)

	sources := mapIDNameSet(ctx, r.Sources, diags,
		func(s sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerSourcesInner,
		) (string, string) {
			return s.GetId(), s.GetName()
		},
		NewSourcesValue,
	)

	v, ruleDiags := NewRulesValue(
		RulesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"applications":     applications,
			"applied_targets":  appliedTargets,
			"config":           configObj,
			"destination_type": convert.StrToType(r.DestinationType),
			"destinations":     destinations,
			"direction":        convert.StrToType(r.Direction),
			"enabled":          convert.BoolToType(r.Enabled),
			"group_name":       convert.StrToType(r.GroupName),
			"id":               convert.Int64ToType(r.Id),
			"name":             convert.StrToType(r.Name),
			"policy":           convert.StrToType(r.Policy),
			"priority":         convert.Int64ToType(r.Priority),
			"profiles":         profiles,
			"rule_group":       ruleGroupVal,
			"scopes":           scopes,
			"source_type":      convert.StrToType(r.SourceType),
			"sources":          sources,
		},
	)
	diags.Append(ruleDiags...)

	return v
}

// idNameExtractor extracts id and name strings from a model.
type idNameExtractor[T any] func(T) (id, name string)

func mapIDNameSet[T any, V interface {
	attr.Value
	AttributeTypes(context.Context) map[string]attr.Type
}](
	ctx context.Context,
	items []T,
	diags *diag.Diagnostics,
	extract idNameExtractor[T],
	construct func(map[string]attr.Type, map[string]attr.Value) (V, diag.Diagnostics),
) types.Set {
	var obj V
	attrTypes := obj.AttributeTypes(ctx)

	if len(items) == 0 {
		return types.SetNull(obj.Type(ctx))
	}

	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		id, name := extract(item)

		v, d := construct(attrTypes, map[string]attr.Value{
			"id":   types.StringValue(id),
			"name": types.StringValue(name),
		})
		diags.Append(d...)

		values = append(values, v)
	}

	set, d := types.SetValue(obj.Type(ctx), values)
	diags.Append(d...)

	return set
}

func mapEmptyObjectSet[V interface {
	attr.Value
	AttributeTypes(context.Context) map[string]attr.Type
}](ctx context.Context, count int, diags *diag.Diagnostics) types.Set {
	var obj V

	if count == 0 {
		return types.SetNull(obj.Type(ctx))
	}

	values := make([]attr.Value, 0, count)
	for range count {
		objVal, d := types.ObjectValue(obj.AttributeTypes(ctx), map[string]attr.Value{})
		diags.Append(d...)

		values = append(values, objVal)
	}

	set, d := types.SetValue(obj.Type(ctx), values)
	diags.Append(d...)

	return set
}

func mapRuleGroupInner(
	ctx context.Context,
	rg *sdk.GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerRuleGroup,
	diags *diag.Diagnostics,
) types.Object {
	attrTypes := RuleGroupValue{}.AttributeTypes(ctx)

	if rg == nil {
		return types.ObjectNull(attrTypes)
	}

	obj, d := types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":   convert.Int64ToType(rg.Id),
		"name": convert.StrToType(rg.Name),
	})
	diags.Append(d...)

	return obj
}
