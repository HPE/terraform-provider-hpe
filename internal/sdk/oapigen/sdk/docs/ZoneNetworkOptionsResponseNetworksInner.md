# ZoneNetworkOptionsResponseNetworksInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DhcpServer** | Pointer to **bool** |  | [optional] 
**AllowStaticOverride** | Pointer to **bool** |  | [optional] 
**Pool** | Pointer to [**NullableZoneNetworkOptionsResponseNetworksInnerPool**](ZoneNetworkOptionsResponseNetworksInnerPool.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneNetworkOptionsResponseNetworksInner{
    // Set fields directly
}
```

### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


