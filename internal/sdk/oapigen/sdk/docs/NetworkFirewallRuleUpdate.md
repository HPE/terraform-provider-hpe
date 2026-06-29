# NetworkFirewallRuleUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RuleGroup** | Pointer to [**NetworkFirewallRuleUpdateRuleGroup**](NetworkFirewallRuleUpdateRuleGroup.md) |  | [optional] 
**Name** | Pointer to **string** | Network firewall rule name | [optional] 
**Description** | Pointer to **NullableString** | Network firewall rule description | [optional] 
**Enabled** | Pointer to **bool** | Use this to set enabled state | [optional] 
**Priority** | Pointer to **NullableInt64** | Network firewall rule priority | [optional] 
**Direction** | Pointer to **string** |  | [optional] 
**Sources** | Pointer to [**NetworkFirewallRuleUpdateSources**](NetworkFirewallRuleUpdateSources.md) |  | [optional] 
**Destinations** | Pointer to [**NetworkFirewallRuleUpdateDestinations**](NetworkFirewallRuleUpdateDestinations.md) |  | [optional] 
**Config** | Pointer to [**NetworkFirewallRuleUpdateConfig**](NetworkFirewallRuleUpdateConfig.md) |  | [optional] 
**Scopes** | Pointer to [**NetworkFirewallRuleUpdateScopes**](NetworkFirewallRuleUpdateScopes.md) |  | [optional] 
**Policy** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkFirewallRuleUpdate{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Priority (Nullable)

Use the Nullable wrapper methods:
- `obj.Priority.IsSet()` — check if set
- `obj.Priority.Get()` — get the inner value (returns pointer)
- `obj.Priority.Set(&val)` — set the value
- `obj.Priority.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


