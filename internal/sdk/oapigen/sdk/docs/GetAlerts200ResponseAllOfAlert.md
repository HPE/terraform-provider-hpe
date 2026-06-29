# GetAlerts200ResponseAllOfAlert

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**AllApps** | Pointer to **bool** |  | [optional] 
**AllChecks** | Pointer to **bool** |  | [optional] 
**AllGroups** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**MinSeverity** | Pointer to **string** |  | [optional] 
**MinDuration** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Checks** | Pointer to **[]int32** |  | [optional] 
**CheckGroups** | Pointer to **[]int32** |  | [optional] 
**Apps** | Pointer to **[]int32** |  | [optional] 
**Contacts** | Pointer to [**[]GetAlerts200ResponseAllOfAlertContactsInner**](GetAlerts200ResponseAllOfAlertContactsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetAlerts200ResponseAllOfAlert{
    // Set fields directly
}
```

### Checks (Nullable)

Use the Nullable wrapper methods:
- `obj.Checks.IsSet()` — check if set
- `obj.Checks.Get()` — get the inner value (returns pointer)
- `obj.Checks.Set(&val)` — set the value
- `obj.Checks.Unset()` — clear the value
### CheckGroups (Nullable)

Use the Nullable wrapper methods:
- `obj.CheckGroups.IsSet()` — check if set
- `obj.CheckGroups.Get()` — get the inner value (returns pointer)
- `obj.CheckGroups.Set(&val)` — set the value
- `obj.CheckGroups.Unset()` — clear the value
### Apps (Nullable)

Use the Nullable wrapper methods:
- `obj.Apps.IsSet()` — check if set
- `obj.Apps.Get()` — get the inner value (returns pointer)
- `obj.Apps.Set(&val)` — set the value
- `obj.Apps.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


