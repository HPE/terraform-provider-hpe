# ClusterContainersStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ts** | Pointer to **time.Time** |  | [optional] 
**Running** | Pointer to **bool** |  | [optional] 
**UserCpuUsage** | Pointer to **float32** |  | [optional] 
**SystemCpuUsage** | Pointer to **float32** |  | [optional] 
**UsedMemory** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**CacheMemory** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**UsedStorage** | Pointer to **int64** |  | [optional] 
**ReadIOPS** | Pointer to **int64** |  | [optional] 
**WriteIOPS** | Pointer to **int64** |  | [optional] 
**TotalIOPS** | Pointer to **int64** |  | [optional] 
**NetTxUsage** | Pointer to **int64** |  | [optional] 
**NetRxUsage** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterContainersStats{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


