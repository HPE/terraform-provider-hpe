# NetworkCreateConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ResourceGroupId** | **string** | Resource Group Name | 
**SubnetName** | **string** | Subnet Name | 
**SubnetCidr** | **string** | The subnet&#39;s address range in CIDR notation (e.g. 192.168.1.0/24). It must be contained by the address space of the virtual network. | 
**AvailabilityZone** | **string** | Availability Zone Name | 
**Cidr** | **string** | Network CIDR | 
**AssignPublicIp** | **bool** | Assign public IPs by default. | 
**ZonePool** | [**NetworkCreateConfigAnyOf1ZonePool**](NetworkCreateConfigAnyOf1ZonePool.md) |  | 
**Mtu** | **string** | GCP MTU | [default to "1460"]
**AutoCreate** | **bool** | Auto create subnets | [default to true]

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkCreateConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


