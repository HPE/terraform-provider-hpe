# AddInstanceDeploy200ResponseAppDeploy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to [**AddInstanceDeploy200ResponseAppDeployInstance**](AddInstanceDeploy200ResponseAppDeployInstance.md) |  | [optional] 
**Deployment** | Pointer to [**AddInstanceDeploy200ResponseAppDeployDeployment**](AddInstanceDeploy200ResponseAppDeployDeployment.md) |  | [optional] 
**DeploymentVersionId** | Pointer to **int64** |  | [optional] 
**DeploymentVersion** | Pointer to [**AddInstanceDeploy200ResponseAppDeployDeploymentVersion**](AddInstanceDeploy200ResponseAppDeployDeploymentVersion.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DeployDate** | Pointer to **time.Time** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddInstanceDeploy200ResponseAppDeploy{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


