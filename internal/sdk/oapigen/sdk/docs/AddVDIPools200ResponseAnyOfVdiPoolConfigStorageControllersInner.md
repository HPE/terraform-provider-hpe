# AddVDIPools200ResponseAnyOfVdiPoolConfigStorageControllersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **NullableBool** |  | [optional] 
**TypeId** | Pointer to **int64** |  | [optional] 
**TypeName** | Pointer to **string** |  | [optional] 
**UnitNumber** | Pointer to **string** |  | [optional] 
**BusNumber** | Pointer to **string** |  | [optional] 
**MaxDevices** | Pointer to **float32** |  | [optional] 
**Removable** | Pointer to **NullableBool** |  | [optional] 
**Editable** | Pointer to **bool** |  | [optional] 
**ReservedUnitNumber** | Pointer to **float32** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**DisplayOrder** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddVDIPools200ResponseAnyOfVdiPoolConfigStorageControllersInner{
    // Set fields directly
}
```

### Active (Nullable)

Use the Nullable wrapper methods:
- `obj.Active.IsSet()` — check if set
- `obj.Active.Get()` — get the inner value (returns pointer)
- `obj.Active.Set(&val)` — set the value
- `obj.Active.Unset()` — clear the value
### Removable (Nullable)

Use the Nullable wrapper methods:
- `obj.Removable.IsSet()` — check if set
- `obj.Removable.Get()` — get the inner value (returns pointer)
- `obj.Removable.Set(&val)` — set the value
- `obj.Removable.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


