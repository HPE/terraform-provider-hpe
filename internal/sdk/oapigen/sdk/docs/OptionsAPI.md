# \OptionsAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetOptionSourceData**](OptionsAPI.md#GetOptionSourceData) | **Get** /api/options/{optionSource} | Get Option Source Data
[**ListCodeRepositories**](OptionsAPI.md#ListCodeRepositories) | **Get** /api/options/codeRepositories | Retrieves a list of Code/GIT Repositories
[**ListOptionAnsibleTowerInventoryOptions**](OptionsAPI.md#ListOptionAnsibleTowerInventoryOptions) | **Get** /api/options/ansibleTowerInventory | Retrieves available Ansible Tower Inventories
[**ListOptionAnsibleTowerJobTemplateOptions**](OptionsAPI.md#ListOptionAnsibleTowerJobTemplateOptions) | **Get** /api/options/ansibleTowerJobTemplate | Retrieves available Ansible Tower Job Templates
[**ListOptionChefServerOptions**](OptionsAPI.md#ListOptionChefServerOptions) | **Get** /api/options/chefServer | Retrieves available Chef Servers
[**ListOptionNetworkOptions**](OptionsAPI.md#ListOptionNetworkOptions) | **Get** /api/options/zoneNetworkOptions | Retrieves network options by zone/cloud
[**ListOptionServiceNowWorkflowsOptions**](OptionsAPI.md#ListOptionServiceNowWorkflowsOptions) | **Get** /api/options/serviceNowWorkflows | Retrieves available ServiceNow workflows
[**ListOptionValues**](OptionsAPI.md#ListOptionValues) | **Get** /api/options/list | Retrieves input option values
[**ListOptionZoneTypesOptions**](OptionsAPI.md#ListOptionZoneTypesOptions) | **Get** /api/options/zoneTypes | Retrieves enabled zones/clouds



## GetOptionSourceData

> GetOptionSourceData200Response GetOptionSourceData(ctx, optionSource).Execute()

Get Option Source Data



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
	optionSource := "keyPairs" // string | `optionSource` to be listed

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.GetOptionSourceData(context.Background(), optionSource).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.GetOptionSourceData``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOptionSourceData`: GetOptionSourceData200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.GetOptionSourceData`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**optionSource** | **string** | &#x60;optionSource&#x60; to be listed | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOptionSourceDataRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetOptionSourceData200Response**](GetOptionSourceData200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCodeRepositories

> ListCodeRepositories200Response ListCodeRepositories(ctx).IntegrationId(integrationId).Execute()

Retrieves a list of Code/GIT Repositories



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
	integrationId := int64(33) // int64 | Filter by an integration Id. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListCodeRepositories(context.Background()).IntegrationId(integrationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListCodeRepositories``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCodeRepositories`: ListCodeRepositories200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListCodeRepositories`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCodeRepositoriesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **integrationId** | **int64** | Filter by an integration Id. | 

### Return type

[**ListCodeRepositories200Response**](ListCodeRepositories200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionAnsibleTowerInventoryOptions

> ListOptionAnsibleTowerInventoryOptions200Response ListOptionAnsibleTowerInventoryOptions(ctx).ZoneId(zoneId).SiteId(siteId).TaskId(taskId).AccountId(accountId).AnsibleTowerIntegrationId(ansibleTowerIntegrationId).Execute()

Retrieves available Ansible Tower Inventories



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
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	siteId := int64(7) // int64 | The Group ID (Site ID) for Filtering (optional)
	taskId := int64(5) // int64 | The Task ID for filtering (optional)
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)
	ansibleTowerIntegrationId := int64(33) // int64 | Filter by an integration Id of an Ansible Tower Integration. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionAnsibleTowerInventoryOptions(context.Background()).ZoneId(zoneId).SiteId(siteId).TaskId(taskId).AccountId(accountId).AnsibleTowerIntegrationId(ansibleTowerIntegrationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionAnsibleTowerInventoryOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionAnsibleTowerInventoryOptions`: ListOptionAnsibleTowerInventoryOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionAnsibleTowerInventoryOptions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionAnsibleTowerInventoryOptionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **siteId** | **int64** | The Group ID (Site ID) for Filtering | 
 **taskId** | **int64** | The Task ID for filtering | 
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 
 **ansibleTowerIntegrationId** | **int64** | Filter by an integration Id of an Ansible Tower Integration. | 

### Return type

[**ListOptionAnsibleTowerInventoryOptions200Response**](ListOptionAnsibleTowerInventoryOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionAnsibleTowerJobTemplateOptions

> ListOptionAnsibleTowerJobTemplateOptions200Response ListOptionAnsibleTowerJobTemplateOptions(ctx).ZoneId(zoneId).SiteId(siteId).TaskId(taskId).AccountId(accountId).AnsibleTowerIntegrationId(ansibleTowerIntegrationId).Execute()

Retrieves available Ansible Tower Job Templates



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
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	siteId := int64(7) // int64 | The Group ID (Site ID) for Filtering (optional)
	taskId := int64(5) // int64 | The Task ID for filtering (optional)
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)
	ansibleTowerIntegrationId := int64(33) // int64 | Filter by an integration Id of an Ansible Tower Integration. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionAnsibleTowerJobTemplateOptions(context.Background()).ZoneId(zoneId).SiteId(siteId).TaskId(taskId).AccountId(accountId).AnsibleTowerIntegrationId(ansibleTowerIntegrationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionAnsibleTowerJobTemplateOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionAnsibleTowerJobTemplateOptions`: ListOptionAnsibleTowerJobTemplateOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionAnsibleTowerJobTemplateOptions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionAnsibleTowerJobTemplateOptionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **siteId** | **int64** | The Group ID (Site ID) for Filtering | 
 **taskId** | **int64** | The Task ID for filtering | 
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 
 **ansibleTowerIntegrationId** | **int64** | Filter by an integration Id of an Ansible Tower Integration. | 

### Return type

[**ListOptionAnsibleTowerJobTemplateOptions200Response**](ListOptionAnsibleTowerJobTemplateOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionChefServerOptions

> ListOptionChefServerOptions200Response ListOptionChefServerOptions(ctx).AccountId(accountId).Execute()

Retrieves available Chef Servers



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
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionChefServerOptions(context.Background()).AccountId(accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionChefServerOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionChefServerOptions`: ListOptionChefServerOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionChefServerOptions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionChefServerOptionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 

### Return type

[**ListOptionChefServerOptions200Response**](ListOptionChefServerOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionNetworkOptions

> ListOptionNetworkOptions200Response ListOptionNetworkOptions(ctx).ZoneId(zoneId).ProvisionTypeId(provisionTypeId).Execute()

Retrieves network options by zone/cloud



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
	zoneId := int64(3) // int64 | The Cloud ID (Zone ID) for Filtering (optional)
	provisionTypeId := int64(22) // int64 | Provision type filter, restricts query to only load service plans of specified provision type (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionNetworkOptions(context.Background()).ZoneId(zoneId).ProvisionTypeId(provisionTypeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionNetworkOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionNetworkOptions`: ListOptionNetworkOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionNetworkOptions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionNetworkOptionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zoneId** | **int64** | The Cloud ID (Zone ID) for Filtering | 
 **provisionTypeId** | **int64** | Provision type filter, restricts query to only load service plans of specified provision type | 

### Return type

[**ListOptionNetworkOptions200Response**](ListOptionNetworkOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionServiceNowWorkflowsOptions

> ListOptionServiceNowWorkflowsOptions200Response ListOptionServiceNowWorkflowsOptions(ctx).AccountId(accountId).Config(config).AccountIntegrationId(accountIntegrationId).Execute()

Retrieves available ServiceNow workflows



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
	accountId := int64(3) // int64 | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. (optional)
	config := map[string]interface{}{ ... } // map[string]interface{} | Input parameters are required if the input is dependent on them.  Fields must be prefixed with `config.` (optional)
	accountIntegrationId := int64(4) // int64 | Filter by Account Integration ID. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionServiceNowWorkflowsOptions(context.Background()).AccountId(accountId).Config(config).AccountIntegrationId(accountIntegrationId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionServiceNowWorkflowsOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionServiceNowWorkflowsOptions`: ListOptionServiceNowWorkflowsOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionServiceNowWorkflowsOptions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionServiceNowWorkflowsOptionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **int64** | Filter by Tenant ID. This is only available to master tenant users with permission to manage tenants and users. | 
 **config** | [**map[string]interface{}**](map[string]interface{}.md) | Input parameters are required if the input is dependent on them.  Fields must be prefixed with &#x60;config.&#x60; | 
 **accountIntegrationId** | **int64** | Filter by Account Integration ID. | 

### Return type

[**ListOptionServiceNowWorkflowsOptions200Response**](ListOptionServiceNowWorkflowsOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionValues

> ListOptionValues200Response ListOptionValues(ctx).OptionTypeId(optionTypeId).Config(config).Execute()

Retrieves input option values



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
	optionTypeId := int64(789) // int64 | Input or Option Type ID
	config := map[string]interface{}{ ... } // map[string]interface{} | Input parameters are required if the input is dependent on them.  Fields must be prefixed with `config.` (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OptionsAPI.ListOptionValues(context.Background()).OptionTypeId(optionTypeId).Config(config).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionValues``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionValues`: ListOptionValues200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionValues`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOptionValuesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **optionTypeId** | **int64** | Input or Option Type ID | 
 **config** | [**map[string]interface{}**](map[string]interface{}.md) | Input parameters are required if the input is dependent on them.  Fields must be prefixed with &#x60;config.&#x60; | 

### Return type

[**ListOptionValues200Response**](ListOptionValues200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOptionZoneTypesOptions

> ListOptionZoneTypesOptions200Response ListOptionZoneTypesOptions(ctx).Execute()

Retrieves enabled zones/clouds



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
	resp, r, err := apiClient.OptionsAPI.ListOptionZoneTypesOptions(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OptionsAPI.ListOptionZoneTypesOptions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOptionZoneTypesOptions`: ListOptionZoneTypesOptions200Response
	fmt.Fprintf(os.Stdout, "Response from `OptionsAPI.ListOptionZoneTypesOptions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOptionZoneTypesOptionsRequest struct via the builder pattern


### Return type

[**ListOptionZoneTypesOptions200Response**](ListOptionZoneTypesOptions200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

