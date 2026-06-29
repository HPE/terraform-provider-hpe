# GetNetworkPoolServerType200ResponseNetworkPoolServerType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**PoolService** | Pointer to **NullableString** |  | [optional] 
**Selectable** | Pointer to **bool** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**IntegrationCode** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetNetworkPoolServerType200ResponseNetworkPoolServerTypeOptionTypesInner**](GetNetworkPoolServerType200ResponseNetworkPoolServerTypeOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkPoolServerType200ResponseNetworkPoolServerType{
    // Set fields directly
}
```

### PoolService (Nullable)

Use the Nullable wrapper methods:
- `obj.PoolService.IsSet()` — check if set
- `obj.PoolService.Get()` — get the inner value (returns pointer)
- `obj.PoolService.Set(&val)` — set the value
- `obj.PoolService.Unset()` — clear the value
### IntegrationCode (Nullable)

Use the Nullable wrapper methods:
- `obj.IntegrationCode.IsSet()` — check if set
- `obj.IntegrationCode.Get()` — get the inner value (returns pointer)
- `obj.IntegrationCode.Set(&val)` — set the value
- `obj.IntegrationCode.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


