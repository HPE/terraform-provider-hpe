# GetOsTypeImage200ResponseOsTypeImage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**VirtualImageId** | Pointer to **int64** | The id of the virtual image.  | [optional] 
**VirtualImageName** | Pointer to **string** | The name of the virtual image.   | [optional] 
**Account** | Pointer to **NullableInt64** | The account attached to the osTypeImage.   | [optional] 
**ProvisionType** | Pointer to **NullableInt64** | The Provision Type of the osTypeImage.  | [optional] 
**ComputeZoneType** | Pointer to **NullableInt64** | The zone type of the osTypeImage.  | [optional] 
**Zone** | Pointer to **NullableInt64** | The cloud that is attached to osTypeImage. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetOsTypeImage200ResponseOsTypeImage{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### ProvisionType (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionType.IsSet()` — check if set
- `obj.ProvisionType.Get()` — get the inner value (returns pointer)
- `obj.ProvisionType.Set(&val)` — set the value
- `obj.ProvisionType.Unset()` — clear the value
### ComputeZoneType (Nullable)

Use the Nullable wrapper methods:
- `obj.ComputeZoneType.IsSet()` — check if set
- `obj.ComputeZoneType.Get()` — get the inner value (returns pointer)
- `obj.ComputeZoneType.Set(&val)` — set the value
- `obj.ComputeZoneType.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


