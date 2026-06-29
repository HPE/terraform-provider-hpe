# GetNetworkRouter200ResponseNetworkRouterFirewall

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**Version** | Pointer to **NullableString** |  | [optional] 
**DefaultPolicy** | Pointer to **NullableString** |  | [optional] 
**Global** | Pointer to **NullableString** |  | [optional] 
**RuleGroups** | Pointer to [**[]GetNetworkRouter200ResponseNetworkRouterFirewallRuleGroupsInner**](GetNetworkRouter200ResponseNetworkRouterFirewallRuleGroupsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouter200ResponseNetworkRouterFirewall{
    // Set fields directly
}
```

### Version (Nullable)

Use the Nullable wrapper methods:
- `obj.Version.IsSet()` — check if set
- `obj.Version.Get()` — get the inner value (returns pointer)
- `obj.Version.Set(&val)` — set the value
- `obj.Version.Unset()` — clear the value
### DefaultPolicy (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultPolicy.IsSet()` — check if set
- `obj.DefaultPolicy.Get()` — get the inner value (returns pointer)
- `obj.DefaultPolicy.Set(&val)` — set the value
- `obj.DefaultPolicy.Unset()` — clear the value
### Global (Nullable)

Use the Nullable wrapper methods:
- `obj.Global.IsSet()` — check if set
- `obj.Global.Get()` — get the inner value (returns pointer)
- `obj.Global.Set(&val)` — set the value
- `obj.Global.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


