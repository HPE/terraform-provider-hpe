# GetFileTemplate200ResponseContainerTemplate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Account** | Pointer to [**GetFileTemplate200ResponseContainerTemplateAccount**](GetFileTemplate200ResponseContainerTemplateAccount.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**FileName** | Pointer to **string** |  | [optional] 
**FilePath** | Pointer to **string** |  | [optional] 
**TemplateType** | Pointer to **NullableString** |  | [optional] 
**TemplatePhase** | Pointer to **string** |  | [optional] 
**Template** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**SettingCategory** | Pointer to **NullableString** |  | [optional] 
**SettingName** | Pointer to **NullableString** |  | [optional] 
**AutoRun** | Pointer to **bool** |  | [optional] 
**RunOnScale** | Pointer to **NullableBool** |  | [optional] 
**RunOnDeploy** | Pointer to **NullableBool** |  | [optional] 
**FileOwner** | Pointer to **NullableString** |  | [optional] 
**FileGroup** | Pointer to **NullableString** |  | [optional] 
**Permissions** | Pointer to **NullableString** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetFileTemplate200ResponseContainerTemplate{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### TemplateType (Nullable)

Use the Nullable wrapper methods:
- `obj.TemplateType.IsSet()` — check if set
- `obj.TemplateType.Get()` — get the inner value (returns pointer)
- `obj.TemplateType.Set(&val)` — set the value
- `obj.TemplateType.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### SettingCategory (Nullable)

Use the Nullable wrapper methods:
- `obj.SettingCategory.IsSet()` — check if set
- `obj.SettingCategory.Get()` — get the inner value (returns pointer)
- `obj.SettingCategory.Set(&val)` — set the value
- `obj.SettingCategory.Unset()` — clear the value
### SettingName (Nullable)

Use the Nullable wrapper methods:
- `obj.SettingName.IsSet()` — check if set
- `obj.SettingName.Get()` — get the inner value (returns pointer)
- `obj.SettingName.Set(&val)` — set the value
- `obj.SettingName.Unset()` — clear the value
### RunOnScale (Nullable)

Use the Nullable wrapper methods:
- `obj.RunOnScale.IsSet()` — check if set
- `obj.RunOnScale.Get()` — get the inner value (returns pointer)
- `obj.RunOnScale.Set(&val)` — set the value
- `obj.RunOnScale.Unset()` — clear the value
### RunOnDeploy (Nullable)

Use the Nullable wrapper methods:
- `obj.RunOnDeploy.IsSet()` — check if set
- `obj.RunOnDeploy.Get()` — get the inner value (returns pointer)
- `obj.RunOnDeploy.Set(&val)` — set the value
- `obj.RunOnDeploy.Unset()` — clear the value
### FileOwner (Nullable)

Use the Nullable wrapper methods:
- `obj.FileOwner.IsSet()` — check if set
- `obj.FileOwner.Get()` — get the inner value (returns pointer)
- `obj.FileOwner.Set(&val)` — set the value
- `obj.FileOwner.Unset()` — clear the value
### FileGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.FileGroup.IsSet()` — check if set
- `obj.FileGroup.Get()` — get the inner value (returns pointer)
- `obj.FileGroup.Set(&val)` — set the value
- `obj.FileGroup.Unset()` — clear the value
### Permissions (Nullable)

Use the Nullable wrapper methods:
- `obj.Permissions.IsSet()` — check if set
- `obj.Permissions.Get()` — get the inner value (returns pointer)
- `obj.Permissions.Set(&val)` — set the value
- `obj.Permissions.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


