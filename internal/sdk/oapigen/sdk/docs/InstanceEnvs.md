# InstanceEnvs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Envs** | Pointer to [**[]InstanceEnvsEnvsInner**](InstanceEnvsEnvsInner.md) |  | [optional] 
**ReadOnlyEnvs** | Pointer to [**map[string]InstanceEnvsReadOnlyEnvsValue**](InstanceEnvsReadOnlyEnvsValue.md) |  | [optional] 
**ImportedEnvs** | Pointer to [**map[string]InstanceEnvsImportedEnvsValue**](InstanceEnvsImportedEnvsValue.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceEnvs{
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


