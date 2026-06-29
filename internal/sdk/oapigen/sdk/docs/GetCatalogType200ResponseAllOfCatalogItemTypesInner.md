# GetCatalogType200ResponseAllOfCatalogItemTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Context** | Pointer to **string** |  | [optional] 
**Featured** | Pointer to **bool** |  | [optional] 
**AllowQuantity** | Pointer to **bool** | Can users order more than one of this item at a time. | [optional] 
**ImagePath** | Pointer to **string** |  | [optional] 
**DarkImagePath** | Pointer to **string** |  | [optional] 
**FormType** | Pointer to **string** |  | [optional] 
**Form** | Pointer to [**GetCatalogType200ResponseAllOfCatalogItemTypesInnerForm**](GetCatalogType200ResponseAllOfCatalogItemTypesInnerForm.md) |  | [optional] 
**OptionTypes** | Pointer to [**[]GetCatalogType200ResponseAllOfCatalogItemTypesInnerOptionTypesInner**](GetCatalogType200ResponseAllOfCatalogItemTypesInnerOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetCatalogType200ResponseAllOfCatalogItemTypesInner{
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


