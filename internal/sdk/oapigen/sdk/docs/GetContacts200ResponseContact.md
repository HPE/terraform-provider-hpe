# GetContacts200ResponseContact

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**EmailAddress** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**SmsAddress** | Pointer to **string** |  | [optional] 
**SlackHook** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetContacts200ResponseContact{
    // Set fields directly
}
```

### SlackHook (Nullable)

Use the Nullable wrapper methods:
- `obj.SlackHook.IsSet()` — check if set
- `obj.SlackHook.Get()` — get the inner value (returns pointer)
- `obj.SlackHook.Set(&val)` — set the value
- `obj.SlackHook.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


