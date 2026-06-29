# UpdatePreseedScript200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PreseedScript** | Pointer to [**UpdatePreseedScript200ResponseAllOfPreseedScript**](UpdatePreseedScript200ResponseAllOfPreseedScript.md) |  | [optional] 
**ErrorCode** | Pointer to **NullableString** |  | [optional] 
**Success** | Pointer to **bool** | Success indicator, true when the request succeeded and false when an error occurred | [optional] [default to true]
**Msg** | Pointer to **NullableString** | Message containing a description of the result, usually a message about the error that occurred | [optional] 
**Errors** | Pointer to **map[string]interface{}** | Validation errors, with a key for Object containing error messages for each invalid parameter (key) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdatePreseedScript200Response{
    // Set fields directly
}
```

### ErrorCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorCode.IsSet()` — check if set
- `obj.ErrorCode.Get()` — get the inner value (returns pointer)
- `obj.ErrorCode.Set(&val)` — set the value
- `obj.ErrorCode.Unset()` — clear the value
### Msg (Nullable)

Use the Nullable wrapper methods:
- `obj.Msg.IsSet()` — check if set
- `obj.Msg.Get()` — get the inner value (returns pointer)
- `obj.Msg.Set(&val)` — set the value
- `obj.Msg.Unset()` — clear the value
### Errors (Nullable)

Use the Nullable wrapper methods:
- `obj.Errors.IsSet()` — check if set
- `obj.Errors.Get()` — get the inner value (returns pointer)
- `obj.Errors.Set(&val)` — set the value
- `obj.Errors.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


