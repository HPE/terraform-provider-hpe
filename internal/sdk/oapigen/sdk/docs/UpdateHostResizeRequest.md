# UpdateHostResizeRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Server** | Pointer to [**UpdateHostResizeRequestServer**](UpdateHostResizeRequestServer.md) |  | [optional] 
**ServicePlanOptions** | Pointer to [**UpdateHostResizeRequestServicePlanOptions**](UpdateHostResizeRequestServicePlanOptions.md) |  | [optional] 
**Volumes** | Pointer to [**[]UpdateHostResizeRequestVolumesInner**](UpdateHostResizeRequestVolumesInner.md) | List of volumes with their new sizes. | [optional] 
**DeleteOriginalVolumes** | Pointer to **bool** | Delete the original volumes after resizing. (Amazon only) | [optional] [default to false]
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces5**](InstancesNetworkInterfaces5.md) | Key for network configurations. Include id to update an existing interface. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateHostResizeRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


