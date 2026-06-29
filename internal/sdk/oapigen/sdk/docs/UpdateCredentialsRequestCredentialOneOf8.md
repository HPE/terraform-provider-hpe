# UpdateCredentialsRequestCredentialOneOf8

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** | Credential Type Code | 
**Name** | **string** | A unique name scoped to your account for the credential | 
**Description** | Pointer to **string** | Optional Description | [optional] 
**Enabled** | Pointer to **bool** | Credential enabled | [optional] [default to true]
**Integration** | Pointer to [**UpdateCredentialsRequestCredentialOneOf8Integration**](UpdateCredentialsRequestCredentialOneOf8Integration.md) |  | [optional] 
**Username** | Pointer to **string** | Username | [optional] 
**Password** | Pointer to **string** | Password | [optional] 
**Config** | [**UpdateCredentialsRequestCredentialOneOf8Config**](UpdateCredentialsRequestCredentialOneOf8Config.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateCredentialsRequestCredentialOneOf8{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


