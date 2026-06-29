# GetEnvVariables200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Envs** | Pointer to [**[]GetEnvVariables200ResponseEnvsInner**](GetEnvVariables200ResponseEnvsInner.md) |  | [optional] 
**ReadOnlyEnvs** | Pointer to [**map[string]GetEnvVariables200ResponseReadOnlyEnvsValue**](GetEnvVariables200ResponseReadOnlyEnvsValue.md) |  | [optional] 
**ImportedEnvs** | Pointer to [**map[string]GetEnvVariables200ResponseImportedEnvsValue**](GetEnvVariables200ResponseImportedEnvsValue.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetEnvVariables200Response{
    // Set fields directly
}
```

### Envs (Nullable)

Use the Nullable wrapper methods:
- `obj.Envs.IsSet()` — check if set
- `obj.Envs.Get()` — get the inner value (returns pointer)
- `obj.Envs.Set(&val)` — set the value
- `obj.Envs.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


