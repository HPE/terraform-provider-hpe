# ClusterCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | [**ClusterCreateType**](ClusterCreateType.md) |  | 
**Name** | **string** | Name of the cluster to be created | 
**Description** | Pointer to **string** | Description of the cluster to be created | [optional] 
**Labels** | Pointer to **[]string** | Array of strings (keywords). This will override labels passed under the &#x60;server&#x60; object. | [optional] 
**Group** | Pointer to [**ClusterCreateGroup**](ClusterCreateGroup.md) |  | [optional] 
**Cloud** | [**ClusterCreateCloud**](ClusterCreateCloud.md) |  | 
**Config** | Pointer to [**ClusterCreateConfig**](ClusterCreateConfig.md) |  | [optional] 
**Layout** | [**ClusterCreateLayout**](ClusterCreateLayout.md) |  | 
**Server** | [**ClusterCreateServer**](ClusterCreateServer.md) |  | 
**AutoRecoverPowerState** | Pointer to **bool** | Automatically Power on VMs | [optional] [default to false]
**TaskSetId** | Pointer to **int64** | Optional Workflow Id desired to be run | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


