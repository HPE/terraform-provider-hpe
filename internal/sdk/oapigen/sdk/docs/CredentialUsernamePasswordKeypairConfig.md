# CredentialUsernamePasswordKeypairConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** | Credential Type Code | 
**Name** | **string** | A unique name scoped to your account for the credential | 
**Description** | Pointer to **string** | Optional Description | [optional] 
**Enabled** | Pointer to **bool** | Credential enabled | [optional] [default to true]
**Integration** | Pointer to [**UpdateCredentialsRequestCredentialOneOf7Integration**](UpdateCredentialsRequestCredentialOneOf7Integration.md) |  | [optional] 
**Username** | **string** | Username | 
**Password** | **string** | User password, API Key, or applicable secret | 
**AuthKey** | [**UpdateCredentialsRequestCredentialOneOf7AuthKey**](UpdateCredentialsRequestCredentialOneOf7AuthKey.md) |  | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CredentialUsernamePasswordKeypairConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


