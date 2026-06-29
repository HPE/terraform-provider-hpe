# GetNetworkFirewallRule200ResponseRule

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
**RuleGroup** | Pointer to [**GetNetworkFirewallRule200ResponseRuleRuleGroup**](GetNetworkFirewallRule200ResponseRuleRuleGroup.md) |  | [optional] 
**GroupName** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Sources** | Pointer to [**[]GetNetworkFirewallRule200ResponseRuleSourcesInner**](GetNetworkFirewallRule200ResponseRuleSourcesInner.md) |  | [optional] 
**Destinations** | Pointer to [**[]GetNetworkFirewallRule200ResponseRuleDestinationsInner**](GetNetworkFirewallRule200ResponseRuleDestinationsInner.md) |  | [optional] 
**Applications** | Pointer to [**[]GetNetworkFirewallRule200ResponseRuleApplicationsInner**](GetNetworkFirewallRule200ResponseRuleApplicationsInner.md) |  | [optional] 
**Scopes** | Pointer to [**[]GetNetworkFirewallRule200ResponseRuleScopesInner**](GetNetworkFirewallRule200ResponseRuleScopesInner.md) |  | [optional] 
**Profiles** | Pointer to [**[]GetNetworkFirewallRule200ResponseRuleProfilesInner**](GetNetworkFirewallRule200ResponseRuleProfilesInner.md) |  | [optional] 
**AppliedTargets** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkFirewallRule200ResponseRule{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


