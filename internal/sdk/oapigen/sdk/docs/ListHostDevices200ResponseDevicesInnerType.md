# ListHostDevices200ResponseDevicesInnerType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Family** | Pointer to **NullableString** |  | [optional] 
**BusType** | Pointer to **NullableString** |  | [optional] 
**Assignable** | Pointer to **bool** |  | [optional] 
**Hotpluggable** | Pointer to **bool** |  | [optional] 
**VendorId** | Pointer to **NullableInt32** |  | [optional] 
**ProductId** | Pointer to **NullableInt32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHostDevices200ResponseDevicesInnerType{
    // Set fields directly
}
```

### Family (Nullable)

Use the Nullable wrapper methods:
- `obj.Family.IsSet()` — check if set
- `obj.Family.Get()` — get the inner value (returns pointer)
- `obj.Family.Set(&val)` — set the value
- `obj.Family.Unset()` — clear the value
### BusType (Nullable)

Use the Nullable wrapper methods:
- `obj.BusType.IsSet()` — check if set
- `obj.BusType.Get()` — get the inner value (returns pointer)
- `obj.BusType.Set(&val)` — set the value
- `obj.BusType.Unset()` — clear the value
### VendorId (Nullable)

Use the Nullable wrapper methods:
- `obj.VendorId.IsSet()` — check if set
- `obj.VendorId.Get()` — get the inner value (returns pointer)
- `obj.VendorId.Set(&val)` — set the value
- `obj.VendorId.Unset()` — clear the value
### ProductId (Nullable)

Use the Nullable wrapper methods:
- `obj.ProductId.IsSet()` — check if set
- `obj.ProductId.Get()` — get the inner value (returns pointer)
- `obj.ProductId.Set(&val)` — set the value
- `obj.ProductId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


