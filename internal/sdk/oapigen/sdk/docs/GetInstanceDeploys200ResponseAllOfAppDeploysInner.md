# GetInstanceDeploys200ResponseAllOfAppDeploysInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to [**GetInstanceDeploys200ResponseAllOfAppDeploysInnerInstance**](GetInstanceDeploys200ResponseAllOfAppDeploysInnerInstance.md) |  | [optional] 
**Deployment** | Pointer to [**GetInstanceDeploys200ResponseAllOfAppDeploysInnerDeployment**](GetInstanceDeploys200ResponseAllOfAppDeploysInnerDeployment.md) |  | [optional] 
**DeploymentVersionId** | Pointer to **int64** |  | [optional] 
**DeploymentVersion** | Pointer to [**GetInstanceDeploys200ResponseAllOfAppDeploysInnerDeploymentVersion**](GetInstanceDeploys200ResponseAllOfAppDeploysInnerDeploymentVersion.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DeployDate** | Pointer to **time.Time** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceDeploys200ResponseAllOfAppDeploysInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


