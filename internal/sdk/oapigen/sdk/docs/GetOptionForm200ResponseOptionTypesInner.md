# GetOptionForm200ResponseOptionTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Context** | Pointer to **NullableString** |  | [optional] 
**Locked** | Pointer to **bool** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Options** | Pointer to [**[]GetOptionForm200ResponseOptionTypesInnerOptionsInner**](GetOptionForm200ResponseOptionTypesInnerOptionsInner.md) |  | [optional] 
**FieldGroups** | Pointer to [**[]GetOptionForm200ResponseOptionTypesInnerFieldGroupsInner**](GetOptionForm200ResponseOptionTypesInnerFieldGroupsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetOptionForm200ResponseOptionTypesInner{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Context (Nullable)

Use the Nullable wrapper methods:
- `obj.Context.IsSet()` — check if set
- `obj.Context.Get()` — get the inner value (returns pointer)
- `obj.Context.Set(&val)` — set the value
- `obj.Context.Unset()` — clear the value
### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


