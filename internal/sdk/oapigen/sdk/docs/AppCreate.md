# AppCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TemplateId** | Pointer to **int64** |  | [optional] 
**BlueprintId** | [**AppCreateBlueprintId**](AppCreateBlueprintId.md) |  | 
**Name** | **string** | A unique name for the app | 
**Description** | Pointer to **string** | Description | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Group** | Pointer to [**AppCreateGroup**](AppCreateGroup.md) |  | [optional] 
**DefaultCloud** | Pointer to [**AppCreateDefaultCloud**](AppCreateDefaultCloud.md) |  | [optional] 
**Environment** | Pointer to **string** | Environment code (appContext) | [optional] 
**Tiers** | Pointer to **map[string]interface{}** | Configuration of app elements | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AppCreate{
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


