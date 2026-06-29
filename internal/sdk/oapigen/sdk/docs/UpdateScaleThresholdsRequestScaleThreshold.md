# UpdateScaleThresholdsRequestScaleThreshold

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the scale threshold | [optional] 
**AutoUp** | Pointer to **bool** | Auto Upscale | [optional] [default to false]
**AutoDown** | Pointer to **bool** | Auto Downscale | [optional] [default to false]
**MinCount** | Pointer to **int32** | The minimum number of nodes to scale down to | [optional] 
**MaxCount** | Pointer to **int32** | The maximum number of nodes to scale up to | [optional] 
**CpuEnabled** | Pointer to **bool** | Enable CPU Threshold | [optional] [default to false]
**MinCpu** | Pointer to **float32** | Min CPU (%) | [optional] [default to 0]
**MaxCpu** | Pointer to **float32** | Max CPU (%) | [optional] [default to 0]
**MemoryEnabled** | Pointer to **bool** | Enable Memory Threshold | [optional] [default to false]
**MinMemory** | Pointer to **float32** | Min Memory (%) | [optional] [default to 0]
**MaxMemory** | Pointer to **float32** | Max Memory (%) | [optional] [default to 0]
**DiskEnabled** | Pointer to **bool** | Enable Disk Threshold | [optional] [default to false]
**MinDisk** | Pointer to **float32** | Min Disk (%) | [optional] [default to 0]
**MaxDisk** | Pointer to **float32** | Max Disk (%) | [optional] [default to 0]

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateScaleThresholdsRequestScaleThreshold{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


