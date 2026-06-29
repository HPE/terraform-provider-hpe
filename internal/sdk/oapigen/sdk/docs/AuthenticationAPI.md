# \AuthenticationAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ForgotPassword**](AuthenticationAPI.md#ForgotPassword) | **Post** /api/forgot/send-email | Request a reset password email
[**GetAccessToken**](AuthenticationAPI.md#GetAccessToken) | **Post** /oauth/token | Get Access Token
[**ResetPassword**](AuthenticationAPI.md#ResetPassword) | **Post** /api/forgot/reset-password | Reset user password
[**Whoami**](AuthenticationAPI.md#Whoami) | **Get** /api/whoami | Whoami



## ForgotPassword

> ForgotPassword200Response ForgotPassword(ctx).ForgotPasswordRequest(forgotPasswordRequest).Execute()

Request a reset password email



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	forgotPasswordRequest := *openapiclient.NewForgotPasswordRequest("Username_example") // ForgotPasswordRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthenticationAPI.ForgotPassword(context.Background()).ForgotPasswordRequest(forgotPasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.ForgotPassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ForgotPassword`: ForgotPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.ForgotPassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiForgotPasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **forgotPasswordRequest** | [**ForgotPasswordRequest**](ForgotPasswordRequest.md) |  | 

### Return type

[**ForgotPassword200Response**](ForgotPassword200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAccessToken

> GetAccessToken200Response GetAccessToken(ctx).ClientId(clientId).GrantType(grantType).Scope(scope).Username(username).Password(password).RefreshToken(refreshToken).ClientSecret(clientSecret).AuthorizationCode(authorizationCode).Execute()

Get Access Token



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	clientId := "clientId_example" // string | Client ID (optional)
	grantType := "grantType_example" // string | OAuth Grant Type, use authorization_code. (optional) (default to "authorization_code")
	scope := "scope_example" // string | OAuth token scope, use `write`. (optional)
	username := "username_example" // string | Username Sub-tenant users must format their username as `subdomain\\\\username` with a prefix that is the tenant subdomain or id by default.  (optional)
	password := "password_example" // string | Password (optional)
	refreshToken := "refreshToken_example" // string | Refresh Token (optional)
	clientSecret := "clientSecret_example" // string | Client Secret (optional)
	authorizationCode := "authorizationCode_example" // string | Authorization code must be obtained with a valid request to `/oauth/authorize`. This code is used to request an access token in the OAuth 2.0 Authorization Code Flow. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthenticationAPI.GetAccessToken(context.Background()).ClientId(clientId).GrantType(grantType).Scope(scope).Username(username).Password(password).RefreshToken(refreshToken).ClientSecret(clientSecret).AuthorizationCode(authorizationCode).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.GetAccessToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAccessToken`: GetAccessToken200Response
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.GetAccessToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAccessTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **clientId** | **string** | Client ID | 
 **grantType** | **string** | OAuth Grant Type, use authorization_code. | [default to &quot;authorization_code&quot;]
 **scope** | **string** | OAuth token scope, use &#x60;write&#x60;. | 
 **username** | **string** | Username Sub-tenant users must format their username as &#x60;subdomain\\\\username&#x60; with a prefix that is the tenant subdomain or id by default.  | 
 **password** | **string** | Password | 
 **refreshToken** | **string** | Refresh Token | 
 **clientSecret** | **string** | Client Secret | 
 **authorizationCode** | **string** | Authorization code must be obtained with a valid request to &#x60;/oauth/authorize&#x60;. This code is used to request an access token in the OAuth 2.0 Authorization Code Flow. | 

### Return type

[**GetAccessToken200Response**](GetAccessToken200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ResetPassword

> ResetPassword200Response ResetPassword(ctx).ResetPasswordRequest(resetPasswordRequest).Execute()

Reset user password



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	resetPasswordRequest := *openapiclient.NewResetPasswordRequest("Token_example", "Password_example") // ResetPasswordRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthenticationAPI.ResetPassword(context.Background()).ResetPasswordRequest(resetPasswordRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.ResetPassword``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ResetPassword`: ResetPassword200Response
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.ResetPassword`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiResetPasswordRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resetPasswordRequest** | [**ResetPasswordRequest**](ResetPasswordRequest.md) |  | 

### Return type

[**ResetPassword200Response**](ResetPassword200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Whoami

> Whoami200Response Whoami(ctx).Execute()

Whoami



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AuthenticationAPI.Whoami(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticationAPI.Whoami``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Whoami`: Whoami200Response
	fmt.Fprintf(os.Stdout, "Response from `AuthenticationAPI.Whoami`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiWhoamiRequest struct via the builder pattern


### Return type

[**Whoami200Response**](Whoami200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

