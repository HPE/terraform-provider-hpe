# HostUpdateConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserData** | Pointer to **NullableString** | User Data. Allows for override of cloud-init based user-data yaml or custom scripts | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &HostUpdateConfig{
    // Set fields directly
}
```

### UserData (Nullable)

Use the Nullable wrapper methods:
- `obj.UserData.IsSet()` — check if set
- `obj.UserData.Get()` — get the inner value (returns pointer)
- `obj.UserData.Set(&val)` — set the value
- `obj.UserData.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


