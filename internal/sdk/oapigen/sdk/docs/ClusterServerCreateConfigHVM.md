# ClusterServerCreateConfigHVM

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CpuArch** | Pointer to **string** |  | [optional] 
**CpuModel** | Pointer to **string** |  | [optional] 
**DynamicPlacementMode** | Pointer to **string** | When enabled, Dynamic Placement will automatically balance VMs across cluster hosts based on resource utilization. When disabled, VMs will only migrate to a new host if they are pinned to a specific host or failed over and not running on the preferred host. | [optional] 
**PowerPolicy** | Pointer to **string** |  | [optional] 
**VcpuPlacementMode** | Pointer to **string** |  | [optional] 
**StorageInterfaceName** | Pointer to **string** |  | [optional] 
**ComputeInterfaceName** | Pointer to **string** |  | [optional] 
**ComputeVlans** | Pointer to **string** |  | [optional] 
**OverlayInterfaceName** | Pointer to **string** |  | [optional] 
**CreateUser** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterServerCreateConfigHVM{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


