# IntegrationSNOWConfigIntegrationConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ServiceNowCustomCmdbMapping** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ServiceNowCmdbClassMapping** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ServiceNowCMDBBusinessObject** | Pointer to **[]map[string]interface{}** |  | [optional] 
**IgnoreCertErrors** | Pointer to **bool** | Ignore SSL Errors. | [optional] [default to false]
**CmdbMode** | Pointer to **string** |  | [optional] [default to "TABLE"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &IntegrationSNOWConfigIntegrationConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


