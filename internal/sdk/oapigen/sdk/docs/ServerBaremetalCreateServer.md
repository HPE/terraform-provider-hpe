# ServerBaremetalCreateServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cloud** | Pointer to [**ServerBaremetalCreateServerCloud**](ServerBaremetalCreateServerCloud.md) |  | [optional] 
**ComputeServerType** | Pointer to [**ServerBaremetalCreateServerComputeServerType**](ServerBaremetalCreateServerComputeServerType.md) |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Group** | Pointer to [**ServerBaremetalCreateServerGroup**](ServerBaremetalCreateServerGroup.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Config** | Pointer to [**ServerBaremetalCreateServerConfig**](ServerBaremetalCreateServerConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ServerBaremetalCreateServer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


