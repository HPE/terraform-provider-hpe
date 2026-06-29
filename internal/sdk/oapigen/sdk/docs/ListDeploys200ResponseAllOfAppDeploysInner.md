# ListDeploys200ResponseAllOfAppDeploysInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InstanceId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to [**NullableListDeploys200ResponseAllOfAppDeploysInnerInstance**](ListDeploys200ResponseAllOfAppDeploysInnerInstance.md) |  | [optional] 
**Deployment** | Pointer to [**ListDeploys200ResponseAllOfAppDeploysInnerDeployment**](ListDeploys200ResponseAllOfAppDeploysInnerDeployment.md) |  | [optional] 
**DeploymentVersionId** | Pointer to **int64** |  | [optional] 
**DeploymentVersion** | Pointer to [**ListDeploys200ResponseAllOfAppDeploysInnerDeploymentVersion**](ListDeploys200ResponseAllOfAppDeploysInnerDeploymentVersion.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DeployDate** | Pointer to **time.Time** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListDeploys200ResponseAllOfAppDeploysInner{
    // Set fields directly
}
```

### Instance (Nullable)

Use the Nullable wrapper methods:
- `obj.Instance.IsSet()` — check if set
- `obj.Instance.Get()` — get the inner value (returns pointer)
- `obj.Instance.Set(&val)` — set the value
- `obj.Instance.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


