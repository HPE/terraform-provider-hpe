# GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**IacId** | Pointer to **NullableString** |  | [optional] 
**Zone** | Pointer to **NullableString** |  | [optional] 
**ZonePool** | Pointer to **NullableString** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **int64** |  | [optional] 
**GroupLayer** | Pointer to **string** |  | [optional] 
**Rules** | Pointer to [**[]GetNetworkRouterFirewallRuleGroup200ResponseRuleGroupRulesInner**](GetNetworkRouterFirewallRuleGroup200ResponseRuleGroupRulesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouterFirewallRuleGroup200ResponseRuleGroup{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### IacId (Nullable)

Use the Nullable wrapper methods:
- `obj.IacId.IsSet()` — check if set
- `obj.IacId.Get()` — get the inner value (returns pointer)
- `obj.IacId.Set(&val)` — set the value
- `obj.IacId.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
### ZonePool (Nullable)

Use the Nullable wrapper methods:
- `obj.ZonePool.IsSet()` — check if set
- `obj.ZonePool.Get()` — get the inner value (returns pointer)
- `obj.ZonePool.Set(&val)` — set the value
- `obj.ZonePool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


