# ZoneNetworkOptionsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Networks** | Pointer to [**[]ZoneNetworkOptionsResponseNetworksInner**](ZoneNetworkOptionsResponseNetworksInner.md) |  | [optional] 
**NetworkGroups** | Pointer to [**[]ZoneNetworkOptionsResponseNetworkGroupsInner**](ZoneNetworkOptionsResponseNetworkGroupsInner.md) |  | [optional] 
**NetworkTypes** | Pointer to [**[]ZoneNetworkOptionsResponseNetworkTypesInner**](ZoneNetworkOptionsResponseNetworkTypesInner.md) |  | [optional] 
**NetworkSubnets** | Pointer to [**[]ZoneNetworkOptionsResponseNetworkSubnetsInner**](ZoneNetworkOptionsResponseNetworkSubnetsInner.md) |  | [optional] 
**HasNetworks** | Pointer to **NullableBool** |  | [optional] 
**MaxNetworks** | Pointer to **NullableInt64** |  | [optional] 
**EnableNetworkTypeSelection** | Pointer to **NullableString** |  | [optional] 
**SupportsNetworkSelection** | Pointer to **NullableBool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneNetworkOptionsResponse{
    // Set fields directly
}
```

### HasNetworks (Nullable)

Use the Nullable wrapper methods:
- `obj.HasNetworks.IsSet()` — check if set
- `obj.HasNetworks.Get()` — get the inner value (returns pointer)
- `obj.HasNetworks.Set(&val)` — set the value
- `obj.HasNetworks.Unset()` — clear the value
### MaxNetworks (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxNetworks.IsSet()` — check if set
- `obj.MaxNetworks.Get()` — get the inner value (returns pointer)
- `obj.MaxNetworks.Set(&val)` — set the value
- `obj.MaxNetworks.Unset()` — clear the value
### EnableNetworkTypeSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.EnableNetworkTypeSelection.IsSet()` — check if set
- `obj.EnableNetworkTypeSelection.Get()` — get the inner value (returns pointer)
- `obj.EnableNetworkTypeSelection.Set(&val)` — set the value
- `obj.EnableNetworkTypeSelection.Unset()` — clear the value
### SupportsNetworkSelection (Nullable)

Use the Nullable wrapper methods:
- `obj.SupportsNetworkSelection.IsSet()` — check if set
- `obj.SupportsNetworkSelection.Get()` — get the inner value (returns pointer)
- `obj.SupportsNetworkSelection.Set(&val)` — set the value
- `obj.SupportsNetworkSelection.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


