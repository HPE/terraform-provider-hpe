# UpdateDeploy200ResponseAppDeploy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to [**UpdateDeploy200ResponseAppDeployInstance**](UpdateDeploy200ResponseAppDeployInstance.md) |  | [optional] 
**Deployment** | Pointer to [**UpdateDeploy200ResponseAppDeployDeployment**](UpdateDeploy200ResponseAppDeployDeployment.md) |  | [optional] 
**DeploymentVersionId** | Pointer to **int64** |  | [optional] 
**DeploymentVersion** | Pointer to [**UpdateDeploy200ResponseAppDeployDeploymentVersion**](UpdateDeploy200ResponseAppDeployDeploymentVersion.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DeployDate** | Pointer to **time.Time** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateDeploy200ResponseAppDeploy{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


