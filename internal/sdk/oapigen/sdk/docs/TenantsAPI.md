# \TenantsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddSubtenantIdentitySources**](TenantsAPI.md#AddSubtenantIdentitySources) | **Post** /api/accounts/{accountId}/user-sources | Creates a Subtenant Identity Source
[**AddTenant**](TenantsAPI.md#AddTenant) | **Post** /api/accounts | Create a Tenant
[**AddUserTenant**](TenantsAPI.md#AddUserTenant) | **Post** /api/accounts/{accountId}/users | Create a User For a Tenant
[**CreateTenantSubtenantGroup**](TenantsAPI.md#CreateTenantSubtenantGroup) | **Post** /api/accounts/{accountId}/groups | Create a Group for Subtenant
[**GetTenant**](TenantsAPI.md#GetTenant) | **Get** /api/accounts/{id} | Get tenant
[**GetTenantSubtenantGroup**](TenantsAPI.md#GetTenantSubtenantGroup) | **Get** /api/accounts/{accountId}/groups/{id} | Get a Specific Group for Subtenant
[**ListTenantSubtenantGroups**](TenantsAPI.md#ListTenantSubtenantGroups) | **Get** /api/accounts/{accountId}/groups | Get Subtenant Groups
[**ListTenantSubtenantIdentitySources**](TenantsAPI.md#ListTenantSubtenantIdentitySources) | **Get** /api/accounts/{accountId}/user-sources | Get Subtenant Identity Sources
[**ListTenants**](TenantsAPI.md#ListTenants) | **Get** /api/accounts | List All Tenants
[**ListTenantsAvailableRoles**](TenantsAPI.md#ListTenantsAvailableRoles) | **Get** /api/accounts/available-roles | List available roles for a tenant
[**RemoveTenant**](TenantsAPI.md#RemoveTenant) | **Delete** /api/accounts/{id} | Delete a Specific Tenant
[**RemoveTenantSubtenantGroup**](TenantsAPI.md#RemoveTenantSubtenantGroup) | **Delete** /api/accounts/{accountId}/groups/{id} | Delete a Group for Subtenant
[**UpdateTenant**](TenantsAPI.md#UpdateTenant) | **Put** /api/accounts/{id} | Update tenant
[**UpdateTenantSubtenantGroup**](TenantsAPI.md#UpdateTenantSubtenantGroup) | **Put** /api/accounts/{accountId}/groups/{id} | Updating a Group for Subtenant
[**UpdateTenantSubtenantGroupZones**](TenantsAPI.md#UpdateTenantSubtenantGroupZones) | **Put** /api/accounts/{accountId}/groups/{id}/update-zones | Updating Group Zones for Subtenant



## AddSubtenantIdentitySources

> UpdateIdentitySources200Response AddSubtenantIdentitySources(ctx, accountId).UpdateIdentitySourcesRequest(updateIdentitySourcesRequest).Execute()

Creates a Subtenant Identity Source



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	updateIdentitySourcesRequest := *openapiclient.NewUpdateIdentitySourcesRequest(*openapiclient.NewUpdateIdentitySourcesRequestUserSource("mydomain AD", "Type_example")) // UpdateIdentitySourcesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.AddSubtenantIdentitySources(context.Background(), accountId).UpdateIdentitySourcesRequest(updateIdentitySourcesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.AddSubtenantIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddSubtenantIdentitySources`: UpdateIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.AddSubtenantIdentitySources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddSubtenantIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateIdentitySourcesRequest** | [**UpdateIdentitySourcesRequest**](UpdateIdentitySourcesRequest.md) |  | 

### Return type

[**UpdateIdentitySources200Response**](UpdateIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddTenant

> AddTenant200Response AddTenant(ctx).AddTenantRequest(addTenantRequest).Execute()

Create a Tenant



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
	addTenantRequest := *openapiclient.NewAddTenantRequest(*openapiclient.NewAddTenantRequestAccount("Test Tenant")) // AddTenantRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.AddTenant(context.Background()).AddTenantRequest(addTenantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.AddTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddTenant`: AddTenant200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.AddTenant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addTenantRequest** | [**AddTenantRequest**](AddTenantRequest.md) |  | 

### Return type

[**AddTenant200Response**](AddTenant200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddUserTenant

> AddUserTenant200Response AddUserTenant(ctx, accountId).AddUserTenantRequest(addUserTenantRequest).Execute()

Create a User For a Tenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	addUserTenantRequest := *openapiclient.NewAddUserTenantRequest(*openapiclient.NewAddUserTenantRequestUser("jsmith", "jsmith@email.com", "Passw0rd!", []openapiclient.AddUserTenantRequestUserRolesInner{*openapiclient.NewAddUserTenantRequestUserRolesInner()})) // AddUserTenantRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.AddUserTenant(context.Background(), accountId).AddUserTenantRequest(addUserTenantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.AddUserTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddUserTenant`: AddUserTenant200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.AddUserTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddUserTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **addUserTenantRequest** | [**AddUserTenantRequest**](AddUserTenantRequest.md) |  | 

### Return type

[**AddUserTenant200Response**](AddUserTenant200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateTenantSubtenantGroup

> CreateTenantSubtenantGroup200Response CreateTenantSubtenantGroup(ctx, accountId).CreateTenantSubtenantGroupRequest(createTenantSubtenantGroupRequest).Execute()

Create a Group for Subtenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	createTenantSubtenantGroupRequest := *openapiclient.NewCreateTenantSubtenantGroupRequest(*openapiclient.NewCreateTenantSubtenantGroupRequestGroup("Name_example")) // CreateTenantSubtenantGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.CreateTenantSubtenantGroup(context.Background(), accountId).CreateTenantSubtenantGroupRequest(createTenantSubtenantGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.CreateTenantSubtenantGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateTenantSubtenantGroup`: CreateTenantSubtenantGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.CreateTenantSubtenantGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateTenantSubtenantGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createTenantSubtenantGroupRequest** | [**CreateTenantSubtenantGroupRequest**](CreateTenantSubtenantGroupRequest.md) |  | 

### Return type

[**CreateTenantSubtenantGroup200Response**](CreateTenantSubtenantGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTenant

> GetTenant200Response GetTenant(ctx, id).Execute()

Get tenant



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
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.GetTenant(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.GetTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTenant`: GetTenant200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.GetTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetTenant200Response**](GetTenant200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTenantSubtenantGroup

> GetTenantSubtenantGroup200Response GetTenantSubtenantGroup(ctx, accountId, id).Execute()

Get a Specific Group for Subtenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.GetTenantSubtenantGroup(context.Background(), accountId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.GetTenantSubtenantGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTenantSubtenantGroup`: GetTenantSubtenantGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.GetTenantSubtenantGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTenantSubtenantGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetTenantSubtenantGroup200Response**](GetTenantSubtenantGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTenantSubtenantGroups

> ListTenantSubtenantGroups200Response ListTenantSubtenantGroups(ctx, accountId).Phrase(phrase).Name(name).LastUpdated(lastUpdated).Execute()

Get Subtenant Groups



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	accountId := int64(7) // int64 | The ID of the subtenant account
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.ListTenantSubtenantGroups(context.Background(), accountId).Phrase(phrase).Name(name).LastUpdated(lastUpdated).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.ListTenantSubtenantGroups``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTenantSubtenantGroups`: ListTenantSubtenantGroups200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.ListTenantSubtenantGroups`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 

### Other Parameters

Other parameters are passed through a pointer to a apiListTenantSubtenantGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 

### Return type

[**ListTenantSubtenantGroups200Response**](ListTenantSubtenantGroups200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTenantSubtenantIdentitySources

> ListTenantSubtenantIdentitySources200Response ListTenantSubtenantIdentitySources(ctx, accountId).Max(max).Offset(offset).Direction(direction).Phrase(phrase).Name(name).Execute()

Get Subtenant Identity Sources



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.ListTenantSubtenantIdentitySources(context.Background(), accountId).Max(max).Offset(offset).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.ListTenantSubtenantIdentitySources``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTenantSubtenantIdentitySources`: ListTenantSubtenantIdentitySources200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.ListTenantSubtenantIdentitySources`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 

### Other Parameters

Other parameters are passed through a pointer to a apiListTenantSubtenantIdentitySourcesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListTenantSubtenantIdentitySources200Response**](ListTenantSubtenantIdentitySources200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTenants

> ListTenants200Response ListTenants(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).LastUpdated(lastUpdated).Execute()

List All Tenants



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func main() {
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	lastUpdated := time.Now() // time.Time | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.ListTenants(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).LastUpdated(lastUpdated).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.ListTenants``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTenants`: ListTenants200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.ListTenants`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTenantsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **lastUpdated** | **time.Time** | Date filter, restricts query to only load resources updated more recently than the date specified (ISO 8601) | 

### Return type

[**ListTenants200Response**](ListTenants200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTenantsAvailableRoles

> ListTenantsAvailableRoles200Response ListTenantsAvailableRoles(ctx).Execute()

List available roles for a tenant



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
	resp, r, err := apiClient.TenantsAPI.ListTenantsAvailableRoles(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.ListTenantsAvailableRoles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTenantsAvailableRoles`: ListTenantsAvailableRoles200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.ListTenantsAvailableRoles`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListTenantsAvailableRolesRequest struct via the builder pattern


### Return type

[**ListTenantsAvailableRoles200Response**](ListTenantsAvailableRoles200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveTenant

> RemoveTenant200Response RemoveTenant(ctx, id).RemoveResources(removeResources).Execute()

Delete a Specific Tenant



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
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	removeResources := true // bool | Remove Resources. This will delete all the managed resources in the tenant. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.RemoveTenant(context.Background(), id).RemoveResources(removeResources).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.RemoveTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveTenant`: RemoveTenant200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.RemoveTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **removeResources** | **bool** | Remove Resources. This will delete all the managed resources in the tenant. | [default to false]

### Return type

[**RemoveTenant200Response**](RemoveTenant200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveTenantSubtenantGroup

> RemoveTenantSubtenantGroup200Response RemoveTenantSubtenantGroup(ctx, accountId, id).Execute()

Delete a Group for Subtenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.RemoveTenantSubtenantGroup(context.Background(), accountId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.RemoveTenantSubtenantGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveTenantSubtenantGroup`: RemoveTenantSubtenantGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.RemoveTenantSubtenantGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveTenantSubtenantGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**RemoveTenantSubtenantGroup200Response**](RemoveTenantSubtenantGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTenant

> UpdateTenant200Response UpdateTenant(ctx, id).UpdateTenantRequest(updateTenantRequest).Execute()

Update tenant



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
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateTenantRequest := *openapiclient.NewUpdateTenantRequest(*openapiclient.NewUpdateTenantRequestAccount()) // UpdateTenantRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.UpdateTenant(context.Background(), id).UpdateTenantRequest(updateTenantRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.UpdateTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTenant`: UpdateTenant200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.UpdateTenant`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateTenantRequest** | [**UpdateTenantRequest**](UpdateTenantRequest.md) |  | 

### Return type

[**UpdateTenant200Response**](UpdateTenant200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTenantSubtenantGroup

> UpdateTenantSubtenantGroup200Response UpdateTenantSubtenantGroup(ctx, accountId, id).UpdateTenantSubtenantGroupRequest(updateTenantSubtenantGroupRequest).Execute()

Updating a Group for Subtenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateTenantSubtenantGroupRequest := *openapiclient.NewUpdateTenantSubtenantGroupRequest(*openapiclient.NewUpdateTenantSubtenantGroupRequestGroup()) // UpdateTenantSubtenantGroupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.UpdateTenantSubtenantGroup(context.Background(), accountId, id).UpdateTenantSubtenantGroupRequest(updateTenantSubtenantGroupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.UpdateTenantSubtenantGroup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTenantSubtenantGroup`: UpdateTenantSubtenantGroup200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.UpdateTenantSubtenantGroup`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTenantSubtenantGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateTenantSubtenantGroupRequest** | [**UpdateTenantSubtenantGroupRequest**](UpdateTenantSubtenantGroupRequest.md) |  | 

### Return type

[**UpdateTenantSubtenantGroup200Response**](UpdateTenantSubtenantGroup200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTenantSubtenantGroupZones

> UpdateTenantSubtenantGroupZones200Response UpdateTenantSubtenantGroupZones(ctx, accountId, id).UpdateTenantSubtenantGroupZonesRequest(updateTenantSubtenantGroupZonesRequest).Execute()

Updating Group Zones for Subtenant



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
	accountId := int64(7) // int64 | The ID of the subtenant account
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateTenantSubtenantGroupZonesRequest := *openapiclient.NewUpdateTenantSubtenantGroupZonesRequest(*openapiclient.NewUpdateTenantSubtenantGroupZonesRequestGroup([]openapiclient.UpdateTenantSubtenantGroupZonesRequestGroupZonesInner{*openapiclient.NewUpdateTenantSubtenantGroupZonesRequestGroupZonesInner(int64(123))})) // UpdateTenantSubtenantGroupZonesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantsAPI.UpdateTenantSubtenantGroupZones(context.Background(), accountId, id).UpdateTenantSubtenantGroupZonesRequest(updateTenantSubtenantGroupZonesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantsAPI.UpdateTenantSubtenantGroupZones``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTenantSubtenantGroupZones`: UpdateTenantSubtenantGroupZones200Response
	fmt.Fprintf(os.Stdout, "Response from `TenantsAPI.UpdateTenantSubtenantGroupZones`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **int64** | The ID of the subtenant account | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTenantSubtenantGroupZonesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateTenantSubtenantGroupZonesRequest** | [**UpdateTenantSubtenantGroupZonesRequest**](UpdateTenantSubtenantGroupZonesRequest.md) |  | 

### Return type

[**UpdateTenantSubtenantGroupZones200Response**](UpdateTenantSubtenantGroupZones200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

