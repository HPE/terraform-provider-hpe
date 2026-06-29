# CreateTenantSubtenantGroup200ResponseAccount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Location** | Pointer to **NullableString** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Zones** | Pointer to [**[]CreateTenantSubtenantGroup200ResponseAccountAllOfZonesInner**](CreateTenantSubtenantGroup200ResponseAccountAllOfZonesInner.md) |  | [optional] 
**Stats** | Pointer to [**CreateTenantSubtenantGroup200ResponseAccountAllOfStats**](CreateTenantSubtenantGroup200ResponseAccountAllOfStats.md) |  | [optional] 
**ServerCount** | Pointer to **int64** |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateTenantSubtenantGroup200ResponseAccount{
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


