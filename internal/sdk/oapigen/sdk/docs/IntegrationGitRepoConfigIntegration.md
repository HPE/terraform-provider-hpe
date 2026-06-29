# IntegrationGitRepoConfigIntegration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name, a unique identifier for the integration | 
**Type** | **string** | Integration Type Code | 
**ServiceUrl** | **string** | Git URL | 
**ServiceUsername** | **string** | Username | 
**ServicePassword** | Pointer to **string** | Password | [optional] 
**ServiceToken** | Pointer to **string** | Access Token | [optional] 
**ServiceKey** | Pointer to **int64** | Key Pair ID | [optional] 
**Config** | Pointer to [**IntegrationGitRepoConfigIntegrationConfig**](IntegrationGitRepoConfigIntegrationConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &IntegrationGitRepoConfigIntegration{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


