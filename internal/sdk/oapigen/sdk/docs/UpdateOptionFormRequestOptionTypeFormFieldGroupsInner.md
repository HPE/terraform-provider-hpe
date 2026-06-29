# UpdateOptionFormRequestOptionTypeFormFieldGroupsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**LocalizedName** | Pointer to **NullableString** |  | [optional] 
**Collapsible** | Pointer to **bool** |  | [optional] 
**DefaultCollapsed** | Pointer to **bool** |  | [optional] 
**VisibleOnCode** | Pointer to **NullableString** |  | [optional] 
**Options** | Pointer to [**[]UpdateOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInner**](UpdateOptionFormRequestOptionTypeFormFieldGroupsInnerOptionsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateOptionFormRequestOptionTypeFormFieldGroupsInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### LocalizedName (Nullable)

Use the Nullable wrapper methods:
- `obj.LocalizedName.IsSet()` — check if set
- `obj.LocalizedName.Get()` — get the inner value (returns pointer)
- `obj.LocalizedName.Set(&val)` — set the value
- `obj.LocalizedName.Unset()` — clear the value
### VisibleOnCode (Nullable)

Use the Nullable wrapper methods:
- `obj.VisibleOnCode.IsSet()` — check if set
- `obj.VisibleOnCode.Get()` — get the inner value (returns pointer)
- `obj.VisibleOnCode.Set(&val)` — set the value
- `obj.VisibleOnCode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


