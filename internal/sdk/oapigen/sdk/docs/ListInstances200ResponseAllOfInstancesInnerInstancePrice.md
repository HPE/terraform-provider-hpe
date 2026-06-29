# ListInstances200ResponseAllOfInstancesInnerInstancePrice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Price** | Pointer to **float32** |  | [optional] 
**Cost** | Pointer to **float32** |  | [optional] 
**Currency** | Pointer to **NullableString** |  | [optional] 
**Unit** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstances200ResponseAllOfInstancesInnerInstancePrice{
    // Set fields directly
}
```

### Currency (Nullable)

Use the Nullable wrapper methods:
- `obj.Currency.IsSet()` — check if set
- `obj.Currency.Get()` — get the inner value (returns pointer)
- `obj.Currency.Set(&val)` — set the value
- `obj.Currency.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


