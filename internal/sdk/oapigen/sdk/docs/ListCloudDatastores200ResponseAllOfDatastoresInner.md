# ListCloudDatastores200ResponseAllOfDatastoresInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Zone** | Pointer to [**NullableListCloudDatastores200ResponseAllOfDatastoresInnerZone**](ListCloudDatastores200ResponseAllOfDatastoresInnerZone.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**FreeSpace** | Pointer to **int64** |  | [optional] 
**Online** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Tenants** | Pointer to [**[]ListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInner**](ListCloudDatastores200ResponseAllOfDatastoresInnerTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**ListCloudDatastores200ResponseAllOfDatastoresInnerResourcePermission**](ListCloudDatastores200ResponseAllOfDatastoresInnerResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListCloudDatastores200ResponseAllOfDatastoresInner{
    // Set fields directly
}
```

### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
### Tenants (Nullable)

Use the Nullable wrapper methods:
- `obj.Tenants.IsSet()` — check if set
- `obj.Tenants.Get()` — get the inner value (returns pointer)
- `obj.Tenants.Set(&val)` — set the value
- `obj.Tenants.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


