# ListBillingServers200ResponseAllOfBillingInfoServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefType** | Pointer to **string** |  | [optional] 
**RefUUID** | Pointer to **string** |  | [optional] 
**RefId** | Pointer to **int64** |  | [optional] 
**StartDate** | Pointer to **time.Time** |  | [optional] 
**EndDate** | Pointer to **time.Time** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Price** | Pointer to **float32** |  | [optional] 
**NumUnits** | Pointer to **float32** |  | [optional] 
**Unit** | Pointer to **string** |  | [optional] 
**Currency** | Pointer to **string** |  | [optional] 
**Usages** | Pointer to [**[]ListBillingServers200ResponseAllOfBillingInfoServersInnerUsagesInner**](ListBillingServers200ResponseAllOfBillingInfoServersInnerUsagesInner.md) |  | [optional] 
**NumUsages** | Pointer to **int64** |  | [optional] 
**TotalUsages** | Pointer to **int64** |  | [optional] 
**HasMoreUsages** | Pointer to **bool** |  | [optional] 
**FoundPricing** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ServerUUID** | Pointer to **string** |  | [optional] 
**ServerUniqueId** | Pointer to **NullableString** |  | [optional] 
**ServerExternalId** | Pointer to **NullableString** |  | [optional] 
**ServerInternalId** | Pointer to **NullableString** |  | [optional] 
**ResourcePoolId** | Pointer to **NullableString** |  | [optional] 
**ResourcePoolName** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingServers200ResponseAllOfBillingInfoServersInner{
    // Set fields directly
}
```

### ServerUniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerUniqueId.IsSet()` — check if set
- `obj.ServerUniqueId.Get()` — get the inner value (returns pointer)
- `obj.ServerUniqueId.Set(&val)` — set the value
- `obj.ServerUniqueId.Unset()` — clear the value
### ServerExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerExternalId.IsSet()` — check if set
- `obj.ServerExternalId.Get()` — get the inner value (returns pointer)
- `obj.ServerExternalId.Set(&val)` — set the value
- `obj.ServerExternalId.Unset()` — clear the value
### ServerInternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ServerInternalId.IsSet()` — check if set
- `obj.ServerInternalId.Get()` — get the inner value (returns pointer)
- `obj.ServerInternalId.Set(&val)` — set the value
- `obj.ServerInternalId.Unset()` — clear the value
### ResourcePoolId (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePoolId.IsSet()` — check if set
- `obj.ResourcePoolId.Get()` — get the inner value (returns pointer)
- `obj.ResourcePoolId.Set(&val)` — set the value
- `obj.ResourcePoolId.Unset()` — clear the value
### ResourcePoolName (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePoolName.IsSet()` — check if set
- `obj.ResourcePoolName.Get()` — get the inner value (returns pointer)
- `obj.ResourcePoolName.Set(&val)` — set the value
- `obj.ResourcePoolName.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


