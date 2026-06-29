# GetStorageVolumeTypes200ResponseStorageVolumeType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**DefaultType** | Pointer to **bool** |  | [optional] 
**CustomLabel** | Pointer to **bool** |  | [optional] 
**CustomSize** | Pointer to **bool** |  | [optional] 
**CustomSizeOptions** | Pointer to **NullableString** |  | [optional] 
**ConfigurableIOPS** | Pointer to **bool** |  | [optional] 
**HasDatastore** | Pointer to **bool** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetStorageVolumeTypes200ResponseStorageVolumeTypeOptionTypesInner**](GetStorageVolumeTypes200ResponseStorageVolumeTypeOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetStorageVolumeTypes200ResponseStorageVolumeType{
    // Set fields directly
}
```

### CustomSizeOptions (Nullable)

Use the Nullable wrapper methods:
- `obj.CustomSizeOptions.IsSet()` — check if set
- `obj.CustomSizeOptions.Get()` — get the inner value (returns pointer)
- `obj.CustomSizeOptions.Set(&val)` — set the value
- `obj.CustomSizeOptions.Unset()` — clear the value
### OptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionTypes.IsSet()` — check if set
- `obj.OptionTypes.Get()` — get the inner value (returns pointer)
- `obj.OptionTypes.Set(&val)` — set the value
- `obj.OptionTypes.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


