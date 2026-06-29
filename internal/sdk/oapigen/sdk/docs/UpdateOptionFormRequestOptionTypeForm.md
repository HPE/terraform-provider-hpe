# UpdateOptionFormRequestOptionTypeForm

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Form name | [optional] 
**Code** | Pointer to **string** | Unique form code | [optional] 
**Description** | Pointer to **NullableString** | A short description of the form | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Options** | Pointer to [**[]UpdateOptionFormRequestOptionTypeFormOptionsInner**](UpdateOptionFormRequestOptionTypeFormOptionsInner.md) | Inputs | [optional] 
**FieldGroups** | Pointer to [**[]UpdateOptionFormRequestOptionTypeFormFieldGroupsInner**](UpdateOptionFormRequestOptionTypeFormFieldGroupsInner.md) | Field Groups | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateOptionFormRequestOptionTypeForm{
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


