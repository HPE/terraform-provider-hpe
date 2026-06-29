# GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner

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
**RuleGroup** | Pointer to [**GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerRuleGroup**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerRuleGroup.md) |  | [optional] 
**GroupName** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Sources** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerSourcesInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerSourcesInner.md) |  | [optional] 
**Destinations** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerDestinationsInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerDestinationsInner.md) |  | [optional] 
**Applications** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerApplicationsInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerApplicationsInner.md) |  | [optional] 
**Scopes** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerScopesInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerScopesInner.md) |  | [optional] 
**Profiles** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerProfilesInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInnerProfilesInner.md) |  | [optional] 
**AppliedTargets** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


