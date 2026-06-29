# App

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Environment** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetApp200ResponseAppAccount**](GetApp200ResponseAppAccount.md) |  | [optional] 
**Owner** | Pointer to [**GetApp200ResponseAppOwner**](GetApp200ResponseAppOwner.md) |  | [optional] 
**SiteId** | Pointer to **int64** |  | [optional] 
**Group** | Pointer to [**GetApp200ResponseAppGroup**](GetApp200ResponseAppGroup.md) |  | [optional] 
**Blueprint** | Pointer to [**GetApp200ResponseAppBlueprint**](GetApp200ResponseAppBlueprint.md) |  | [optional] 
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
**Instances** | Pointer to [**[]GetApp200ResponseAppInstancesInner**](GetApp200ResponseAppInstancesInner.md) |  | [optional] 
**Stats** | Pointer to [**GetApp200ResponseAppStats**](GetApp200ResponseAppStats.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &App{
    // Set fields directly
}
```

### RemovalDate (Nullable)

Use the Nullable wrapper methods:
- `obj.RemovalDate.IsSet()` — check if set
- `obj.RemovalDate.Get()` — get the inner value (returns pointer)
- `obj.RemovalDate.Set(&val)` — set the value
- `obj.RemovalDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


