// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data NetworkFirewallRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := data.Id.ValueInt64()
	networkIntegrationId := data.NetworkIntegrationId.ValueInt64()

	state, notFound, diag := getNetworkFirewallRuleAsState(ctx, id, networkIntegrationId, client)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}

	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getNetworkFirewallRuleAsState(
	ctx context.Context,
	id int64,
	networkIntegrationId int64,
	client *sdk.APIClient,
) (state NetworkFirewallRuleModel, notFound bool, diags diag.Diagnostics) {
	ruleResp, httpResp, err := client.NetworksAPI.
		GetNetworkFirewallRule(ctx, id, networkIntegrationId).Execute()
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return state, true, diags
		}

		diags.AddError(
			"error reading network firewall rule",
			fmt.Sprintf("network firewall rule %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	rule := ruleResp.Rule
	if rule == nil {
		diags.AddError("API returned nil", "Rule is nil in the response")

		return state, false, diags
	}

	state.Id = types.Int64Value(*rule.Id)
	state.NetworkIntegrationId = types.Int64Value(networkIntegrationId)
	state.Name = convert.StrToType(rule.Name)
	state.Direction = convert.StrToType(rule.Direction)
	state.Policy = convert.StrToType(rule.Policy)
	state.Enabled = convert.BoolToType(rule.Enabled)

	if rule.Priority != nil {
		state.Priority = types.Int64Value(*rule.Priority)
	} else {
		state.Priority = types.Int64Null()
	}

	// Description is a NullableString in the SDK response
	state.Description = convert.StrToType(rule.Description.Get())

	sources, srcDiags := mapSourcesFromResponse(rule.Sources)
	diags.Append(srcDiags...)
	state.Sources = sources

	destinations, dstDiags := mapDestinationsFromResponse(rule.Destinations)
	diags.Append(dstDiags...)
	state.Destinations = destinations

	scopes, scopeDiags := mapScopesFromResponse(rule.Scopes)
	diags.Append(scopeDiags...)
	state.Scopes = scopes

	application, appDiags := mapApplicationsFromResponse(rule.Applications)
	diags.Append(appDiags...)
	state.Application = application

	profile, profileDiags := mapProfilesFromResponse(rule.Profiles)
	diags.Append(profileDiags...)
	state.Profile = profile

	state.RuleGroupId = mapRuleGroupFromResponse(rule.RuleGroup)

	return state, false, diags
}

func extractStringIDs[T any](
	items []T,
	getID func(T) *string,
) (basetypes.SetValue, diag.Diagnostics) {
	ids := make([]attr.Value, 0, len(items))

	for _, item := range items {
		id := getID(item)
		if id != nil {
			ids = append(ids, types.StringValue(*id))
		}
	}

	return types.SetValue(types.StringType, ids)
}

func mapSourcesFromResponse(
	items []sdk.GetNetworkFirewallRule200ResponseRuleSourcesInner,
) (SourcesValue, diag.Diagnostics) {
	idSet, diags := extractStringIDs(items, func(i sdk.GetNetworkFirewallRule200ResponseRuleSourcesInner) *string {
		return i.Id
	})

	return SourcesValue{
		Id:    idSet,
		state: attr.ValueStateKnown,
	}, diags
}

func mapDestinationsFromResponse(
	items []sdk.GetNetworkFirewallRule200ResponseRuleDestinationsInner,
) (DestinationsValue, diag.Diagnostics) {
	idSet, diags := extractStringIDs(items, func(i sdk.GetNetworkFirewallRule200ResponseRuleDestinationsInner) *string {
		return i.Id
	})

	return DestinationsValue{
		Id:    idSet,
		state: attr.ValueStateKnown,
	}, diags
}

func mapScopesFromResponse(
	items []sdk.GetNetworkFirewallRule200ResponseRuleScopesInner,
) (ScopesValue, diag.Diagnostics) {
	idSet, diags := extractStringIDs(items, func(i sdk.GetNetworkFirewallRule200ResponseRuleScopesInner) *string {
		return i.Id
	})

	return ScopesValue{
		Id:    idSet,
		state: attr.ValueStateKnown,
	}, diags
}

func mapApplicationsFromResponse(
	applications []sdk.GetNetworkFirewallRule200ResponseRuleApplicationsInner,
) (basetypes.SetValue, diag.Diagnostics) {
	return extractStringIDs(
		applications,
		func(i sdk.GetNetworkFirewallRule200ResponseRuleApplicationsInner) *string {
			return i.Id
		},
	)
}

func mapProfilesFromResponse(
	profiles []sdk.GetNetworkFirewallRule200ResponseRuleProfilesInner,
) (basetypes.SetValue, diag.Diagnostics) {
	return extractStringIDs(
		profiles,
		func(i sdk.GetNetworkFirewallRule200ResponseRuleProfilesInner) *string {
			return i.Id
		},
	)
}

func mapRuleGroupFromResponse(
	ruleGroup *sdk.GetNetworkFirewallRule200ResponseRuleRuleGroup,
) RuleGroupIdValue {
	if ruleGroup == nil {
		return NewRuleGroupIdValueNull()
	}

	var idVal basetypes.Int64Value
	if ruleGroup.Id != nil {
		idVal = types.Int64Value(*ruleGroup.Id)
	} else {
		idVal = types.Int64Null()
	}

	var nameVal basetypes.StringValue
	if ruleGroup.Name != nil {
		nameVal = types.StringValue(*ruleGroup.Name)
	} else {
		nameVal = types.StringNull()
	}

	return RuleGroupIdValue{
		Id:    idVal,
		Name:  nameVal,
		state: attr.ValueStateKnown,
	}
}

func setValueToStringSlice(set basetypes.SetValue) []string {
	elements := set.Elements()
	result := make([]string, 0, len(elements))

	for _, e := range elements {
		if sv, ok := e.(types.String); ok {
			result = append(result, sv.ValueString())
		}
	}

	return result
}
