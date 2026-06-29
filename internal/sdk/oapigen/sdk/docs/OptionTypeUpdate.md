# OptionTypeUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the option type for handy reference | [optional] 
**Description** | Pointer to **NullableString** | Short description of the option type | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**FieldName** | Pointer to **string** | Field Name, the name for user input. This along with fieldContext determines the configuration property name.  The property key for when posting this option type to a JSON POST request | [optional] 
**Type** | Pointer to **string** | Type, the type of input. eg. text, checkbox, select, etc. | [optional] [default to "text"]
**FieldLabel** | Pointer to **string** | Field Label, the label for user input. | [optional] 
**PlaceHolder** | Pointer to **string** | Any placeholder text when nothing is yet entered | [optional] 
**VerifyPattern** | Pointer to **string** | Verify Pattern, A regexp string that validates the input, use (?i) to make the matcher case insensitive | [optional] 
**DefaultValue** | Pointer to **string** | The default value if no user entry is specified. This value should be passed to the desired JSON Map if nothing else is entered | [optional] 
**Required** | Pointer to **bool** | Is this field entry required for the request | [optional] [default to false]
**ExportMeta** | Pointer to **bool** | Export as Tag | [optional] [default to false]
**Editable** | Pointer to **bool** | Used primarily on tasks and workflows. Basically wether or not the field can be overridden optionally when the object is run | [optional] [default to false]
**OptionList** | Pointer to [**OptionTypeUpdateOptionList**](OptionTypeUpdateOptionList.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &OptionTypeUpdate{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


