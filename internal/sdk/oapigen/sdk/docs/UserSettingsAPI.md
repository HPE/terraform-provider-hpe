# \UserSettingsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteUserSettingsAccessToken**](UserSettingsAPI.md#DeleteUserSettingsAccessToken) | **Put** /api/user-settings/clear-access-token | Revoke API Access Token
[**DeleteUserSettingsAvatar**](UserSettingsAPI.md#DeleteUserSettingsAvatar) | **Delete** /api/user-settings/avatar | Delete Avatar
[**DeleteUserSettingsDesktopBackground**](UserSettingsAPI.md#DeleteUserSettingsDesktopBackground) | **Delete** /api/user-settings/desktop-background | Delete Desktop Background
[**GetUserSettingsApiClients**](UserSettingsAPI.md#GetUserSettingsApiClients) | **Get** /api/user-settings/api-clients | Get Available API Clients
[**ListUserSettings**](UserSettingsAPI.md#ListUserSettings) | **Get** /api/user-settings | User Settings
[**UpdateUserSettings**](UserSettingsAPI.md#UpdateUserSettings) | **Put** /api/user-settings | Update User Settings
[**UpdateUserSettingsAccessToken**](UserSettingsAPI.md#UpdateUserSettingsAccessToken) | **Put** /api/user-settings/regenerate-access-token | Regenerate API Access Token
[**UpdateUserSettingsAvatar**](UserSettingsAPI.md#UpdateUserSettingsAvatar) | **Post** /api/user-settings/avatar | Update Avatar
[**UpdateUserSettingsDesktopBackground**](UserSettingsAPI.md#UpdateUserSettingsDesktopBackground) | **Post** /api/user-settings/desktop-background | Update Desktop Background



## DeleteUserSettingsAccessToken

> DeleteUserSettingsAccessToken200Response DeleteUserSettingsAccessToken(ctx).UserId(userId).ClientId(clientId).Execute()

Revoke API Access Token



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)
	clientId := "morph-cli" // string | Client ID.  See `Get Available API Clients` for a list of valid `clientId` values. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.DeleteUserSettingsAccessToken(context.Background()).UserId(userId).ClientId(clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.DeleteUserSettingsAccessToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteUserSettingsAccessToken`: DeleteUserSettingsAccessToken200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.DeleteUserSettingsAccessToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserSettingsAccessTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 
 **clientId** | **string** | Client ID.  See &#x60;Get Available API Clients&#x60; for a list of valid &#x60;clientId&#x60; values. | 

### Return type

[**DeleteUserSettingsAccessToken200Response**](DeleteUserSettingsAccessToken200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUserSettingsAvatar

> DeleteUserSettingsAvatar200Response DeleteUserSettingsAvatar(ctx).UserId(userId).Execute()

Delete Avatar



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.DeleteUserSettingsAvatar(context.Background()).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.DeleteUserSettingsAvatar``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteUserSettingsAvatar`: DeleteUserSettingsAvatar200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.DeleteUserSettingsAvatar`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserSettingsAvatarRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 

### Return type

[**DeleteUserSettingsAvatar200Response**](DeleteUserSettingsAvatar200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUserSettingsDesktopBackground

> DeleteUserSettingsDesktopBackground200Response DeleteUserSettingsDesktopBackground(ctx).UserId(userId).Execute()

Delete Desktop Background



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.DeleteUserSettingsDesktopBackground(context.Background()).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.DeleteUserSettingsDesktopBackground``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteUserSettingsDesktopBackground`: DeleteUserSettingsDesktopBackground200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.DeleteUserSettingsDesktopBackground`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserSettingsDesktopBackgroundRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 

### Return type

[**DeleteUserSettingsDesktopBackground200Response**](DeleteUserSettingsDesktopBackground200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUserSettingsApiClients

> GetUserSettingsApiClients200Response GetUserSettingsApiClients(ctx).UserId(userId).Execute()

Get Available API Clients



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.GetUserSettingsApiClients(context.Background()).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.GetUserSettingsApiClients``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetUserSettingsApiClients`: GetUserSettingsApiClients200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.GetUserSettingsApiClients`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetUserSettingsApiClientsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 

### Return type

[**GetUserSettingsApiClients200Response**](GetUserSettingsApiClients200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUserSettings

> ListUserSettings200Response ListUserSettings(ctx).UserId(userId).Execute()

User Settings



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.ListUserSettings(context.Background()).UserId(userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.ListUserSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListUserSettings`: ListUserSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.ListUserSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListUserSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 

### Return type

[**ListUserSettings200Response**](ListUserSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUserSettings

> UpdateUserSettings200Response UpdateUserSettings(ctx).UserId(userId).UpdateUserSettingsRequest(updateUserSettingsRequest).Execute()

Update User Settings



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)
	updateUserSettingsRequest := *openapiclient.NewUpdateUserSettingsRequest() // UpdateUserSettingsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.UpdateUserSettings(context.Background()).UserId(userId).UpdateUserSettingsRequest(updateUserSettingsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.UpdateUserSettings``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateUserSettings`: UpdateUserSettings200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.UpdateUserSettings`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserSettingsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 
 **updateUserSettingsRequest** | [**UpdateUserSettingsRequest**](UpdateUserSettingsRequest.md) |  | 

### Return type

[**UpdateUserSettings200Response**](UpdateUserSettings200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json, multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUserSettingsAccessToken

> UpdateUserSettingsAccessToken200Response UpdateUserSettingsAccessToken(ctx).Id(id).UserId(userId).ClientId(clientId).Execute()

Regenerate API Access Token



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
	id := int64(789) // int64 | ID of Access Token to be regenerated. If specified, the existing token will be revoked and a new token will be generated for the same client. (optional)
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)
	clientId := "morph-cli" // string | Client ID.  See `Get Available API Clients` for a list of valid `clientId` values. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.UpdateUserSettingsAccessToken(context.Background()).Id(id).UserId(userId).ClientId(clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.UpdateUserSettingsAccessToken``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateUserSettingsAccessToken`: UpdateUserSettingsAccessToken200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.UpdateUserSettingsAccessToken`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserSettingsAccessTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int64** | ID of Access Token to be regenerated. If specified, the existing token will be revoked and a new token will be generated for the same client. | 
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 
 **clientId** | **string** | Client ID.  See &#x60;Get Available API Clients&#x60; for a list of valid &#x60;clientId&#x60; values. | 

### Return type

[**UpdateUserSettingsAccessToken200Response**](UpdateUserSettingsAccessToken200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUserSettingsAvatar

> UpdateUserSettingsAvatar200Response UpdateUserSettingsAvatar(ctx).UserId(userId).UserAvatar(userAvatar).Execute()

Update Avatar



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)
	userAvatar := "userAvatar_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.UpdateUserSettingsAvatar(context.Background()).UserId(userId).UserAvatar(userAvatar).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.UpdateUserSettingsAvatar``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateUserSettingsAvatar`: UpdateUserSettingsAvatar200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.UpdateUserSettingsAvatar`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserSettingsAvatarRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 
 **userAvatar** | **string** |  | 

### Return type

[**UpdateUserSettingsAvatar200Response**](UpdateUserSettingsAvatar200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateUserSettingsDesktopBackground

> UpdateUserSettingsDesktopBackground200Response UpdateUserSettingsDesktopBackground(ctx).UserId(userId).UserDesktopBackground(userDesktopBackground).Execute()

Update Desktop Background



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
	userId := int64(789) // int64 | ID of User (Only available to the master tenant.)  Default is the current user (optional)
	userDesktopBackground := "userDesktopBackground_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.UserSettingsAPI.UpdateUserSettingsDesktopBackground(context.Background()).UserId(userId).UserDesktopBackground(userDesktopBackground).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `UserSettingsAPI.UpdateUserSettingsDesktopBackground``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateUserSettingsDesktopBackground`: UpdateUserSettingsDesktopBackground200Response
	fmt.Fprintf(os.Stdout, "Response from `UserSettingsAPI.UpdateUserSettingsDesktopBackground`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateUserSettingsDesktopBackgroundRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int64** | ID of User (Only available to the master tenant.)  Default is the current user | 
 **userDesktopBackground** | **string** |  | 

### Return type

[**UpdateUserSettingsDesktopBackground200Response**](UpdateUserSettingsDesktopBackground200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

