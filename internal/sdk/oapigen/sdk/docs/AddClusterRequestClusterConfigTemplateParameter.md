# AddClusterRequestClusterConfigTemplateParameter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DnsPrefix** | Pointer to **string** | Optional DNS prefix to use with hosted Kubernetes API server FQDN. | [optional] 
**OnDiskSizeGB** | Pointer to **int64** | Disk size (in GB) to provision for each of the agent pool nodes. This value ranges from 0 to 1023. Specifying 0 will apply the default disk size for that agentVMSize. | [optional] 
**NodeVMSize** | Pointer to **int64** | The size of the Virtual Machine. | [optional] 
**VnetSubnetID** | Pointer to **string** | Resource ID of virtual network subnet used for nodes and/or pods IP assignment. | [optional] 
**ServiceCidr** | Pointer to **string** | A CIDR notation IP range from which to assign service cluster IPs. | [optional] 
**DnsServiceIP** | Pointer to **string** | Containers DNS server IP address. | [optional] 
**DockerBridgeCidr** | Pointer to **string** | A CIDR notation IP for Docker bridge. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterRequestClusterConfigTemplateParameter{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


