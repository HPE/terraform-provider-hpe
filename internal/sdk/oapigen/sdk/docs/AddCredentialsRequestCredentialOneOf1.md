# AddCredentialsRequestCredentialOneOf1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** | Credential Type Code | 
**Name** | **string** | A unique name scoped to your account for the credential | 
**Description** | Pointer to **string** | Optional Description | [optional] 
**Enabled** | Pointer to **bool** | Credential enabled | [optional] [default to true]
**Integration** | Pointer to [**AddCredentialsRequestCredentialOneOf1Integration**](AddCredentialsRequestCredentialOneOf1Integration.md) |  | [optional] 
**Username** | **string** | Client ID | 
**Password** | **string** | Client Secret | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCredentialsRequestCredentialOneOf1{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


