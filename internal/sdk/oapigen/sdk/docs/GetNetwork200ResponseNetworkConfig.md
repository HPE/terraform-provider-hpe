# GetNetwork200ResponseNetworkConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VlanIDs** | Pointer to **NullableString** |  | [optional] 
**ConnectedGateway** | Pointer to **string** |  | [optional] 
**SubnetIpManagementType** | Pointer to **string** |  | [optional] 
**SubnetIpServerId** | Pointer to **string** |  | [optional] 
**DhcpRange** | Pointer to **string** |  | [optional] 
**SubnetDhcpServerAddress** | Pointer to **string** |  | [optional] 
**SubnetDhcpLeaseTime** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetwork200ResponseNetworkConfig{
    // Set fields directly
}
```

### VlanIDs (Nullable)

Use the Nullable wrapper methods:
- `obj.VlanIDs.IsSet()` — check if set
- `obj.VlanIDs.Get()` — get the inner value (returns pointer)
- `obj.VlanIDs.Set(&val)` — set the value
- `obj.VlanIDs.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


