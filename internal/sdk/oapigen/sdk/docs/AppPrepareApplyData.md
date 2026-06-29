# AppPrepareApplyData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Image** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**AutoValidate** | Pointer to **bool** |  | [optional] 
**Terraform** | Pointer to [**AppPrepareApplyDataTerraform**](AppPrepareApplyDataTerraform.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**BlueprintName** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**TemplateId** | Pointer to **int64** |  | [optional] 
**BlueprintId** | Pointer to **int64** |  | [optional] 
**Group** | Pointer to [**AppPrepareApplyDataGroup**](AppPrepareApplyDataGroup.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AppPrepareApplyData{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


