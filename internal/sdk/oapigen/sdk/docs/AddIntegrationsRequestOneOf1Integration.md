# AddIntegrationsRequestOneOf1Integration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name, a unique identifier for the integration | 
**Type** | **string** | Integration Type Code | 
**Enabled** | Pointer to **bool** | Set &#x60;true&#x60; to enable integration | [optional] 
**Refresh** | Pointer to **bool** | Pass &#x60;false&#x60; to skip refresh.  By default, refresh is done on update, when it is supported by the integration type.  | [optional] [default to true]
**ServiceUrl** | **string** | Ansible Git URL | 
**ServiceUsername** | Pointer to **string** | Git Username | [optional] 
**ServicePassword** | Pointer to **string** | Git Password or Token depending on the Git host | [optional] 
**ServiceToken** | Pointer to **string** | Git Token | [optional] 
**ServiceKey** | Pointer to **int64** | Keypair ID | [optional] 
**Config** | Pointer to [**AddIntegrationsRequestOneOf1IntegrationConfig**](AddIntegrationsRequestOneOf1IntegrationConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddIntegrationsRequestOneOf1Integration{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


