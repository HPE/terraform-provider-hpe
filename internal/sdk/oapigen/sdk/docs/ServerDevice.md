# ServerDevice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**RefType** | Pointer to **NullableString** |  | [optional] 
**RefId** | Pointer to **NullableInt64** | (Assignee) Target Server ID | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**DomainId** | Pointer to **NullableInt32** |  | [optional] 
**Bus** | Pointer to **NullableInt32** |  | [optional] 
**Slot** | Pointer to **NullableInt32** |  | [optional] 
**Device** | Pointer to **NullableInt32** |  | [optional] 
**VendorId** | Pointer to **NullableString** |  | [optional] 
**ProductId** | Pointer to **NullableString** |  | [optional] 
**FunctionId** | Pointer to **NullableInt32** |  | [optional] 
**UniqueId** | Pointer to **NullableString** |  | [optional] 
**IommuGroup** | Pointer to **NullableInt32** |  | [optional] 
**IommuDeviceCount** | Pointer to **NullableInt32** |  | [optional] 
**Type** | Pointer to [**ServerDeviceType**](ServerDeviceType.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ServerDevice{
    // Set fields directly
}
```

### RefType (Nullable)

Use the Nullable wrapper methods:
- `obj.RefType.IsSet()` — check if set
- `obj.RefType.Get()` — get the inner value (returns pointer)
- `obj.RefType.Set(&val)` — set the value
- `obj.RefType.Unset()` — clear the value
### RefId (Nullable)

Use the Nullable wrapper methods:
- `obj.RefId.IsSet()` — check if set
- `obj.RefId.Get()` — get the inner value (returns pointer)
- `obj.RefId.Set(&val)` — set the value
- `obj.RefId.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value
### DomainId (Nullable)

Use the Nullable wrapper methods:
- `obj.DomainId.IsSet()` — check if set
- `obj.DomainId.Get()` — get the inner value (returns pointer)
- `obj.DomainId.Set(&val)` — set the value
- `obj.DomainId.Unset()` — clear the value
### Bus (Nullable)

Use the Nullable wrapper methods:
- `obj.Bus.IsSet()` — check if set
- `obj.Bus.Get()` — get the inner value (returns pointer)
- `obj.Bus.Set(&val)` — set the value
- `obj.Bus.Unset()` — clear the value
### Slot (Nullable)

Use the Nullable wrapper methods:
- `obj.Slot.IsSet()` — check if set
- `obj.Slot.Get()` — get the inner value (returns pointer)
- `obj.Slot.Set(&val)` — set the value
- `obj.Slot.Unset()` — clear the value
### Device (Nullable)

Use the Nullable wrapper methods:
- `obj.Device.IsSet()` — check if set
- `obj.Device.Get()` — get the inner value (returns pointer)
- `obj.Device.Set(&val)` — set the value
- `obj.Device.Unset()` — clear the value
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
### FunctionId (Nullable)

Use the Nullable wrapper methods:
- `obj.FunctionId.IsSet()` — check if set
- `obj.FunctionId.Get()` — get the inner value (returns pointer)
- `obj.FunctionId.Set(&val)` — set the value
- `obj.FunctionId.Unset()` — clear the value
### UniqueId (Nullable)

Use the Nullable wrapper methods:
- `obj.UniqueId.IsSet()` — check if set
- `obj.UniqueId.Get()` — get the inner value (returns pointer)
- `obj.UniqueId.Set(&val)` — set the value
- `obj.UniqueId.Unset()` — clear the value
### IommuGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.IommuGroup.IsSet()` — check if set
- `obj.IommuGroup.Get()` — get the inner value (returns pointer)
- `obj.IommuGroup.Set(&val)` — set the value
- `obj.IommuGroup.Unset()` — clear the value
### IommuDeviceCount (Nullable)

Use the Nullable wrapper methods:
- `obj.IommuDeviceCount.IsSet()` — check if set
- `obj.IommuDeviceCount.Get()` — get the inner value (returns pointer)
- `obj.IommuDeviceCount.Set(&val)` — set the value
- `obj.IommuDeviceCount.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


