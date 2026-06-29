# GetNetworkFirewallRules200ResponseAllOfRulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Direction** | Pointer to **string** |  | [optional] 
**SourceType** | Pointer to **string** |  | [optional] 
**DestinationType** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Policy** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **int64** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**RuleGroup** | Pointer to [**NullableGetNetworkFirewallRules200ResponseAllOfRulesInnerRuleGroup**](GetNetworkFirewallRules200ResponseAllOfRulesInnerRuleGroup.md) |  | [optional] 
**GroupName** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Sources** | Pointer to [**[]GetNetworkFirewallRules200ResponseAllOfRulesInnerSourcesInner**](GetNetworkFirewallRules200ResponseAllOfRulesInnerSourcesInner.md) |  | [optional] 
**Destinations** | Pointer to [**[]GetNetworkFirewallRules200ResponseAllOfRulesInnerDestinationsInner**](GetNetworkFirewallRules200ResponseAllOfRulesInnerDestinationsInner.md) |  | [optional] 
**Applications** | Pointer to [**[]GetNetworkFirewallRules200ResponseAllOfRulesInnerApplicationsInner**](GetNetworkFirewallRules200ResponseAllOfRulesInnerApplicationsInner.md) |  | [optional] 
**Scopes** | Pointer to [**[]GetNetworkFirewallRules200ResponseAllOfRulesInnerScopesInner**](GetNetworkFirewallRules200ResponseAllOfRulesInnerScopesInner.md) |  | [optional] 
**Profiles** | Pointer to [**[]GetNetworkFirewallRules200ResponseAllOfRulesInnerProfilesInner**](GetNetworkFirewallRules200ResponseAllOfRulesInnerProfilesInner.md) |  | [optional] 
**AppliedTargets** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkFirewallRules200ResponseAllOfRulesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### RuleGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.RuleGroup.IsSet()` — check if set
- `obj.RuleGroup.Get()` — get the inner value (returns pointer)
- `obj.RuleGroup.Set(&val)` — set the value
- `obj.RuleGroup.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


