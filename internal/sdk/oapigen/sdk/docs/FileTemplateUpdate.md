# FileTemplateUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | File template name | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**FileName** | Pointer to **string** | Filename for the file template | [optional] 
**FilePath** | Pointer to **string** | Path for the file template | [optional] 
**Category** | Pointer to **string** | Category | [optional] 
**TemplatePhase** | Pointer to **string** | Template Phase, provision, start, etc. | [optional] 
**Template** | Pointer to **string** | Template content, that is, the file template content itself. | [optional] 
**FileOwner** | Pointer to **int64** | File Owner | [optional] 
**SettingName** | Pointer to **string** | Setting Name | [optional] 
**SettingCategory** | Pointer to **string** | Setting Category | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &FileTemplateUpdate{
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


