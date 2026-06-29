# GetNetworkRouterType200ResponseNetworkRouterTypeRuleGroupOptionTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**FieldName** | Pointer to **string** |  | [optional] 
**FieldLabel** | Pointer to **string** |  | [optional] 
**FieldCode** | Pointer to **NullableString** |  | [optional] 
**FieldContext** | Pointer to **string** |  | [optional] 
**FieldGroup** | Pointer to **NullableString** |  | [optional] 
**FieldClass** | Pointer to **NullableString** |  | [optional] 
**FieldAddOn** | Pointer to **NullableString** |  | [optional] 
**FieldComponent** | Pointer to **NullableString** |  | [optional] 
**FieldInput** | Pointer to **NullableString** |  | [optional] 
**PlaceHolder** | Pointer to **NullableString** |  | [optional] 
**VerifyPattern** | Pointer to **NullableString** |  | [optional] 
**HelpBlock** | Pointer to **NullableString** |  | [optional] 
**HelpBlockFieldCode** | Pointer to **NullableString** |  | [optional] 
**DefaultValue** | Pointer to **NullableString** |  | [optional] 
**OptionSource** | Pointer to **NullableString** |  | [optional] 
**OptionSourceType** | Pointer to **NullableString** |  | [optional] 
**OptionList** | Pointer to [**ListClusterTypes200ResponseAllOfClusterTypesInnerOptionTypesInnerOptionList**](ListClusterTypes200ResponseAllOfClusterTypesInnerOptionTypesInnerOptionList.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Advanced** | Pointer to **bool** |  | [optional] 
**Required** | Pointer to **bool** |  | [optional] 
**ExportMeta** | Pointer to **bool** |  | [optional] 
**Editable** | Pointer to **bool** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**DisplayOrder** | Pointer to **int64** |  | [optional] 
**WrapperClass** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**NoBlank** | Pointer to **bool** |  | [optional] 
**DependsOnCode** | Pointer to **NullableString** |  | [optional] 
**VisibleOnCode** | Pointer to **NullableString** |  | [optional] 
**RequireOnCode** | Pointer to **NullableString** |  | [optional] 
**ContextualDefault** | Pointer to **NullableBool** |  | [optional] 
**DisplayValueOnDetails** | Pointer to **NullableBool** |  | [optional] 
**ShowOnCreate** | Pointer to **NullableBool** |  | [optional] 
**ShowOnEdit** | Pointer to **NullableBool** |  | [optional] 
**LocalCredential** | Pointer to **NullableBool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouterType200ResponseNetworkRouterTypeRuleGroupOptionTypesInner{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### FieldCode (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldCode.IsSet()` — check if set
- `obj.FieldCode.Get()` — get the inner value (returns pointer)
- `obj.FieldCode.Set(&val)` — set the value
- `obj.FieldCode.Unset()` — clear the value
### FieldGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldGroup.IsSet()` — check if set
- `obj.FieldGroup.Get()` — get the inner value (returns pointer)
- `obj.FieldGroup.Set(&val)` — set the value
- `obj.FieldGroup.Unset()` — clear the value
### FieldClass (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldClass.IsSet()` — check if set
- `obj.FieldClass.Get()` — get the inner value (returns pointer)
- `obj.FieldClass.Set(&val)` — set the value
- `obj.FieldClass.Unset()` — clear the value
### FieldAddOn (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldAddOn.IsSet()` — check if set
- `obj.FieldAddOn.Get()` — get the inner value (returns pointer)
- `obj.FieldAddOn.Set(&val)` — set the value
- `obj.FieldAddOn.Unset()` — clear the value
### FieldComponent (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldComponent.IsSet()` — check if set
- `obj.FieldComponent.Get()` — get the inner value (returns pointer)
- `obj.FieldComponent.Set(&val)` — set the value
- `obj.FieldComponent.Unset()` — clear the value
### FieldInput (Nullable)

Use the Nullable wrapper methods:
- `obj.FieldInput.IsSet()` — check if set
- `obj.FieldInput.Get()` — get the inner value (returns pointer)
- `obj.FieldInput.Set(&val)` — set the value
- `obj.FieldInput.Unset()` — clear the value
### PlaceHolder (Nullable)

Use the Nullable wrapper methods:
- `obj.PlaceHolder.IsSet()` — check if set
- `obj.PlaceHolder.Get()` — get the inner value (returns pointer)
- `obj.PlaceHolder.Set(&val)` — set the value
- `obj.PlaceHolder.Unset()` — clear the value
### VerifyPattern (Nullable)

Use the Nullable wrapper methods:
- `obj.VerifyPattern.IsSet()` — check if set
- `obj.VerifyPattern.Get()` — get the inner value (returns pointer)
- `obj.VerifyPattern.Set(&val)` — set the value
- `obj.VerifyPattern.Unset()` — clear the value
### HelpBlock (Nullable)

Use the Nullable wrapper methods:
- `obj.HelpBlock.IsSet()` — check if set
- `obj.HelpBlock.Get()` — get the inner value (returns pointer)
- `obj.HelpBlock.Set(&val)` — set the value
- `obj.HelpBlock.Unset()` — clear the value
### HelpBlockFieldCode (Nullable)

Use the Nullable wrapper methods:
- `obj.HelpBlockFieldCode.IsSet()` — check if set
- `obj.HelpBlockFieldCode.Get()` — get the inner value (returns pointer)
- `obj.HelpBlockFieldCode.Set(&val)` — set the value
- `obj.HelpBlockFieldCode.Unset()` — clear the value
### DefaultValue (Nullable)

Use the Nullable wrapper methods:
- `obj.DefaultValue.IsSet()` — check if set
- `obj.DefaultValue.Get()` — get the inner value (returns pointer)
- `obj.DefaultValue.Set(&val)` — set the value
- `obj.DefaultValue.Unset()` — clear the value
### OptionSource (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionSource.IsSet()` — check if set
- `obj.OptionSource.Get()` — get the inner value (returns pointer)
- `obj.OptionSource.Set(&val)` — set the value
- `obj.OptionSource.Unset()` — clear the value
### OptionSourceType (Nullable)

Use the Nullable wrapper methods:
- `obj.OptionSourceType.IsSet()` — check if set
- `obj.OptionSourceType.Get()` — get the inner value (returns pointer)
- `obj.OptionSourceType.Set(&val)` — set the value
- `obj.OptionSourceType.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### WrapperClass (Nullable)

Use the Nullable wrapper methods:
- `obj.WrapperClass.IsSet()` — check if set
- `obj.WrapperClass.Get()` — get the inner value (returns pointer)
- `obj.WrapperClass.Set(&val)` — set the value
- `obj.WrapperClass.Unset()` — clear the value
### DependsOnCode (Nullable)

Use the Nullable wrapper methods:
- `obj.DependsOnCode.IsSet()` — check if set
- `obj.DependsOnCode.Get()` — get the inner value (returns pointer)
- `obj.DependsOnCode.Set(&val)` — set the value
- `obj.DependsOnCode.Unset()` — clear the value
### VisibleOnCode (Nullable)

Use the Nullable wrapper methods:
- `obj.VisibleOnCode.IsSet()` — check if set
- `obj.VisibleOnCode.Get()` — get the inner value (returns pointer)
- `obj.VisibleOnCode.Set(&val)` — set the value
- `obj.VisibleOnCode.Unset()` — clear the value
### RequireOnCode (Nullable)

Use the Nullable wrapper methods:
- `obj.RequireOnCode.IsSet()` — check if set
- `obj.RequireOnCode.Get()` — get the inner value (returns pointer)
- `obj.RequireOnCode.Set(&val)` — set the value
- `obj.RequireOnCode.Unset()` — clear the value
### ContextualDefault (Nullable)

Use the Nullable wrapper methods:
- `obj.ContextualDefault.IsSet()` — check if set
- `obj.ContextualDefault.Get()` — get the inner value (returns pointer)
- `obj.ContextualDefault.Set(&val)` — set the value
- `obj.ContextualDefault.Unset()` — clear the value
### DisplayValueOnDetails (Nullable)

Use the Nullable wrapper methods:
- `obj.DisplayValueOnDetails.IsSet()` — check if set
- `obj.DisplayValueOnDetails.Get()` — get the inner value (returns pointer)
- `obj.DisplayValueOnDetails.Set(&val)` — set the value
- `obj.DisplayValueOnDetails.Unset()` — clear the value
### ShowOnCreate (Nullable)

Use the Nullable wrapper methods:
- `obj.ShowOnCreate.IsSet()` — check if set
- `obj.ShowOnCreate.Get()` — get the inner value (returns pointer)
- `obj.ShowOnCreate.Set(&val)` — set the value
- `obj.ShowOnCreate.Unset()` — clear the value
### ShowOnEdit (Nullable)

Use the Nullable wrapper methods:
- `obj.ShowOnEdit.IsSet()` — check if set
- `obj.ShowOnEdit.Get()` — get the inner value (returns pointer)
- `obj.ShowOnEdit.Set(&val)` — set the value
- `obj.ShowOnEdit.Unset()` — clear the value
### LocalCredential (Nullable)

Use the Nullable wrapper methods:
- `obj.LocalCredential.IsSet()` — check if set
- `obj.LocalCredential.Get()` — get the inner value (returns pointer)
- `obj.LocalCredential.Set(&val)` — set the value
- `obj.LocalCredential.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


