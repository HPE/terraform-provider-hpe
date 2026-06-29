# GetBlueprint200ResponseBlueprint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**ResourcePermission** | Pointer to **map[string]interface{}** |  | [optional] 
**Owner** | Pointer to [**GetBlueprint200ResponseBlueprintOwner**](GetBlueprint200ResponseBlueprintOwner.md) |  | [optional] 
**Tenant** | Pointer to [**GetBlueprint200ResponseBlueprintTenant**](GetBlueprint200ResponseBlueprintTenant.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetBlueprint200ResponseBlueprint{
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
### ResourcePermission (Nullable)

Use the Nullable wrapper methods:
- `obj.ResourcePermission.IsSet()` — check if set
- `obj.ResourcePermission.Get()` — get the inner value (returns pointer)
- `obj.ResourcePermission.Set(&val)` — set the value
- `obj.ResourcePermission.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


