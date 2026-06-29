# InstancesNetworkInterfaces1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Network** | [**InstancesNetworkInterfaces1Network**](InstancesNetworkInterfaces1Network.md) |  | 
**NetworkInterfaceTypeId** | Pointer to **int64** | The id of type of the network interface. | [optional] 
**IpMode** | Pointer to **string** | The mode for determining ip address. Can be &#39;static&#39;, &#39;dhcp&#39; or empty string. | [optional] [default to ""]
**IpAddress** | Pointer to **string** | The ip address. Not applicable when using DHCP or IP Pools. | [optional] 
**MacAddress** | Pointer to **string** | The MAC address. | [optional] 
**Id** | Pointer to **int64** | The interface id. Applicable when resizing and you want to identify an interface to update that already exists. | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces1NetworkInterfacesInner**](InstancesNetworkInterfaces1NetworkInterfacesInner.md) | The nested networkInterfaces can be used to define child virtual network interfaces. The Options API &#x60;/api/options/zoneNetworkOptions?zoneId&#x3D;5&amp;provisionTypeId&#x3D;10&#x60; can be used to see which types support this (&#x60;hasVirtualInvirtualInterfaces &#x3D; true&#x60; and list of available &#x60;virtualInterfaces&#x60; will be defined.  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstancesNetworkInterfaces1{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


