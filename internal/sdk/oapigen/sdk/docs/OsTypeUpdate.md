# OsTypeUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the osType.  | [optional] 
**Description** | Pointer to **NullableString** | The description of the osType.   | [optional] 
**Platform** | Pointer to **string** | The platform of the osType.   | [optional] 
**Category** | Pointer to **NullableString** | The category of the osType.  | [optional] 
**Vendor** | Pointer to **NullableString** | The vendor of the osType.  | [optional] 
**OsName** | Pointer to **NullableString** | The osName of the osType.  | [optional] 
**OsVersion** | Pointer to **NullableString** | The osVersion of the osType.  | [optional] 
**OsCodename** | Pointer to **NullableString** | The osCodename of the osType.  | [optional] 
**OsFamily** | Pointer to **NullableString** | The family of the osType.  | [optional] 
**BitCount** | Pointer to **int64** | The bitCount/architecture of the osType.  | [optional] 
**CloudInitVersion** | Pointer to **string** | The version of CloudInit being used.  | [optional] 
**InstallAgent** | Pointer to **NullableBool** | Whether the morpheus agent is installed.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &OsTypeUpdate{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Vendor (Nullable)

Use the Nullable wrapper methods:
- `obj.Vendor.IsSet()` — check if set
- `obj.Vendor.Get()` — get the inner value (returns pointer)
- `obj.Vendor.Set(&val)` — set the value
- `obj.Vendor.Unset()` — clear the value
### OsName (Nullable)

Use the Nullable wrapper methods:
- `obj.OsName.IsSet()` — check if set
- `obj.OsName.Get()` — get the inner value (returns pointer)
- `obj.OsName.Set(&val)` — set the value
- `obj.OsName.Unset()` — clear the value
### OsVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.OsVersion.IsSet()` — check if set
- `obj.OsVersion.Get()` — get the inner value (returns pointer)
- `obj.OsVersion.Set(&val)` — set the value
- `obj.OsVersion.Unset()` — clear the value
### OsCodename (Nullable)

Use the Nullable wrapper methods:
- `obj.OsCodename.IsSet()` — check if set
- `obj.OsCodename.Get()` — get the inner value (returns pointer)
- `obj.OsCodename.Set(&val)` — set the value
- `obj.OsCodename.Unset()` — clear the value
### OsFamily (Nullable)

Use the Nullable wrapper methods:
- `obj.OsFamily.IsSet()` — check if set
- `obj.OsFamily.Get()` — get the inner value (returns pointer)
- `obj.OsFamily.Set(&val)` — set the value
- `obj.OsFamily.Unset()` — clear the value
### InstallAgent (Nullable)

Use the Nullable wrapper methods:
- `obj.InstallAgent.IsSet()` — check if set
- `obj.InstallAgent.Get()` — get the inner value (returns pointer)
- `obj.InstallAgent.Set(&val)` — set the value
- `obj.InstallAgent.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


