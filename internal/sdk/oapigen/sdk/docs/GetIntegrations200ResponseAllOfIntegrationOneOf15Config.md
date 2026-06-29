# GetIntegrations200ResponseAllOfIntegrationOneOf15Config

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IncidentAccess** | Pointer to **bool** |  | [optional] 
**RequestAccess** | Pointer to **bool** |  | [optional] 
**ServiceNowCMDBBusinessObject** | Pointer to **string** |  | [optional] 
**ServiceNowCustomCmdbMapping** | Pointer to **string** |  | [optional] 
**ServiceNowCmdbClassMapping** | Pointer to [**[]GetIntegrations200ResponseAllOfIntegrationOneOf15ConfigServiceNowCmdbClassMappingInner**](GetIntegrations200ResponseAllOfIntegrationOneOf15ConfigServiceNowCmdbClassMappingInner.md) |  | [optional] 
**WebServiceImportUrl** | Pointer to **string** |  | [optional] 
**WebServiceImportSysId** | Pointer to **string** |  | [optional] 
**WebServiceOperationUrl** | Pointer to **string** |  | [optional] 
**CmdbMode** | Pointer to **string** |  | [optional] [default to "TABLE"]
**PreparedForSync** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetIntegrations200ResponseAllOfIntegrationOneOf15Config{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


