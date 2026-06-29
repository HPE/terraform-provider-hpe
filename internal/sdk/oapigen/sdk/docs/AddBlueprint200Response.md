# AddBlueprint200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Blueprint** | Pointer to [**AddBlueprint200ResponseAllOfBlueprint**](AddBlueprint200ResponseAllOfBlueprint.md) |  | [optional] 
**Msg** | Pointer to **NullableString** |  | [optional] 
**Errors** | Pointer to **NullableString** |  | [optional] 
**ErrorCode** | Pointer to **NullableString** |  | [optional] 
**InProgress** | Pointer to **bool** |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprint200Response{
    // Set fields directly
}
```

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
### ErrorCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ErrorCode.IsSet()` — check if set
- `obj.ErrorCode.Get()` — get the inner value (returns pointer)
- `obj.ErrorCode.Set(&val)` — set the value
- `obj.ErrorCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


