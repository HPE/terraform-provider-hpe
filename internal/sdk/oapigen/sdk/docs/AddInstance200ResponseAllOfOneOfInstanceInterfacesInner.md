# AddInstance200ResponseAllOfOneOfInstanceInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to [**AddInstance200ResponseAllOfOneOfInstanceInterfacesInnerId**](AddInstance200ResponseAllOfOneOfInstanceInterfacesInnerId.md) |  | [optional] 
**Network** | Pointer to [**AddInstance200ResponseAllOfOneOfInstanceInterfacesInnerNetwork**](AddInstance200ResponseAllOfOneOfInstanceInterfacesInnerNetwork.md) |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**NetworkInterfaceTypeId** | Pointer to **int64** |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstanceInterfacesNetworkInterfacesInner1**](InstanceInterfacesNetworkInterfacesInner1.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddInstance200ResponseAllOfOneOfInstanceInterfacesInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


