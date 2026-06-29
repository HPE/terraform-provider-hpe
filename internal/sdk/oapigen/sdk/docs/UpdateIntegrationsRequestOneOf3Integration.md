# UpdateIntegrationsRequestOneOf3Integration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name, a unique identifier for the integration | 
**Type** | **string** | Integration Type Code | 
**ServiceMode** | Pointer to **string** | Topology | [optional] [default to "single"]
**ServiceUrl** | **string** | Salt Master | 
**Secondary** | Pointer to **string** | Salt Syndic | [optional] 
**ServicePort** | Pointer to **int32** | SSH Port | [optional] [default to 22]
**ServiceUsername** | **string** | Username | 
**ServicePassword** | Pointer to **string** | Password | [optional] 
**ServiceKey** | Pointer to **string** | Master Key Pair | [optional] 
**AuthKey** | Pointer to **string** | Signing Key | [optional] 
**ServicePath** | Pointer to **string** | Working Directory | [optional] 
**ServiceVersion** | Pointer to **string** | Salt Version | [optional] 
**ServiceWindowsVersion** | Pointer to **string** | Salt Version (Windows) | [optional] 
**RepoUrl** | Pointer to **string** | Repo URL | [optional] 
**ServiceConfig** | Pointer to **string** | Minion Config | [optional] 
**ServiceCommand** | Pointer to **string** | Post Provision Commands | [optional] 
**Config** | Pointer to [**UpdateIntegrationsRequestOneOf3IntegrationConfig**](UpdateIntegrationsRequestOneOf3IntegrationConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIntegrationsRequestOneOf3Integration{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


