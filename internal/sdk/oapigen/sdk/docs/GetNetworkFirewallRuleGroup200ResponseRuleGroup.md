# GetNetworkFirewallRuleGroup200ResponseRuleGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Priority** | Pointer to **int64** |  | [optional] 
**GroupLayer** | Pointer to **string** |  | [optional] 
**Rules** | Pointer to [**[]GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner**](GetNetworkFirewallRuleGroup200ResponseRuleGroupRulesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkFirewallRuleGroup200ResponseRuleGroup{
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


