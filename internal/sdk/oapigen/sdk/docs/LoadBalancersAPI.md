# \LoadBalancersAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateLoadBalancer**](LoadBalancersAPI.md#CreateLoadBalancer) | **Post** /api/load-balancers | Create a Load Balancer
[**CreateLoadBalancerMonitor**](LoadBalancersAPI.md#CreateLoadBalancerMonitor) | **Post** /api/load-balancers/{loadBalancerId}/monitors | Create a Load Balancer Monitor
[**CreateLoadBalancerPool**](LoadBalancersAPI.md#CreateLoadBalancerPool) | **Post** /api/load-balancers/{loadBalancerId}/pools | Create a Load Balancer Pool
[**CreateLoadBalancerPoolNode**](LoadBalancersAPI.md#CreateLoadBalancerPoolNode) | **Post** /api/load-balancer-pools/{loadBalancerPoolId}/nodes | Create a Load Balancer Pool Node
[**CreateLoadBalancerProfile**](LoadBalancersAPI.md#CreateLoadBalancerProfile) | **Post** /api/load-balancers/{loadBalancerId}/profiles | Create a Load Balancer Profile
[**CreateLoadBalancerVirtualServer**](LoadBalancersAPI.md#CreateLoadBalancerVirtualServer) | **Post** /api/load-balancers/{loadBalancerId}/virtual-servers | Create a Load Balancer Virtual Server
[**DeleteLoadBalancer**](LoadBalancersAPI.md#DeleteLoadBalancer) | **Delete** /api/load-balancers/{loadBalancerId} | Delete a Load Balancer
[**DeleteLoadBalancerMonitor**](LoadBalancersAPI.md#DeleteLoadBalancerMonitor) | **Delete** /api/load-balancers/{loadBalancerId}/monitors/{id} | Delete a Load Balancer Monitor
[**DeleteLoadBalancerPool**](LoadBalancersAPI.md#DeleteLoadBalancerPool) | **Delete** /api/load-balancers/{loadBalancerId}/pools/{id} | Delete a Load Balancer Pool
[**DeleteLoadBalancerPoolNode**](LoadBalancersAPI.md#DeleteLoadBalancerPoolNode) | **Delete** /api/load-balancer-pools/{loadBalancerPoolId}/nodes/{id} | Delete a Load Balancer Pool Node
[**DeleteLoadBalancerProfile**](LoadBalancersAPI.md#DeleteLoadBalancerProfile) | **Delete** /api/load-balancers/{loadBalancerId}/profiles/{id} | Delete a Load Balancer Profile
[**DeleteLoadBalancerVirtualServer**](LoadBalancersAPI.md#DeleteLoadBalancerVirtualServer) | **Delete** /api/load-balancers/{loadBalancerId}/virtual-servers/{id} | Delete a Load Balancer Virtual Server
[**GetLoadBalancer**](LoadBalancersAPI.md#GetLoadBalancer) | **Get** /api/load-balancers/{loadBalancerId} | Get a Specific Load Balancer
[**GetLoadBalancerMonitor**](LoadBalancersAPI.md#GetLoadBalancerMonitor) | **Get** /api/load-balancers/{loadBalancerId}/monitors/{id} | Get a Specific Load Balancer Monitor
[**GetLoadBalancerPool**](LoadBalancersAPI.md#GetLoadBalancerPool) | **Get** /api/load-balancers/{loadBalancerId}/pools/{id} | Get a Specific Load Balancer Pool
[**GetLoadBalancerPoolNode**](LoadBalancersAPI.md#GetLoadBalancerPoolNode) | **Get** /api/load-balancer-pools/{loadBalancerPoolId}/nodes/{id} | Get a Specific Load Balancer Pool Node
[**GetLoadBalancerProfile**](LoadBalancersAPI.md#GetLoadBalancerProfile) | **Get** /api/load-balancers/{loadBalancerId}/profiles/{id} | Get a Specific Load Balancer Profile
[**GetLoadBalancerType**](LoadBalancersAPI.md#GetLoadBalancerType) | **Get** /api/load-balancer-types/{id} | Get a Specific Load Balancer Type
[**GetLoadBalancerVirtualServer**](LoadBalancersAPI.md#GetLoadBalancerVirtualServer) | **Get** /api/load-balancers/{loadBalancerId}/virtual-servers/{id} | Get a Specific Load Balancer Virtual Server
[**ListLoadBalancerMonitors**](LoadBalancersAPI.md#ListLoadBalancerMonitors) | **Get** /api/load-balancers/{loadBalancerId}/monitors | Get All Load Balancer Monitors For Load Balancer
[**ListLoadBalancerPoolNodes**](LoadBalancersAPI.md#ListLoadBalancerPoolNodes) | **Get** /api/load-balancer-pools/{loadBalancerPoolId}/nodes | Get All Load Balancer Pool Nodes For Load Balancer Pool
[**ListLoadBalancerPools**](LoadBalancersAPI.md#ListLoadBalancerPools) | **Get** /api/load-balancers/{loadBalancerId}/pools | Get All Load Balancer Pools For Load Balancer
[**ListLoadBalancerProfiles**](LoadBalancersAPI.md#ListLoadBalancerProfiles) | **Get** /api/load-balancers/{loadBalancerId}/profiles | Get All Load Balancer Profiles For Load Balancer
[**ListLoadBalancerTypes**](LoadBalancersAPI.md#ListLoadBalancerTypes) | **Get** /api/load-balancer-types | Get All Load Balancer Types
[**ListLoadBalancerVirtualServers**](LoadBalancersAPI.md#ListLoadBalancerVirtualServers) | **Get** /api/load-balancers/{loadBalancerId}/virtual-servers | Get All Load Balancer Virtual Servers For Load Balancer
[**ListLoadBalancers**](LoadBalancersAPI.md#ListLoadBalancers) | **Get** /api/load-balancers | Get All Load Balancers
[**RefreshLoadBalancer**](LoadBalancersAPI.md#RefreshLoadBalancer) | **Put** /api/load-balancers/{loadBalancerId}/refresh | Refresh a Load Balancer
[**UpdateLoadBalancer**](LoadBalancersAPI.md#UpdateLoadBalancer) | **Put** /api/load-balancers/{loadBalancerId} | Update a Load Balancer
[**UpdateLoadBalancerMonitor**](LoadBalancersAPI.md#UpdateLoadBalancerMonitor) | **Put** /api/load-balancers/{loadBalancerId}/monitors/{id} | Update a Load Balancer Monitor
[**UpdateLoadBalancerPool**](LoadBalancersAPI.md#UpdateLoadBalancerPool) | **Put** /api/load-balancers/{loadBalancerId}/pools/{id} | Update a Load Balancer Pool
[**UpdateLoadBalancerPoolNode**](LoadBalancersAPI.md#UpdateLoadBalancerPoolNode) | **Put** /api/load-balancer-pools/{loadBalancerPoolId}/nodes/{id} | Update a Load Balancer Pool Node
[**UpdateLoadBalancerProfile**](LoadBalancersAPI.md#UpdateLoadBalancerProfile) | **Put** /api/load-balancers/{loadBalancerId}/profiles/{id} | Update a Load Balancer Profile
[**UpdateLoadBalancerVirtualServer**](LoadBalancersAPI.md#UpdateLoadBalancerVirtualServer) | **Put** /api/load-balancers/{loadBalancerId}/virtual-servers/{id} | Update a Load Balancer Virtual Server



## CreateLoadBalancer

> CreateLoadBalancer200Response CreateLoadBalancer(ctx).CreateLoadBalancerRequest(createLoadBalancerRequest).Execute()

Create a Load Balancer



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
	createLoadBalancerRequest := *openapiclient.NewCreateLoadBalancerRequest() // CreateLoadBalancerRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancer(context.Background()).CreateLoadBalancerRequest(createLoadBalancerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancer`: CreateLoadBalancer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createLoadBalancerRequest** | [**CreateLoadBalancerRequest**](CreateLoadBalancerRequest.md) |  | 

### Return type

[**CreateLoadBalancer200Response**](CreateLoadBalancer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLoadBalancerMonitor

> CreateLoadBalancerMonitor200Response CreateLoadBalancerMonitor(ctx, loadBalancerId).CreateLoadBalancerMonitorRequest(createLoadBalancerMonitorRequest).Execute()

Create a Load Balancer Monitor



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	createLoadBalancerMonitorRequest := *openapiclient.NewCreateLoadBalancerMonitorRequest() // CreateLoadBalancerMonitorRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancerMonitor(context.Background(), loadBalancerId).CreateLoadBalancerMonitorRequest(createLoadBalancerMonitorRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancerMonitor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancerMonitor`: CreateLoadBalancerMonitor200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancerMonitor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerMonitorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createLoadBalancerMonitorRequest** | [**CreateLoadBalancerMonitorRequest**](CreateLoadBalancerMonitorRequest.md) |  | 

### Return type

[**CreateLoadBalancerMonitor200Response**](CreateLoadBalancerMonitor200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLoadBalancerPool

> CreateLoadBalancerPool200Response CreateLoadBalancerPool(ctx, loadBalancerId).CreateLoadBalancerPoolRequest(createLoadBalancerPoolRequest).Execute()

Create a Load Balancer Pool



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	createLoadBalancerPoolRequest := *openapiclient.NewCreateLoadBalancerPoolRequest() // CreateLoadBalancerPoolRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancerPool(context.Background(), loadBalancerId).CreateLoadBalancerPoolRequest(createLoadBalancerPoolRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancerPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancerPool`: CreateLoadBalancerPool200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancerPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createLoadBalancerPoolRequest** | [**CreateLoadBalancerPoolRequest**](CreateLoadBalancerPoolRequest.md) |  | 

### Return type

[**CreateLoadBalancerPool200Response**](CreateLoadBalancerPool200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLoadBalancerPoolNode

> CreateLoadBalancerPoolNode200Response CreateLoadBalancerPoolNode(ctx, loadBalancerPoolId).CreateLoadBalancerPoolNodeRequest(createLoadBalancerPoolNodeRequest).Execute()

Create a Load Balancer Pool Node



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
	loadBalancerPoolId := int64(4) // int64 | Load Balancer Pool ID
	createLoadBalancerPoolNodeRequest := *openapiclient.NewCreateLoadBalancerPoolNodeRequest() // CreateLoadBalancerPoolNodeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancerPoolNode(context.Background(), loadBalancerPoolId).CreateLoadBalancerPoolNodeRequest(createLoadBalancerPoolNodeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancerPoolNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancerPoolNode`: CreateLoadBalancerPoolNode200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancerPoolNode`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerPoolId** | **int64** | Load Balancer Pool ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerPoolNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createLoadBalancerPoolNodeRequest** | [**CreateLoadBalancerPoolNodeRequest**](CreateLoadBalancerPoolNodeRequest.md) |  | 

### Return type

[**CreateLoadBalancerPoolNode200Response**](CreateLoadBalancerPoolNode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLoadBalancerProfile

> CreateLoadBalancerProfile200Response CreateLoadBalancerProfile(ctx, loadBalancerId).CreateLoadBalancerProfileRequest(createLoadBalancerProfileRequest).Execute()

Create a Load Balancer Profile



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	createLoadBalancerProfileRequest := *openapiclient.NewCreateLoadBalancerProfileRequest() // CreateLoadBalancerProfileRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancerProfile(context.Background(), loadBalancerId).CreateLoadBalancerProfileRequest(createLoadBalancerProfileRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancerProfile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancerProfile`: CreateLoadBalancerProfile200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancerProfile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerProfileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createLoadBalancerProfileRequest** | [**CreateLoadBalancerProfileRequest**](CreateLoadBalancerProfileRequest.md) |  | 

### Return type

[**CreateLoadBalancerProfile200Response**](CreateLoadBalancerProfile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLoadBalancerVirtualServer

> CreateLoadBalancerVirtualServer200Response CreateLoadBalancerVirtualServer(ctx, loadBalancerId).CreateLoadBalancerVirtualServerRequest(createLoadBalancerVirtualServerRequest).Execute()

Create a Load Balancer Virtual Server



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	createLoadBalancerVirtualServerRequest := *openapiclient.NewCreateLoadBalancerVirtualServerRequest() // CreateLoadBalancerVirtualServerRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.CreateLoadBalancerVirtualServer(context.Background(), loadBalancerId).CreateLoadBalancerVirtualServerRequest(createLoadBalancerVirtualServerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.CreateLoadBalancerVirtualServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLoadBalancerVirtualServer`: CreateLoadBalancerVirtualServer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.CreateLoadBalancerVirtualServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateLoadBalancerVirtualServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createLoadBalancerVirtualServerRequest** | [**CreateLoadBalancerVirtualServerRequest**](CreateLoadBalancerVirtualServerRequest.md) |  | 

### Return type

[**CreateLoadBalancerVirtualServer200Response**](CreateLoadBalancerVirtualServer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancer

> DeleteLoadBalancer200Response DeleteLoadBalancer(ctx, loadBalancerId).Execute()

Delete a Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancer(context.Background(), loadBalancerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancer`: DeleteLoadBalancer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteLoadBalancer200Response**](DeleteLoadBalancer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancerMonitor

> DeleteLoadBalancerMonitor200Response DeleteLoadBalancerMonitor(ctx, loadBalancerId, id).Execute()

Delete a Load Balancer Monitor



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancerMonitor(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancerMonitor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancerMonitor`: DeleteLoadBalancerMonitor200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancerMonitor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerMonitorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteLoadBalancerMonitor200Response**](DeleteLoadBalancerMonitor200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancerPool

> DeleteLoadBalancerPool200Response DeleteLoadBalancerPool(ctx, loadBalancerId, id).Execute()

Delete a Load Balancer Pool



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancerPool(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancerPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancerPool`: DeleteLoadBalancerPool200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancerPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteLoadBalancerPool200Response**](DeleteLoadBalancerPool200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancerPoolNode

> DeleteLoadBalancerPoolNode200Response DeleteLoadBalancerPoolNode(ctx, loadBalancerPoolId, id).Execute()

Delete a Load Balancer Pool Node



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
	loadBalancerPoolId := int64(4) // int64 | Load Balancer Pool ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancerPoolNode(context.Background(), loadBalancerPoolId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancerPoolNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancerPoolNode`: DeleteLoadBalancerPoolNode200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancerPoolNode`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerPoolId** | **int64** | Load Balancer Pool ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerPoolNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteLoadBalancerPoolNode200Response**](DeleteLoadBalancerPoolNode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancerProfile

> DeleteLoadBalancerProfile200Response DeleteLoadBalancerProfile(ctx, loadBalancerId, id).Execute()

Delete a Load Balancer Profile



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancerProfile(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancerProfile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancerProfile`: DeleteLoadBalancerProfile200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancerProfile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerProfileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteLoadBalancerProfile200Response**](DeleteLoadBalancerProfile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLoadBalancerVirtualServer

> DeleteLoadBalancerVirtualServer200Response DeleteLoadBalancerVirtualServer(ctx, loadBalancerId, id).Execute()

Delete a Load Balancer Virtual Server



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.DeleteLoadBalancerVirtualServer(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.DeleteLoadBalancerVirtualServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLoadBalancerVirtualServer`: DeleteLoadBalancerVirtualServer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.DeleteLoadBalancerVirtualServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLoadBalancerVirtualServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteLoadBalancerVirtualServer200Response**](DeleteLoadBalancerVirtualServer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancer

> GetLoadBalancer200Response GetLoadBalancer(ctx, loadBalancerId).Execute()

Get a Specific Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancer(context.Background(), loadBalancerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancer`: GetLoadBalancer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetLoadBalancer200Response**](GetLoadBalancer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerMonitor

> GetLoadBalancerMonitor200Response GetLoadBalancerMonitor(ctx, loadBalancerId, id).Execute()

Get a Specific Load Balancer Monitor



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerMonitor(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerMonitor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerMonitor`: GetLoadBalancerMonitor200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerMonitor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerMonitorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetLoadBalancerMonitor200Response**](GetLoadBalancerMonitor200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerPool

> GetLoadBalancerPool200Response GetLoadBalancerPool(ctx, loadBalancerId, id).Execute()

Get a Specific Load Balancer Pool



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerPool(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerPool`: GetLoadBalancerPool200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetLoadBalancerPool200Response**](GetLoadBalancerPool200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerPoolNode

> GetLoadBalancerPoolNode200Response GetLoadBalancerPoolNode(ctx, loadBalancerPoolId, id).Execute()

Get a Specific Load Balancer Pool Node



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
	loadBalancerPoolId := int64(4) // int64 | Load Balancer Pool ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerPoolNode(context.Background(), loadBalancerPoolId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerPoolNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerPoolNode`: GetLoadBalancerPoolNode200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerPoolNode`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerPoolId** | **int64** | Load Balancer Pool ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerPoolNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetLoadBalancerPoolNode200Response**](GetLoadBalancerPoolNode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerProfile

> GetLoadBalancerProfile200Response GetLoadBalancerProfile(ctx, loadBalancerId, id).Execute()

Get a Specific Load Balancer Profile



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerProfile(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerProfile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerProfile`: GetLoadBalancerProfile200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerProfile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerProfileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetLoadBalancerProfile200Response**](GetLoadBalancerProfile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerType

> GetLoadBalancerType200Response GetLoadBalancerType(ctx, id).Execute()

Get a Specific Load Balancer Type



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
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerType(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerType``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerType`: GetLoadBalancerType200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerType`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerTypeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetLoadBalancerType200Response**](GetLoadBalancerType200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLoadBalancerVirtualServer

> GetLoadBalancerVirtualServer200Response GetLoadBalancerVirtualServer(ctx, loadBalancerId, id).Execute()

Get a Specific Load Balancer Virtual Server



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.GetLoadBalancerVirtualServer(context.Background(), loadBalancerId, id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.GetLoadBalancerVirtualServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLoadBalancerVirtualServer`: GetLoadBalancerVirtualServer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.GetLoadBalancerVirtualServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLoadBalancerVirtualServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**GetLoadBalancerVirtualServer200Response**](GetLoadBalancerVirtualServer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerMonitors

> ListLoadBalancerMonitors200Response ListLoadBalancerMonitors(ctx, loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()

Get All Load Balancer Monitors For Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerMonitors(context.Background(), loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerMonitors``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerMonitors`: ListLoadBalancerMonitors200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerMonitors`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerMonitorsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListLoadBalancerMonitors200Response**](ListLoadBalancerMonitors200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerPoolNodes

> ListLoadBalancerPoolNodes200Response ListLoadBalancerPoolNodes(ctx, loadBalancerPoolId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Execute()

Get All Load Balancer Pool Nodes For Load Balancer Pool



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
	loadBalancerPoolId := int64(4) // int64 | Load Balancer Pool ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerPoolNodes(context.Background(), loadBalancerPoolId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerPoolNodes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerPoolNodes`: ListLoadBalancerPoolNodes200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerPoolNodes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerPoolId** | **int64** | Load Balancer Pool ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerPoolNodesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListLoadBalancerPoolNodes200Response**](ListLoadBalancerPoolNodes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerPools

> ListLoadBalancerPools200Response ListLoadBalancerPools(ctx, loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()

Get All Load Balancer Pools For Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerPools(context.Background(), loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerPools`: ListLoadBalancerPools200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerPools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerPoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListLoadBalancerPools200Response**](ListLoadBalancerPools200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerProfiles

> ListLoadBalancerProfiles200Response ListLoadBalancerProfiles(ctx, loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()

Get All Load Balancer Profiles For Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerProfiles(context.Background(), loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerProfiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerProfiles`: ListLoadBalancerProfiles200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerProfiles`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerProfilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListLoadBalancerProfiles200Response**](ListLoadBalancerProfiles200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerTypes

> ListLoadBalancerTypes200Response ListLoadBalancerTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).OptionTypes(optionTypes).Phrase(phrase).Name(name).Code(code).Execute()

Get All Load Balancer Types



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
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	optionTypes := true // bool | Pass true to include optionTypes in the response for each entry. (optional) (default to false)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	code := "azr" // string | If specified will return an exact match on code (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).OptionTypes(optionTypes).Phrase(phrase).Name(name).Code(code).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerTypes`: ListLoadBalancerTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **optionTypes** | **bool** | Pass true to include optionTypes in the response for each entry. | [default to false]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 

### Return type

[**ListLoadBalancerTypes200Response**](ListLoadBalancerTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancerVirtualServers

> ListLoadBalancerVirtualServers200Response ListLoadBalancerVirtualServers(ctx, loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).VipName(vipName).VipAddress(vipAddress).VipHostname(vipHostname).Execute()

Get All Load Balancer Virtual Servers For Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	vipName := "lb-114" // string | If specified will return an exact match on vipName (optional)
	vipAddress := "192.168.123.44" // string | If specified will return an exact match on vipAddress (optional)
	vipHostname := "mylb" // string | If specified will return an exact match on vipHostname (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancerVirtualServers(context.Background(), loadBalancerId).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).VipName(vipName).VipAddress(vipAddress).VipHostname(vipHostname).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancerVirtualServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancerVirtualServers`: ListLoadBalancerVirtualServers200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancerVirtualServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancerVirtualServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **vipName** | **string** | If specified will return an exact match on vipName | 
 **vipAddress** | **string** | If specified will return an exact match on vipAddress | 
 **vipHostname** | **string** | If specified will return an exact match on vipHostname | 

### Return type

[**ListLoadBalancerVirtualServers200Response**](ListLoadBalancerVirtualServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLoadBalancers

> ListLoadBalancers200Response ListLoadBalancers(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()

Get All Load Balancers



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
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.ListLoadBalancers(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.ListLoadBalancers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLoadBalancers`: ListLoadBalancers200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.ListLoadBalancers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListLoadBalancersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListLoadBalancers200Response**](ListLoadBalancers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RefreshLoadBalancer

> RefreshLoadBalancer200Response RefreshLoadBalancer(ctx, loadBalancerId).Execute()

Refresh a Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.RefreshLoadBalancer(context.Background(), loadBalancerId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.RefreshLoadBalancer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RefreshLoadBalancer`: RefreshLoadBalancer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.RefreshLoadBalancer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRefreshLoadBalancerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RefreshLoadBalancer200Response**](RefreshLoadBalancer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancer

> CreateLoadBalancer200Response UpdateLoadBalancer(ctx, loadBalancerId).UpdateLoadBalancerRequest(updateLoadBalancerRequest).Execute()

Update a Load Balancer



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	updateLoadBalancerRequest := *openapiclient.NewUpdateLoadBalancerRequest() // UpdateLoadBalancerRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancer(context.Background(), loadBalancerId).UpdateLoadBalancerRequest(updateLoadBalancerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancer`: CreateLoadBalancer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateLoadBalancerRequest** | [**UpdateLoadBalancerRequest**](UpdateLoadBalancerRequest.md) |  | 

### Return type

[**CreateLoadBalancer200Response**](CreateLoadBalancer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancerMonitor

> UpdateLoadBalancerMonitor200Response UpdateLoadBalancerMonitor(ctx, loadBalancerId, id).UpdateLoadBalancerMonitorRequest(updateLoadBalancerMonitorRequest).Execute()

Update a Load Balancer Monitor



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateLoadBalancerMonitorRequest := *openapiclient.NewUpdateLoadBalancerMonitorRequest() // UpdateLoadBalancerMonitorRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancerMonitor(context.Background(), loadBalancerId, id).UpdateLoadBalancerMonitorRequest(updateLoadBalancerMonitorRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancerMonitor``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancerMonitor`: UpdateLoadBalancerMonitor200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancerMonitor`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerMonitorRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateLoadBalancerMonitorRequest** | [**UpdateLoadBalancerMonitorRequest**](UpdateLoadBalancerMonitorRequest.md) |  | 

### Return type

[**UpdateLoadBalancerMonitor200Response**](UpdateLoadBalancerMonitor200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancerPool

> UpdateLoadBalancerPool200Response UpdateLoadBalancerPool(ctx, loadBalancerId, id).UpdateLoadBalancerPoolRequest(updateLoadBalancerPoolRequest).Execute()

Update a Load Balancer Pool



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateLoadBalancerPoolRequest := *openapiclient.NewUpdateLoadBalancerPoolRequest() // UpdateLoadBalancerPoolRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancerPool(context.Background(), loadBalancerId, id).UpdateLoadBalancerPoolRequest(updateLoadBalancerPoolRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancerPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancerPool`: UpdateLoadBalancerPool200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancerPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateLoadBalancerPoolRequest** | [**UpdateLoadBalancerPoolRequest**](UpdateLoadBalancerPoolRequest.md) |  | 

### Return type

[**UpdateLoadBalancerPool200Response**](UpdateLoadBalancerPool200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancerPoolNode

> UpdateLoadBalancerPoolNode200Response UpdateLoadBalancerPoolNode(ctx, loadBalancerPoolId, id).UpdateLoadBalancerPoolNodeRequest(updateLoadBalancerPoolNodeRequest).Execute()

Update a Load Balancer Pool Node



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
	loadBalancerPoolId := int64(4) // int64 | Load Balancer Pool ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateLoadBalancerPoolNodeRequest := *openapiclient.NewUpdateLoadBalancerPoolNodeRequest() // UpdateLoadBalancerPoolNodeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancerPoolNode(context.Background(), loadBalancerPoolId, id).UpdateLoadBalancerPoolNodeRequest(updateLoadBalancerPoolNodeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancerPoolNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancerPoolNode`: UpdateLoadBalancerPoolNode200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancerPoolNode`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerPoolId** | **int64** | Load Balancer Pool ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerPoolNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateLoadBalancerPoolNodeRequest** | [**UpdateLoadBalancerPoolNodeRequest**](UpdateLoadBalancerPoolNodeRequest.md) |  | 

### Return type

[**UpdateLoadBalancerPoolNode200Response**](UpdateLoadBalancerPoolNode200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancerProfile

> UpdateLoadBalancerProfile200Response UpdateLoadBalancerProfile(ctx, loadBalancerId, id).UpdateLoadBalancerProfileRequest(updateLoadBalancerProfileRequest).Execute()

Update a Load Balancer Profile



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateLoadBalancerProfileRequest := *openapiclient.NewUpdateLoadBalancerProfileRequest() // UpdateLoadBalancerProfileRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancerProfile(context.Background(), loadBalancerId, id).UpdateLoadBalancerProfileRequest(updateLoadBalancerProfileRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancerProfile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancerProfile`: UpdateLoadBalancerProfile200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancerProfile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerProfileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateLoadBalancerProfileRequest** | [**UpdateLoadBalancerProfileRequest**](UpdateLoadBalancerProfileRequest.md) |  | 

### Return type

[**UpdateLoadBalancerProfile200Response**](UpdateLoadBalancerProfile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateLoadBalancerVirtualServer

> UpdateLoadBalancerVirtualServer200Response UpdateLoadBalancerVirtualServer(ctx, loadBalancerId, id).UpdateLoadBalancerVirtualServerRequest(updateLoadBalancerVirtualServerRequest).Execute()

Update a Load Balancer Virtual Server



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
	loadBalancerId := int64(4) // int64 | Load Balancer ID
	id := int64(1) // int64 | Morpheus ID of the Object being referenced
	updateLoadBalancerVirtualServerRequest := *openapiclient.NewUpdateLoadBalancerVirtualServerRequest() // UpdateLoadBalancerVirtualServerRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LoadBalancersAPI.UpdateLoadBalancerVirtualServer(context.Background(), loadBalancerId, id).UpdateLoadBalancerVirtualServerRequest(updateLoadBalancerVirtualServerRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LoadBalancersAPI.UpdateLoadBalancerVirtualServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLoadBalancerVirtualServer`: UpdateLoadBalancerVirtualServer200Response
	fmt.Fprintf(os.Stdout, "Response from `LoadBalancersAPI.UpdateLoadBalancerVirtualServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**loadBalancerId** | **int64** | Load Balancer ID | 
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLoadBalancerVirtualServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **updateLoadBalancerVirtualServerRequest** | [**UpdateLoadBalancerVirtualServerRequest**](UpdateLoadBalancerVirtualServerRequest.md) |  | 

### Return type

[**UpdateLoadBalancerVirtualServer200Response**](UpdateLoadBalancerVirtualServer200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

