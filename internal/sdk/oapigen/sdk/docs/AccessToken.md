# AccessToken

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | Pointer to **string** | Access token that is used as a bearer token in the Authorization header of API requests to protected resources | [optional] 
**RefreshToken** | Pointer to **string** | Refresh token that allows refreshing access token without re-authenticating | [optional] 
**ExpiresIn** | Pointer to **float32** | Seconds until token expires | [optional] 
**TokenType** | Pointer to **string** | Token type granted | [optional] 
**Scope** | Pointer to **string** | Scope granted | [optional] 
**IdToken** | Pointer to **string** | ID token that is only returned when using the &#x60;openid&#x60; scope. The ID token is a JSON Web Token (JWT) that contains claims about the authenticated user and can be decoded and verified using the tenant&#39;s public key. The claims in the ID token can be used to obtain user information such as the user&#39;s name, email, and roles without needing to make additional API calls. The ID token is intended for use in OIDC-compliant applications that require user authentication and authorization based on OIDC claims. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AccessToken{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


