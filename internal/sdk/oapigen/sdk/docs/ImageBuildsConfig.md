# ImageBuildsConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Instance** | Pointer to [**ImageBuildsConfigInstance**](ImageBuildsConfigInstance.md) |  | [optional] 
**NetworkInterfaces** | Pointer to [**[]ImageBuildsConfigNetworkInterfacesInner**](ImageBuildsConfigNetworkInterfacesInner.md) |  | [optional] 
**Volumes** | Pointer to [**[]ImageBuildsConfigVolumesInner**](ImageBuildsConfigVolumesInner.md) |  | [optional] 
**StorageControllers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**Config** | Pointer to [**ImageBuildsConfigConfig**](ImageBuildsConfigConfig.md) |  | [optional] 
**Plan** | Pointer to [**ImageBuildsConfigPlan**](ImageBuildsConfigPlan.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ImageBuildsConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


