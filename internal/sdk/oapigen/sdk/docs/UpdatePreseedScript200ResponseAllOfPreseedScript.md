# UpdatePreseedScript200ResponseAllOfPreseedScript

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**UpdatePreseedScript200ResponseAllOfPreseedScriptAccount**](UpdatePreseedScript200ResponseAllOfPreseedScriptAccount.md) |  | [optional] 
**FileName** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Content** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to [**UpdatePreseedScript200ResponseAllOfPreseedScriptCreatedBy**](UpdatePreseedScript200ResponseAllOfPreseedScriptCreatedBy.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdatePreseedScript200ResponseAllOfPreseedScript{
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


