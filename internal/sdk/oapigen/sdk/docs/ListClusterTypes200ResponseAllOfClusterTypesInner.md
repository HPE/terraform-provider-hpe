# ListClusterTypes200ResponseAllOfClusterTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**DeployTargetService** | Pointer to **string** |  | [optional] 
**ShortName** | Pointer to **string** |  | [optional] 
**ProviderType** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**HostService** | Pointer to **string** |  | [optional] 
**Managed** | Pointer to **bool** |  | [optional] 
**HasMasters** | Pointer to **bool** |  | [optional] 
**HasWorkers** | Pointer to **bool** |  | [optional] 
**ViewSet** | Pointer to **string** |  | [optional] 
**ImageCode** | Pointer to **string** |  | [optional] 
**KubeCtlLocal** | Pointer to **bool** |  | [optional] 
**HasDatastore** | Pointer to **bool** |  | [optional] 
**SupportsCloudScaling** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**HasDefaultDataDisk** | Pointer to **bool** |  | [optional] 
**CanManage** | Pointer to **bool** |  | [optional] 
**HasCluster** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**OptionTypes** | Pointer to [**[]ListClusterTypes200ResponseAllOfClusterTypesInnerOptionTypesInner**](ListClusterTypes200ResponseAllOfClusterTypesInnerOptionTypesInner.md) |  | [optional] 
**ControllerTypes** | Pointer to [**[]ListClusterTypes200ResponseAllOfClusterTypesInnerControllerTypesInner**](ListClusterTypes200ResponseAllOfClusterTypesInnerControllerTypesInner.md) |  | [optional] 
**WorkerTypes** | Pointer to [**[]ListClusterTypes200ResponseAllOfClusterTypesInnerWorkerTypesInner**](ListClusterTypes200ResponseAllOfClusterTypesInnerWorkerTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterTypes200ResponseAllOfClusterTypesInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


