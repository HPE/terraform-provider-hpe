# UpdateNetworkFirewallRuleGroupRequestRuleGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Network firewall rule group name | [optional] 
**Description** | Pointer to **NullableString** | Network firewall rule group description | [optional] 
**Priority** | Pointer to **NullableInt64** | Network firewall rule group priority | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateNetworkFirewallRuleGroupRequestRuleGroup{
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


