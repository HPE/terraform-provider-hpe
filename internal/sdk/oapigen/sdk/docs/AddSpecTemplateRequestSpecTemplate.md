# AddSpecTemplateRequestSpecTemplate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Spec template name | 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**Type** | [**AddSpecTemplateRequestSpecTemplateType**](AddSpecTemplateRequestSpecTemplateType.md) |  | 
**File** | [**AddSpecTemplateRequestSpecTemplateFile**](AddSpecTemplateRequestSpecTemplateFile.md) |  | 
**Config** | Pointer to [**AddSpecTemplateRequestSpecTemplateConfig**](AddSpecTemplateRequestSpecTemplateConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddSpecTemplateRequestSpecTemplate{
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


