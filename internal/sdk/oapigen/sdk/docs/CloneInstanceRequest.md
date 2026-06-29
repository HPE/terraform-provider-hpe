# CloneInstanceRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the new cloned instance. If none is specified the existing name will be duplicated with the &#39;clone&#39; suffix added. | [optional] 
**Group** | Pointer to [**CloneInstanceRequestGroup**](CloneInstanceRequestGroup.md) |  | [optional] 
**Cloud** | Pointer to [**CloneInstanceRequestCloud**](CloneInstanceRequestCloud.md) |  | [optional] 
**Plan** | Pointer to [**CloneInstanceRequestPlan**](CloneInstanceRequestPlan.md) |  | [optional] 
**UserGroup** | Pointer to [**CloneInstanceRequestUserGroup**](CloneInstanceRequestUserGroup.md) |  | [optional] 
**Volumes** | Pointer to [**[]CloneInstanceRequestVolumesInner**](CloneInstanceRequestVolumesInner.md) | The complete set of volumes for the clone. When provided, this replaces (does not merge with) the source instance&#39;s volumes.  Volume specifications define infrastructure parameters (size, datastore, type) only - the disk contents are always copied from the source instance.  | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces3**](InstancesNetworkInterfaces3.md) | The complete set of network interfaces for the clone. When provided, this replaces (does not merge with) the source instance&#39;s network interfaces.  The Options API &#x60;/api/options/zoneNetworkOptions?zoneId&#x3D;5&amp;provisionTypeId&#x3D;10&#x60; can be used to see which options are available.  | [optional] 
**Config** | Pointer to [**CloneInstanceRequestConfig**](CloneInstanceRequestConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CloneInstanceRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


