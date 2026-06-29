# TokenAuthorizationCodeRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | **string** | Client ID | 
**ClientSecret** | **string** | Client Secret | 
**GrantType** | **string** | OAuth Grant Type, use authorization_code. | [default to "authorization_code"]
**AuthorizationCode** | **string** | Authorization code must be obtained with a valid request to &#x60;/oauth/authorize&#x60;. This code is used to request an access token in the OAuth 2.0 Authorization Code Flow. | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &TokenAuthorizationCodeRequest{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


