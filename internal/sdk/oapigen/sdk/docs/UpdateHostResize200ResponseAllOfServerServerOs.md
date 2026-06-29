# UpdateHostResize200ResponseAllOfServerServerOs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Vendor** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**OsFamily** | Pointer to **NullableString** |  | [optional] 
**OsVersion** | Pointer to **string** |  | [optional] 
**BitCount** | Pointer to **int64** |  | [optional] 
**Platform** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostResize200ResponseAllOfServerServerOs{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### OsFamily (Nullable)

Use the Nullable wrapper methods:
- `obj.OsFamily.IsSet()` — check if set
- `obj.OsFamily.Get()` — get the inner value (returns pointer)
- `obj.OsFamily.Set(&val)` — set the value
- `obj.OsFamily.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


