# GetGroups200ResponseGroup

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
**Account** | Pointer to [**GetGroups200ResponseGroupAccount**](GetGroups200ResponseGroupAccount.md) |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**Config** | Pointer to [**GetGroups200ResponseGroupConfig**](GetGroups200ResponseGroupConfig.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Zones** | Pointer to [**[]GetGroups200ResponseGroupZonesInner**](GetGroups200ResponseGroupZonesInner.md) |  | [optional] 
**Stats** | Pointer to [**GetGroups200ResponseGroupStats**](GetGroups200ResponseGroupStats.md) |  | [optional] 
**ServerCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetGroups200ResponseGroup{
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


