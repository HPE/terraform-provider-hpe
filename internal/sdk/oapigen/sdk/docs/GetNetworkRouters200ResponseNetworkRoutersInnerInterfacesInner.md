# GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**InterfaceType** | Pointer to **NullableString** |  | [optional] 
**NetworkPosition** | Pointer to **NullableString** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Cidr** | Pointer to **string** |  | [optional] 
**ExternalLink** | Pointer to **NullableString** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Network** | Pointer to [**GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInnerNetwork**](GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInnerNetwork.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkRouters200ResponseNetworkRoutersInnerInterfacesInner{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### InterfaceType (Nullable)

Use the Nullable wrapper methods:
- `obj.InterfaceType.IsSet()` — check if set
- `obj.InterfaceType.Get()` — get the inner value (returns pointer)
- `obj.InterfaceType.Set(&val)` — set the value
- `obj.InterfaceType.Unset()` — clear the value
### NetworkPosition (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPosition.IsSet()` — check if set
- `obj.NetworkPosition.Get()` — get the inner value (returns pointer)
- `obj.NetworkPosition.Set(&val)` — set the value
- `obj.NetworkPosition.Unset()` — clear the value
### ExternalLink (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalLink.IsSet()` — check if set
- `obj.ExternalLink.Get()` — get the inner value (returns pointer)
- `obj.ExternalLink.Set(&val)` — set the value
- `obj.ExternalLink.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


