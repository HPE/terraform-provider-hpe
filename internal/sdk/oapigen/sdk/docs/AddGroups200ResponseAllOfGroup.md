# AddGroups200ResponseAllOfGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Location** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**AddGroups200ResponseAllOfGroupAccount**](AddGroups200ResponseAllOfGroupAccount.md) |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**AddGroups200ResponseAllOfGroupConfig**](AddGroups200ResponseAllOfGroupConfig.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Zones** | Pointer to [**[]AddGroups200ResponseAllOfGroupZonesInner**](AddGroups200ResponseAllOfGroupZonesInner.md) |  | [optional] 
**Stats** | Pointer to [**AddGroups200ResponseAllOfGroupStats**](AddGroups200ResponseAllOfGroupStats.md) |  | [optional] 
**ServerCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddGroups200ResponseAllOfGroup{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Location (Nullable)

Use the Nullable wrapper methods:
- `obj.Location.IsSet()` — check if set
- `obj.Location.Get()` — get the inner value (returns pointer)
- `obj.Location.Set(&val)` — set the value
- `obj.Location.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


