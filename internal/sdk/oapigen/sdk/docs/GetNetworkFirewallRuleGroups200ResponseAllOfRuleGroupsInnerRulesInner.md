# GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Direction** | Pointer to **string** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**DestinationType** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Policy** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **int64** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**RuleGroup** | Pointer to [**NullableGetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerRuleGroup**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerRuleGroup.md) |  | [optional] 
**GroupName** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Sources** | Pointer to [**[]GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerSourcesInner**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerSourcesInner.md) |  | [optional] 
**Destinations** | Pointer to [**[]GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerDestinationsInner**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerDestinationsInner.md) |  | [optional] 
**Applications** | Pointer to [**[]GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerApplicationsInner**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerApplicationsInner.md) |  | [optional] 
**Scopes** | Pointer to [**[]GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerScopesInner**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerScopesInner.md) |  | [optional] 
**Profiles** | Pointer to [**[]GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerProfilesInner**](GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInnerProfilesInner.md) |  | [optional] 
**AppliedTargets** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkFirewallRuleGroups200ResponseAllOfRuleGroupsInnerRulesInner{
    // Set fields directly
}
```

### RuleGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.RuleGroup.IsSet()` — check if set
- `obj.RuleGroup.Get()` — get the inner value (returns pointer)
- `obj.RuleGroup.Set(&val)` — set the value
- `obj.RuleGroup.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


