# ListPreseedScripts200ResponseAllOfPreseedScriptsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableListPreseedScripts200ResponseAllOfPreseedScriptsInnerAccount**](ListPreseedScripts200ResponseAllOfPreseedScriptsInnerAccount.md) |  | [optional] 
**FileName** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Content** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to [**ListPreseedScripts200ResponseAllOfPreseedScriptsInnerCreatedBy**](ListPreseedScripts200ResponseAllOfPreseedScriptsInnerCreatedBy.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListPreseedScripts200ResponseAllOfPreseedScriptsInner{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


