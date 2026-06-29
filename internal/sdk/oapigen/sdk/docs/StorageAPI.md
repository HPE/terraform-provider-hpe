# \StorageAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddStorageBuckets**](StorageAPI.md#AddStorageBuckets) | **Post** /api/storage-buckets | Creates a Storage Bucket
[**AddStorageServers**](StorageAPI.md#AddStorageServers) | **Post** /api/storage-servers | Creates a Storage Server
[**AddStorageVolumes**](StorageAPI.md#AddStorageVolumes) | **Post** /api/storage-volumes | Creates a Storage Volume
[**AttachServerVolume**](StorageAPI.md#AttachServerVolume) | **Put** /api/servers/{id}/volumes/{volumeId}/attach | Attach existing storage volume
[**DetachServerVolume**](StorageAPI.md#DetachServerVolume) | **Put** /api/servers/{id}/volumes/{volumeId}/detach | Detach storage volume
[**GetStorageBuckets**](StorageAPI.md#GetStorageBuckets) | **Get** /api/storage-buckets/{id} | Retrieves a Specific Storage Bucket
[**GetStorageServerTypes**](StorageAPI.md#GetStorageServerTypes) | **Get** /api/storage-server-types/{id} | Retrieves a Specific Storage Server Type
[**GetStorageServers**](StorageAPI.md#GetStorageServers) | **Get** /api/storage-servers/{id} | Retrieves a Specific Storage Server
[**GetStorageVolumeTypes**](StorageAPI.md#GetStorageVolumeTypes) | **Get** /api/storage-volume-types/{id} | Retrieves a Specific Storage Volume Type
[**GetStorageVolumes**](StorageAPI.md#GetStorageVolumes) | **Get** /api/storage-volumes/{id} | Retrieves a Specific Storage Volume
[**ListStorageBuckets**](StorageAPI.md#ListStorageBuckets) | **Get** /api/storage-buckets | Retrieves all Storage Buckets
[**ListStorageServerTypes**](StorageAPI.md#ListStorageServerTypes) | **Get** /api/storage-server-types | Retrieves all Storage Server Types
[**ListStorageServers**](StorageAPI.md#ListStorageServers) | **Get** /api/storage-servers | Retrieves all Storage Servers
[**ListStorageVolumeTypes**](StorageAPI.md#ListStorageVolumeTypes) | **Get** /api/storage-volume-types | Retrieves all Storage Volume Types
[**ListStorageVolumes**](StorageAPI.md#ListStorageVolumes) | **Get** /api/storage-volumes | Retrieves all Storage Volumes
[**RemoveStorageBuckets**](StorageAPI.md#RemoveStorageBuckets) | **Delete** /api/storage-buckets/{id} | Deletes a Storage Bucket
[**RemoveStorageServers**](StorageAPI.md#RemoveStorageServers) | **Delete** /api/storage-servers/{id} | Deletes a Storage Server
[**RemoveStorageVolumes**](StorageAPI.md#RemoveStorageVolumes) | **Delete** /api/storage-volumes/{id} | Deletes a Storage Volume
[**UpdateStorageBuckets**](StorageAPI.md#UpdateStorageBuckets) | **Put** /api/storage-buckets/{id} | Updates a Storage Bucket
[**UpdateStorageServers**](StorageAPI.md#UpdateStorageServers) | **Put** /api/storage-servers/{id} | Updates a Storage Server
[**UpdateStorageVolumes**](StorageAPI.md#UpdateStorageVolumes) | **Put** /api/storage-volumes/{id} | Updates a Storage Volume



## AddStorageBuckets

> AddStorageBuckets200Response AddStorageBuckets(ctx).AddStorageBucketsRequest(addStorageBucketsRequest).Execute()

Creates a Storage Bucket



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
	addStorageBucketsRequest := *openapiclient.NewAddStorageBucketsRequest(*openapiclient.NewAddStorageBucketsRequestStorageBucket("Name_example", "ProviderType_example", "myBucket", openapiclient.addStorageBuckets_request_storageBucket_config{AddStorageBucketsRequestStorageBucketConfigOneOf: openapiclient.NewAddStorageBucketsRequestStorageBucketConfigOneOf()})) // AddStorageBucketsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.AddStorageBuckets(context.Background()).AddStorageBucketsRequest(addStorageBucketsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.AddStorageBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddStorageBuckets`: AddStorageBuckets200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.AddStorageBuckets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddStorageBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addStorageBucketsRequest** | [**AddStorageBucketsRequest**](AddStorageBucketsRequest.md) |  | 

### Return type

[**AddStorageBuckets200Response**](AddStorageBuckets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddStorageServers

> AddStorageServers200Response AddStorageServers(ctx).AddStorageServersRequest(addStorageServersRequest).Execute()

Creates a Storage Server



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
	addStorageServersRequest := *openapiclient.NewAddStorageServersRequest(*openapiclient.NewAddStorageServersRequestStorageServer("Name_example", "Type_example", map[string]interface{}(123))) // AddStorageServersRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.AddStorageServers(context.Background()).AddStorageServersRequest(addStorageServersRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.AddStorageServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddStorageServers`: AddStorageServers200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.AddStorageServers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddStorageServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addStorageServersRequest** | [**AddStorageServersRequest**](AddStorageServersRequest.md) |  | 

### Return type

[**AddStorageServers200Response**](AddStorageServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddStorageVolumes

> AddStorageVolumes200Response AddStorageVolumes(ctx).AddStorageVolumesRequest(addStorageVolumesRequest).Execute()

Creates a Storage Volume



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
	addStorageVolumesRequest := *openapiclient.NewAddStorageVolumesRequest(*openapiclient.NewAddStorageVolumesRequestStorageVolume("Name_example", "Type_example", *openapiclient.NewAddStorageVolumesRequestStorageVolumeStorageServer(int64(123)))) // AddStorageVolumesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.AddStorageVolumes(context.Background()).AddStorageVolumesRequest(addStorageVolumesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.AddStorageVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddStorageVolumes`: AddStorageVolumes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.AddStorageVolumes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddStorageVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addStorageVolumesRequest** | [**AddStorageVolumesRequest**](AddStorageVolumesRequest.md) |  | 

### Return type

[**AddStorageVolumes200Response**](AddStorageVolumes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachServerVolume

> AttachServerVolume200Response AttachServerVolume(ctx, id, volumeId).AttachServerVolumeRequest(attachServerVolumeRequest).Execute()

Attach existing storage volume



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
	id := openapiclient.updateHost_id_parameter{Int64: new(int64)} // UpdateHostIdParameter | The `id` or `uuid` to identify the server record.
	volumeId := openapiclient.attachServerVolume_volumeId_parameter{Int64: new(int64)} // AttachServerVolumeVolumeIdParameter | The `id` or `uuid` to identify the storage volume record.
	attachServerVolumeRequest := *openapiclient.NewAttachServerVolumeRequest(*openapiclient.NewAttachServerVolumeRequestMountPoint(*openapiclient.NewAttachServerVolumeRequestMountPointController(int64(123)), "UnitNumber_example")) // AttachServerVolumeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.AttachServerVolume(context.Background(), id, volumeId).AttachServerVolumeRequest(attachServerVolumeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.AttachServerVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachServerVolume`: AttachServerVolume200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.AttachServerVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**UpdateHostIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the server record. | 
**volumeId** | [**AttachServerVolumeVolumeIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the storage volume record. | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachServerVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **attachServerVolumeRequest** | [**AttachServerVolumeRequest**](AttachServerVolumeRequest.md) |  | 

### Return type

[**AttachServerVolume200Response**](AttachServerVolume200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DetachServerVolume

> DetachServerVolume200Response DetachServerVolume(ctx, id, volumeId).Force(force).Execute()

Detach storage volume



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
	id := openapiclient.updateHost_id_parameter{Int64: new(int64)} // UpdateHostIdParameter | The `id` or `uuid` to identify the server record.
	volumeId := openapiclient.detachServerVolume_volumeId_parameter{Int64: new(int64)} // DetachServerVolumeVolumeIdParameter | The `id` or `uuid` to identify the storage volume record.
	force := true // bool | Force the detachment of the volume (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.DetachServerVolume(context.Background(), id, volumeId).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.DetachServerVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachServerVolume`: DetachServerVolume200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.DetachServerVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**UpdateHostIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the server record. | 
**volumeId** | [**DetachServerVolumeVolumeIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the storage volume record. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachServerVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **bool** | Force the detachment of the volume | [default to false]

### Return type

[**DetachServerVolume200Response**](DetachServerVolume200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStorageBuckets

> GetStorageBuckets200Response GetStorageBuckets(ctx, id).Execute()

Retrieves a Specific Storage Bucket



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
	resp, r, err := apiClient.StorageAPI.GetStorageBuckets(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.GetStorageBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStorageBuckets`: GetStorageBuckets200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.GetStorageBuckets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStorageBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStorageBuckets200Response**](GetStorageBuckets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStorageServerTypes

> GetStorageServerTypes200Response GetStorageServerTypes(ctx, id).Execute()

Retrieves a Specific Storage Server Type



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
	resp, r, err := apiClient.StorageAPI.GetStorageServerTypes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.GetStorageServerTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStorageServerTypes`: GetStorageServerTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.GetStorageServerTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStorageServerTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStorageServerTypes200Response**](GetStorageServerTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStorageServers

> GetStorageServers200Response GetStorageServers(ctx, id).Execute()

Retrieves a Specific Storage Server



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
	resp, r, err := apiClient.StorageAPI.GetStorageServers(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.GetStorageServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStorageServers`: GetStorageServers200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.GetStorageServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStorageServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStorageServers200Response**](GetStorageServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStorageVolumeTypes

> GetStorageVolumeTypes200Response GetStorageVolumeTypes(ctx, id).Execute()

Retrieves a Specific Storage Volume Type



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
	resp, r, err := apiClient.StorageAPI.GetStorageVolumeTypes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.GetStorageVolumeTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStorageVolumeTypes`: GetStorageVolumeTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.GetStorageVolumeTypes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStorageVolumeTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStorageVolumeTypes200Response**](GetStorageVolumeTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetStorageVolumes

> GetStorageVolumes200Response GetStorageVolumes(ctx, id).Execute()

Retrieves a Specific Storage Volume



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
	id := openapiclient.getStorageVolumes_id_parameter{Int64: new(int64)} // GetStorageVolumesIdParameter | The `id` or `uuid` to identify the storage volume record.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.GetStorageVolumes(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.GetStorageVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetStorageVolumes`: GetStorageVolumes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.GetStorageVolumes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**GetStorageVolumesIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the storage volume record. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetStorageVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetStorageVolumes200Response**](GetStorageVolumes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageBuckets

> ListStorageBuckets200Response ListStorageBuckets(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Storage Buckets



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.ListStorageBuckets(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.ListStorageBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListStorageBuckets`: ListStorageBuckets200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.ListStorageBuckets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListStorageBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListStorageBuckets200Response**](ListStorageBuckets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageServerTypes

> ListStorageServerTypes200Response ListStorageServerTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Code(code).Phrase(phrase).Execute()

Retrieves all Storage Server Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.ListStorageServerTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Code(code).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.ListStorageServerTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListStorageServerTypes`: ListStorageServerTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.ListStorageServerTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListStorageServerTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListStorageServerTypes200Response**](ListStorageServerTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageServers

> ListStorageServers200Response ListStorageServers(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()

Retrieves all Storage Servers



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.ListStorageServers(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.ListStorageServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListStorageServers`: ListStorageServers200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.ListStorageServers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListStorageServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 

### Return type

[**ListStorageServers200Response**](ListStorageServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageVolumeTypes

> ListStorageVolumeTypes200Response ListStorageVolumeTypes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Code(code).Phrase(phrase).Execute()

Retrieves all Storage Volume Types



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
	code := "azr" // string | If specified will return an exact match on code (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.ListStorageVolumeTypes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Name(name).Code(code).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.ListStorageVolumeTypes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListStorageVolumeTypes`: ListStorageVolumeTypes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.ListStorageVolumeTypes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListStorageVolumeTypesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **name** | **string** | Filter by name | 
 **code** | **string** | If specified will return an exact match on code | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListStorageVolumeTypes200Response**](ListStorageVolumeTypes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListStorageVolumes

> ListStorageVolumes200Response ListStorageVolumes(ctx).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Uuid(uuid).Execute()

Retrieves all Storage Volumes



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
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	name := "example" // string | Filter by name (optional)
	uuid := "uuid_example" // string | Filter by UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.ListStorageVolumes(context.Background()).Max(max).Offset(offset).Sort(sort).Direction(direction).Phrase(phrase).Name(name).Uuid(uuid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.ListStorageVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListStorageVolumes`: ListStorageVolumes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.ListStorageVolumes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListStorageVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **name** | **string** | Filter by name | 
 **uuid** | **string** | Filter by UUID | 

### Return type

[**ListStorageVolumes200Response**](ListStorageVolumes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveStorageBuckets

> RemoveExecuteSchedules200Response RemoveStorageBuckets(ctx, id).RemoveResources(removeResources).Execute()

Deletes a Storage Bucket



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
	removeResources := true // bool | Removing resources will permanently delete the bucket and all files within the bucket. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.RemoveStorageBuckets(context.Background(), id).RemoveResources(removeResources).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.RemoveStorageBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveStorageBuckets`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.RemoveStorageBuckets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveStorageBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **removeResources** | **bool** | Removing resources will permanently delete the bucket and all files within the bucket. | [default to false]

### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveStorageServers

> RemoveExecuteSchedules200Response RemoveStorageServers(ctx, id).Execute()

Deletes a Storage Server



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
	resp, r, err := apiClient.StorageAPI.RemoveStorageServers(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.RemoveStorageServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveStorageServers`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.RemoveStorageServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveStorageServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveStorageVolumes

> RemoveExecuteSchedules200Response RemoveStorageVolumes(ctx, id).Force(force).Execute()

Deletes a Storage Volume



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
	id := openapiclient.updateStorageVolumes_id_parameter{Int64: new(int64)} // UpdateStorageVolumesIdParameter | The `id` or `uuid` to identify the storage volume record.
	force := "on" // string | Force Delete (optional) (default to "off")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.RemoveStorageVolumes(context.Background(), id).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.RemoveStorageVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveStorageVolumes`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.RemoveStorageVolumes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**UpdateStorageVolumesIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the storage volume record. | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveStorageVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **force** | **string** | Force Delete | [default to &quot;off&quot;]

### Return type

[**RemoveExecuteSchedules200Response**](RemoveExecuteSchedules200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateStorageBuckets

> UpdateStorageBuckets200Response UpdateStorageBuckets(ctx, id).UpdateStorageBucketsRequest(updateStorageBucketsRequest).Execute()

Updates a Storage Bucket



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
	updateStorageBucketsRequest := *openapiclient.NewUpdateStorageBucketsRequest(*openapiclient.NewUpdateStorageBucketsRequestStorageBucket()) // UpdateStorageBucketsRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.UpdateStorageBuckets(context.Background(), id).UpdateStorageBucketsRequest(updateStorageBucketsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.UpdateStorageBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateStorageBuckets`: UpdateStorageBuckets200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.UpdateStorageBuckets`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateStorageBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateStorageBucketsRequest** | [**UpdateStorageBucketsRequest**](UpdateStorageBucketsRequest.md) |  | 

### Return type

[**UpdateStorageBuckets200Response**](UpdateStorageBuckets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateStorageServers

> UpdateStorageServers200Response UpdateStorageServers(ctx, id).UpdateStorageServersRequest(updateStorageServersRequest).Execute()

Updates a Storage Server



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
	updateStorageServersRequest := *openapiclient.NewUpdateStorageServersRequest(*openapiclient.NewUpdateStorageServersRequestStorageServer()) // UpdateStorageServersRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.UpdateStorageServers(context.Background(), id).UpdateStorageServersRequest(updateStorageServersRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.UpdateStorageServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateStorageServers`: UpdateStorageServers200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.UpdateStorageServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateStorageServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateStorageServersRequest** | [**UpdateStorageServersRequest**](UpdateStorageServersRequest.md) |  | 

### Return type

[**UpdateStorageServers200Response**](UpdateStorageServers200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateStorageVolumes

> UpdateStorageVolumes200Response UpdateStorageVolumes(ctx, id).UpdateStorageVolumesRequest(updateStorageVolumesRequest).Execute()

Updates a Storage Volume



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
	id := openapiclient.updateStorageVolumes_id_parameter{Int64: new(int64)} // UpdateStorageVolumesIdParameter | The `id` or `uuid` to identify the storage volume record.
	updateStorageVolumesRequest := *openapiclient.NewUpdateStorageVolumesRequest(*openapiclient.NewUpdateStorageVolumesRequestStorageVolume()) // UpdateStorageVolumesRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.StorageAPI.UpdateStorageVolumes(context.Background(), id).UpdateStorageVolumesRequest(updateStorageVolumesRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StorageAPI.UpdateStorageVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateStorageVolumes`: UpdateStorageVolumes200Response
	fmt.Fprintf(os.Stdout, "Response from `StorageAPI.UpdateStorageVolumes`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**UpdateStorageVolumesIdParameter**](.md) | The &#x60;id&#x60; or &#x60;uuid&#x60; to identify the storage volume record. | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateStorageVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateStorageVolumesRequest** | [**UpdateStorageVolumesRequest**](UpdateStorageVolumesRequest.md) |  | 

### Return type

[**UpdateStorageVolumes200Response**](UpdateStorageVolumes200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

