# CatalogTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Context** | Pointer to **string** |  | [optional] 
**Featured** | Pointer to **bool** |  | [optional] 
**AllowQuantity** | Pointer to **bool** |  | [optional] 
**ImagePath** | Pointer to **string** |  | [optional] 
**DarkImagePath** | Pointer to **string** |  | [optional] 
**FormType** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CatalogTypesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


