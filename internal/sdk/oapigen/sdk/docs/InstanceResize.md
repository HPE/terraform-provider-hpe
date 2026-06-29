# InstanceResize

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instance** | Pointer to [**InstanceResizeInstance**](InstanceResizeInstance.md) |  | [optional] 
**ServicePlanOptions** | Pointer to [**InstanceResizeServicePlanOptions**](InstanceResizeServicePlanOptions.md) |  | [optional] 
**Volumes** | Pointer to [**[]InstanceResizeVolumesInner**](InstanceResizeVolumesInner.md) | Can be used to grow just the logical volume of the instance instead of choosing a plan | [optional] 
**DeleteOriginalVolumes** | Pointer to **bool** | Delete the original volumes after resizing. (Amazon only) | [optional] [default to false]
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces11**](InstancesNetworkInterfaces11.md) | Key for network configuration. Include id to update an existing interface. The existing interfaces and their id can be retrieved with the hosts API. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceResize{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


