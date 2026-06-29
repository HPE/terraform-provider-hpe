# ListInstances200ResponseAllOfInstancesInnerInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerInterfacesInnerId**](ListInstances200ResponseAllOfInstancesInnerInterfacesInnerId.md) |  | [optional] 
**Network** | Pointer to [**ListInstances200ResponseAllOfInstancesInnerInterfacesInnerNetwork**](ListInstances200ResponseAllOfInstancesInnerInterfacesInnerNetwork.md) |  | [optional] 
**IpAddress** | Pointer to **NullableString** |  | [optional] 
**NetworkInterfaceTypeId** | Pointer to **NullableInt64** |  | [optional] 
**IpMode** | Pointer to **NullableString** |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstanceInterfacesNetworkInterfacesInner**](InstanceInterfacesNetworkInterfacesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListInstances200ResponseAllOfInstancesInnerInterfacesInner{
    // Set fields directly
}
```

### IpAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.IpAddress.IsSet()` — check if set
- `obj.IpAddress.Get()` — get the inner value (returns pointer)
- `obj.IpAddress.Set(&val)` — set the value
- `obj.IpAddress.Unset()` — clear the value
### NetworkInterfaceTypeId (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkInterfaceTypeId.IsSet()` — check if set
- `obj.NetworkInterfaceTypeId.Get()` — get the inner value (returns pointer)
- `obj.NetworkInterfaceTypeId.Set(&val)` — set the value
- `obj.NetworkInterfaceTypeId.Unset()` — clear the value
### IpMode (Nullable)

Use the Nullable wrapper methods:
- `obj.IpMode.IsSet()` — check if set
- `obj.IpMode.Get()` — get the inner value (returns pointer)
- `obj.IpMode.Set(&val)` — set the value
- `obj.IpMode.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


