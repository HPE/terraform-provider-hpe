# GetInstanceStats200ResponseInstanceStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaxCpu** | Pointer to **float32** | Total CPUs for all instances | [optional] 
**MaxCores** | Pointer to **float32** | Total cores for all instances | [optional] 
**CpuUsageAverage** | Pointer to **float32** | Average current CPU Usage across all instances | [optional] 
**CpuUsagePeak** | Pointer to **float32** | Max current CPU Usage across all instances | [optional] 
**UsedMemory** | Pointer to **float32** | Total used memory across all instances | [optional] 
**MaxMemory** | Pointer to **float32** | Total memory across all instances | [optional] 
**UsedStorage** | Pointer to **float32** | Total used storage total across all instances | [optional] 
**MaxStorage** | Pointer to **float32** | Total storage across all instances | [optional] 
**Running** | Pointer to **float32** | Total number of running instances | [optional] 
**Total** | Pointer to **float32** | Total number of instances | [optional] 
**TotalContainers** | Pointer to **float32** | Total number of containers across all instances | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceStats200ResponseInstanceStats{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


