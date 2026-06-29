# AppUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A unique name for the app | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Environment** | Pointer to **string** | Environment code (appContext) | [optional] 
**OwnerId** | Pointer to **int64** | User ID, can be used to change app owner. This also changes the owner for each instance in the app. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AppUpdate{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


