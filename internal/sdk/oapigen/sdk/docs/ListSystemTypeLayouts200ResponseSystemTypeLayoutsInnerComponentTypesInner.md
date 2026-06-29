# ListSystemTypeLayouts200ResponseSystemTypeLayoutsInnerComponentTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Unique identifier for the component type. | [optional] 
**Code** | Pointer to **string** | Unique code identifier for the component type. | [optional] 
**Name** | Pointer to **string** | Display name of the component type. | [optional] 
**Category** | Pointer to **NullableString** | Category grouping for the component type. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListSystemTypeLayouts200ResponseSystemTypeLayoutsInnerComponentTypesInner{
    // Set fields directly
}
```

### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


