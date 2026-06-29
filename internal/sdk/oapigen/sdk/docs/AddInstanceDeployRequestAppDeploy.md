# AddInstanceDeployRequestAppDeploy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeploymentId** | Pointer to **int64** | Deployment ID. | [optional] 
**Version** | Pointer to **string** | Deployment Version number identifier (userVersion). Can be passed along with deploymentId to identify the version | [optional] 
**VersionId** | Pointer to **int64** | Deployment Version ID. This can be passed instead of deploymentId and version. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Map of configuration properties that vary by instance type. | [optional] 
**StageOnly** | Pointer to **bool** | Stage Only, do not run the deploy right away and instead set status to staged so it can be deployed later on. | [optional] [default to false]

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddInstanceDeployRequestAppDeploy{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


