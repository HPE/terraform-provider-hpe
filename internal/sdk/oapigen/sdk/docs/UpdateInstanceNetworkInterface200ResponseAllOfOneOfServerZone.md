# UpdateInstanceNetworkInterface200ResponseAllOfOneOfServerZone

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Groups** | Pointer to **[]int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Location** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**ZoneTypeId** | Pointer to **int64** |  | [optional] 
**NetworkServer** | Pointer to [**UpdateInstanceNetworkInterface200ResponseAllOfOneOfServerZoneNetworkServer**](UpdateInstanceNetworkInterface200ResponseAllOfOneOfServerZoneNetworkServer.md) |  | [optional] 
**SecurityServer** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInstanceNetworkInterface200ResponseAllOfOneOfServerZone{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### SecurityServer (Nullable)

Use the Nullable wrapper methods:
- `obj.SecurityServer.IsSet()` — check if set
- `obj.SecurityServer.Get()` — get the inner value (returns pointer)
- `obj.SecurityServer.Set(&val)` — set the value
- `obj.SecurityServer.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


