# AddOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | The ID of the input, include this to use an existing input or to update an existing form input record instead of creating a new one. | [optional] 
**Code** | Pointer to **string** | The code of the option type as a globally unique identifier. By default a UUID will be used. | [optional] 
**Name** | Pointer to **string** | The name of the option type for handy reference. By default a UUID will be used. | [optional] 
**Description** | Pointer to **NullableString** | Short description of the option type | [optional] 
**FieldName** | **string** | Field Name, the name for user input. This along with fieldContext determines the configuration property name.  The property key for when posting this option type to a JSON POST request | 
**Type** | Pointer to **string** | Type, the type of input. eg. text, checkbox, select, etc. | [optional] [default to "text"]
**FieldLabel** | **string** | Field Label, the label for user input. | 
**FieldCode** | Pointer to **NullableString** | Localized Label, i18n code for the label | [optional] 
**PlaceHolder** | Pointer to **string** | Any placeholder text when nothing is yet entered | [optional] 
**HelpBlock** | Pointer to **NullableString** | This is the explaination of the input that shows typically underneath the option | [optional] 
**HelpBlockFieldCode** | Pointer to **NullableString** | Localized Help Block, i18n code for the help block | [optional] 
**DefaultValue** | Pointer to **string** | The default value if no user entry is specified. This value should be passed to the desired JSON Map if nothing else is entered | [optional] 
**Required** | Pointer to **bool** | Is this field entry required for the request | [optional] [default to false]
**ExportMeta** | Pointer to **bool** | Export as Tag | [optional] [default to false]
**Editable** | Pointer to **bool** | Used primarily on tasks and workflows. Basically wether or not the field can be overridden optionally when the object is run | [optional] [default to false]
**OptionList** | Pointer to [**AddOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInnerOptionList**](AddOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInnerOptionList.md) |  | [optional] 
**DisplayValueOnDetails** | Pointer to **bool** | Display Value On Details | [optional] [default to false]
**IsLocked** | Pointer to **bool** | Locked | [optional] [default to false]
**IsHidden** | Pointer to **bool** | Hidden | [optional] [default to false]
**ExcludeFromSearch** | Pointer to **bool** | Exclude From Search | [optional] [default to false]
**DependsOnCode** | Pointer to **NullableString** | A fieldName that will trigger reloading this input | [optional] 
**VisibleOnCode** | Pointer to **NullableString** | A fieldName that will trigger visibility of this input | [optional] 
**VerifyPattern** | Pointer to **string** | Verify Pattern, A regexp string that validates the input, use (?i) to make the matcher case insensitive | [optional] 
**RequireOnCode** | Pointer to **NullableString** | A fieldName that will trigger required attribute of this input | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInner{
    // Set fields directly
}
```

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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


