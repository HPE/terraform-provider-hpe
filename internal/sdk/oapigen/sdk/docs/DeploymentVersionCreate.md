# DeploymentVersionCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | Pointer to **string** | Version number (userVersion), a unique version identifier for the deployment version. | [optional] 
**UserVersion** | Pointer to **string** | Alias for version | [optional] 
**DeployType** | Pointer to **string** | Deploy Type, eg. file, git, fetch | [optional] 
**GitUrl** | Pointer to **NullableString** |  | [optional] 
**GitRef** | Pointer to **NullableString** |  | [optional] 
**FetchUrl** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &DeploymentVersionCreate{
    // Set fields directly
}
```

### GitUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.GitUrl.IsSet()` — check if set
- `obj.GitUrl.Get()` — get the inner value (returns pointer)
- `obj.GitUrl.Set(&val)` — set the value
- `obj.GitUrl.Unset()` — clear the value
### GitRef (Nullable)

Use the Nullable wrapper methods:
- `obj.GitRef.IsSet()` — check if set
- `obj.GitRef.Get()` — get the inner value (returns pointer)
- `obj.GitRef.Set(&val)` — set the value
- `obj.GitRef.Unset()` — clear the value
### FetchUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.FetchUrl.IsSet()` — check if set
- `obj.FetchUrl.Get()` — get the inner value (returns pointer)
- `obj.FetchUrl.Set(&val)` — set the value
- `obj.FetchUrl.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


