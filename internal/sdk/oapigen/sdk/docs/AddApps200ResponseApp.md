# AddApps200ResponseApp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Environment** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableAddApps200ResponseAppAccount**](AddApps200ResponseAppAccount.md) |  | [optional] 
**Owner** | Pointer to [**AddApps200ResponseAppOwner**](AddApps200ResponseAppOwner.md) |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**Group** | Pointer to [**NullableAddApps200ResponseAppGroup**](AddApps200ResponseAppGroup.md) |  | [optional] 
**Blueprint** | Pointer to [**AddApps200ResponseAppBlueprint**](AddApps200ResponseAppBlueprint.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**RemovalDate** | Pointer to **NullableTime** |  | [optional] 
**AppContext** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**AppStatus** | Pointer to **string** |  | [optional] 
**InstanceCount** | Pointer to **int64** |  | [optional] 
**ContainerCount** | Pointer to **int64** |  | [optional] 
**AppTiers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Instances** | Pointer to [**[]AddApps200ResponseAppInstancesInner**](AddApps200ResponseAppInstancesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddApps200ResponseApp{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Group (Nullable)

Use the Nullable wrapper methods:
- `obj.Group.IsSet()` — check if set
- `obj.Group.Get()` — get the inner value (returns pointer)
- `obj.Group.Set(&val)` — set the value
- `obj.Group.Unset()` — clear the value
### RemovalDate (Nullable)

Use the Nullable wrapper methods:
- `obj.RemovalDate.IsSet()` — check if set
- `obj.RemovalDate.Get()` — get the inner value (returns pointer)
- `obj.RemovalDate.Set(&val)` — set the value
- `obj.RemovalDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


