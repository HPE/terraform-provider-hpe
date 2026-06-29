# AddClusterRequestCluster

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | [**AddClusterRequestClusterType**](AddClusterRequestClusterType.md) |  | 
**Name** | **string** | Name of the cluster to be created | 
**Description** | Pointer to **string** | Description of the cluster to be created | [optional] 
**Labels** | Pointer to **[]string** | Array of strings (keywords). This will override labels passed under the &#x60;server&#x60; object. | [optional] 
**Group** | Pointer to [**AddClusterRequestClusterGroup**](AddClusterRequestClusterGroup.md) |  | [optional] 
**Cloud** | [**AddClusterRequestClusterCloud**](AddClusterRequestClusterCloud.md) |  | 
**Config** | Pointer to [**AddClusterRequestClusterConfig**](AddClusterRequestClusterConfig.md) |  | [optional] 
**Layout** | [**AddClusterRequestClusterLayout**](AddClusterRequestClusterLayout.md) |  | 
**Server** | [**AddClusterRequestClusterServer**](AddClusterRequestClusterServer.md) |  | 
**AutoRecoverPowerState** | Pointer to **bool** | Automatically Power on VMs | [optional] [default to false]
**TaskSetId** | Pointer to **int64** | Optional Workflow Id desired to be run | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterRequestCluster{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


