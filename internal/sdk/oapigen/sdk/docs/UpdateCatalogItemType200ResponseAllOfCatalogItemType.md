# UpdateCatalogItemType200ResponseAllOfCatalogItemType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** | Useful shortcode for provisioning naming schemes and export reference. | [optional] 
**Category** | Pointer to **NullableString** | Catalog Item Type category | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Featured** | Pointer to **bool** |  | [optional] 
**AllowQuantity** | Pointer to **bool** | Can users order more than one of this item at a time. | [optional] 
**IconPath** | Pointer to **string** |  | [optional] 
**ImagePath** | Pointer to **string** |  | [optional] 
**DarkImagePath** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**LayoutCode** | Pointer to **NullableString** |  | [optional] 
**Blueprint** | Pointer to **map[string]interface{}** |  | [optional] 
**AppSpec** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**InstanceSpec** | Pointer to **NullableString** |  | [optional] 
**Workflow** | Pointer to [**UpdateCatalogItemType200ResponseAllOfCatalogItemTypeWorkflow**](UpdateCatalogItemType200ResponseAllOfCatalogItemTypeWorkflow.md) |  | [optional] 
**Content** | Pointer to **NullableString** |  | [optional] 
**FormType** | Pointer to **string** |  | [optional] 
**Form** | Pointer to **map[string]interface{}** | Form object that contains input options and/or field groups | [optional] 
**FormConfig** | Pointer to **map[string]interface{}** | Form config object | [optional] 
**OptionTypes** | Pointer to [**[]UpdateCatalogItemType200ResponseAllOfCatalogItemTypeOptionTypesInner**](UpdateCatalogItemType200ResponseAllOfCatalogItemTypeOptionTypesInner.md) |  | [optional] 
**CreatedBy** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to [**UpdateCatalogItemType200ResponseAllOfCatalogItemTypeOwner**](UpdateCatalogItemType200ResponseAllOfCatalogItemTypeOwner.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateCatalogItemType200ResponseAllOfCatalogItemType{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LayoutCode (Nullable)

Use the Nullable wrapper methods:
- `obj.LayoutCode.IsSet()` — check if set
- `obj.LayoutCode.Get()` — get the inner value (returns pointer)
- `obj.LayoutCode.Set(&val)` — set the value
- `obj.LayoutCode.Unset()` — clear the value
### Blueprint (Nullable)

Use the Nullable wrapper methods:
- `obj.Blueprint.IsSet()` — check if set
- `obj.Blueprint.Get()` — get the inner value (returns pointer)
- `obj.Blueprint.Set(&val)` — set the value
- `obj.Blueprint.Unset()` — clear the value
### AppSpec (Nullable)

Use the Nullable wrapper methods:
- `obj.AppSpec.IsSet()` — check if set
- `obj.AppSpec.Get()` — get the inner value (returns pointer)
- `obj.AppSpec.Set(&val)` — set the value
- `obj.AppSpec.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### InstanceSpec (Nullable)

Use the Nullable wrapper methods:
- `obj.InstanceSpec.IsSet()` — check if set
- `obj.InstanceSpec.Get()` — get the inner value (returns pointer)
- `obj.InstanceSpec.Set(&val)` — set the value
- `obj.InstanceSpec.Unset()` — clear the value
### Content (Nullable)

Use the Nullable wrapper methods:
- `obj.Content.IsSet()` — check if set
- `obj.Content.Get()` — get the inner value (returns pointer)
- `obj.Content.Set(&val)` — set the value
- `obj.Content.Unset()` — clear the value
### Form (Nullable)

Use the Nullable wrapper methods:
- `obj.Form.IsSet()` — check if set
- `obj.Form.Get()` — get the inner value (returns pointer)
- `obj.Form.Set(&val)` — set the value
- `obj.Form.Unset()` — clear the value
### FormConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.FormConfig.IsSet()` — check if set
- `obj.FormConfig.Get()` — get the inner value (returns pointer)
- `obj.FormConfig.Set(&val)` — set the value
- `obj.FormConfig.Unset()` — clear the value
### OptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionTypes.IsSet()` — check if set
- `obj.OptionTypes.Get()` — get the inner value (returns pointer)
- `obj.OptionTypes.Set(&val)` — set the value
- `obj.OptionTypes.Unset()` — clear the value
### CreatedBy (Nullable)

Use the Nullable wrapper methods:
- `obj.CreatedBy.IsSet()` — check if set
- `obj.CreatedBy.Get()` — get the inner value (returns pointer)
- `obj.CreatedBy.Set(&val)` — set the value
- `obj.CreatedBy.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


