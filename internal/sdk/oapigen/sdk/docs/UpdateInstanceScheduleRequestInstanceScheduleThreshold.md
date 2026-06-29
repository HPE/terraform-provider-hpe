# UpdateInstanceScheduleRequestInstanceScheduleThreshold

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SourceThresholdId** | Pointer to **int64** | Source Scale Threshold to apply as a template. All threshold settings with be copied from this threshold, and can be overridden by also passing each setting explicitly. | [optional] 
**AutoUp** | Pointer to **bool** | Auto Upscale | [optional] 
**AutoDown** | Pointer to **bool** | Auto Downscale | [optional] 
**MinCount** | Pointer to **int32** | The minimum number of nodes to scale down to | [optional] 
**MaxCount** | Pointer to **int32** | The maximum number of nodes to scale up to | [optional] 
**CpuEnabled** | Pointer to **bool** | Enable CPU Threshold | [optional] 
**MinCpu** | Pointer to **float64** | Min CPU (%) | [optional] 
**MaxCpu** | Pointer to **float64** | Max CPU (%) | [optional] 
**MemoryEnabled** | Pointer to **bool** | Enable Memory Threshold | [optional] 
**MinMemory** | Pointer to **float64** | Min Memory (%) | [optional] 
**MaxMemory** | Pointer to **float64** | Max Memory (%) | [optional] 
**DiskEnabled** | Pointer to **bool** | Enable Disk Threshold | [optional] 
**MinDisk** | Pointer to **float64** | Min Disk (%) | [optional] 
**MaxDisk** | Pointer to **float64** | Max Disk (%) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInstanceScheduleRequestInstanceScheduleThreshold{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


