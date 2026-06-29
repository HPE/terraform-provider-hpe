# InstanceClone

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the new cloned instance. If none is specified the existing name will be duplicated with the &#39;clone&#39; suffix added. | [optional] 
**Group** | Pointer to [**InstanceCloneGroup**](InstanceCloneGroup.md) |  | [optional] 
**Cloud** | Pointer to [**InstanceCloneCloud**](InstanceCloneCloud.md) |  | [optional] 
**Plan** | Pointer to [**InstanceClonePlan**](InstanceClonePlan.md) |  | [optional] 
**UserGroup** | Pointer to [**InstanceCloneUserGroup**](InstanceCloneUserGroup.md) |  | [optional] 
**Volumes** | Pointer to [**[]InstanceCloneVolumesInner**](InstanceCloneVolumesInner.md) | The complete set of volumes for the clone. When provided, this replaces (does not merge with) the source instance&#39;s volumes.  Volume specifications define infrastructure parameters (size, datastore, type) only - the disk contents are always copied from the source instance.  | [optional] 
**NetworkInterfaces** | Pointer to [**[]InstancesNetworkInterfaces10**](InstancesNetworkInterfaces10.md) | The complete set of network interfaces for the clone. When provided, this replaces (does not merge with) the source instance&#39;s network interfaces.  The Options API &#x60;/api/options/zoneNetworkOptions?zoneId&#x3D;5&amp;provisionTypeId&#x3D;10&#x60; can be used to see which options are available.  | [optional] 
**Config** | Pointer to [**InstanceCloneConfig**](InstanceCloneConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceClone{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


