# AddBlueprint200ResponseAllOfBlueprint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | Blueprint ID | [optional] 
**Name** | Pointer to **string** | A name for the blueprint | [optional] 
**Description** | Pointer to **NullableString** | A description for the blueprint | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Category** | Pointer to **NullableString** | Category | [optional] 
**Config** | Pointer to [**AddBlueprint200ResponseAllOfBlueprintConfig**](AddBlueprint200ResponseAllOfBlueprintConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddBlueprint200ResponseAllOfBlueprint{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


