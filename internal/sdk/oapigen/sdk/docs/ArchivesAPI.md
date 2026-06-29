# \ArchivesAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddArchiveBucket**](ArchivesAPI.md#AddArchiveBucket) | **Post** /api/archives/buckets | Create an Archive Bucket
[**AddArchiveFile**](ArchivesAPI.md#AddArchiveFile) | **Post** /api/archives/buckets/{bucket}/files/{filepath} | Upload Archive File
[**AddArchiveFileLink**](ArchivesAPI.md#AddArchiveFileLink) | **Post** /api/archives/files/{id}/links | Create an Archive File Link
[**DeleteArchiveBucket**](ArchivesAPI.md#DeleteArchiveBucket) | **Delete** /api/archives/buckets/{id} | Delete an Archive Bucket
[**DeleteArchiveFile**](ArchivesAPI.md#DeleteArchiveFile) | **Delete** /api/archives/files/{id} | Delete Archive File
[**DeleteArchiveFileLink**](ArchivesAPI.md#DeleteArchiveFileLink) | **Delete** /api/archives/files/{id}/links/{linkId} | Delete an Archive File Link
[**GetArchiveBucket**](ArchivesAPI.md#GetArchiveBucket) | **Get** /api/archives/buckets/{id} | Get a Specific Archive Bucket
[**GetArchiveFile**](ArchivesAPI.md#GetArchiveFile) | **Get** /api/archives/download/{bucket}/{filepath} | Download an Archive File
[**GetArchiveFileDetail**](ArchivesAPI.md#GetArchiveFileDetail) | **Get** /api/archives/files/{id} | Get Archive File Details
[**GetArchiveFileLinks**](ArchivesAPI.md#GetArchiveFileLinks) | **Get** /api/archives/files/{id}/links | Get Archive File Links
[**GetArchivePublicFile**](ArchivesAPI.md#GetArchivePublicFile) | **Get** /api/public-archives/download/{bucket}/{filepath} | Download a Public Archive File
[**GetArchivePublicFileLink**](ArchivesAPI.md#GetArchivePublicFileLink) | **Get** /api/public-archives/link | Download an Archive File Link
[**ListArchiveBuckets**](ArchivesAPI.md#ListArchiveBuckets) | **Get** /api/archives/buckets | Get All Archive Buckets
[**ListArchiveFiles**](ArchivesAPI.md#ListArchiveFiles) | **Get** /api/archives/buckets/{bucket}/files/{filepath} | Get All Archive Files
[**UpdateArchiveBucket**](ArchivesAPI.md#UpdateArchiveBucket) | **Put** /api/archives/buckets/{id} | Update an Archive Bucket



## AddArchiveBucket

> AddArchiveBucket200Response AddArchiveBucket(ctx).AddArchiveBucketRequest(addArchiveBucketRequest).Execute()

Create an Archive Bucket



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
	addArchiveBucketRequest := *openapiclient.NewAddArchiveBucketRequest() // AddArchiveBucketRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.AddArchiveBucket(context.Background()).AddArchiveBucketRequest(addArchiveBucketRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.AddArchiveBucket``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddArchiveBucket`: AddArchiveBucket200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.AddArchiveBucket`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddArchiveBucketRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **addArchiveBucketRequest** | [**AddArchiveBucketRequest**](AddArchiveBucketRequest.md) |  | 

### Return type

[**AddArchiveBucket200Response**](AddArchiveBucket200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddArchiveFile

> AddArchiveFile200Response AddArchiveFile(ctx, bucket, filepath).Filename(filename).File(file).Execute()

Upload Archive File



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
	bucket := "bucket_example" // string | Bucket name
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")
	filename := "/path/to/file" // string | Specify a filename for archive file. The base filename of the uploaded file is the default. (optional)
	file := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.AddArchiveFile(context.Background(), bucket, filepath).Filename(filename).File(file).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.AddArchiveFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddArchiveFile`: AddArchiveFile200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.AddArchiveFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**bucket** | **string** | Bucket name | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiAddArchiveFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **filename** | **string** | Specify a filename for archive file. The base filename of the uploaded file is the default. | 
 **file** | ***os.File** |  | 

### Return type

[**AddArchiveFile200Response**](AddArchiveFile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AddArchiveFileLink

> AddArchiveFileLink200Response AddArchiveFileLink(ctx, id).ExpireSeconds(expireSeconds).Execute()

Create an Archive File Link



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
	expireSeconds := int64(600) // int64 | Time to live in seconds. 0 means do not expire. (optional) (default to 0)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.AddArchiveFileLink(context.Background(), id).ExpireSeconds(expireSeconds).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.AddArchiveFileLink``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddArchiveFileLink`: AddArchiveFileLink200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.AddArchiveFileLink`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddArchiveFileLinkRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **expireSeconds** | **int64** | Time to live in seconds. 0 means do not expire. | [default to 0]

### Return type

[**AddArchiveFileLink200Response**](AddArchiveFileLink200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteArchiveBucket

> DeleteArchiveBucket200Response DeleteArchiveBucket(ctx, id).Execute()

Delete an Archive Bucket



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
	resp, r, err := apiClient.ArchivesAPI.DeleteArchiveBucket(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.DeleteArchiveBucket``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteArchiveBucket`: DeleteArchiveBucket200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.DeleteArchiveBucket`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteArchiveBucketRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteArchiveBucket200Response**](DeleteArchiveBucket200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteArchiveFile

> DeleteArchiveFile200Response DeleteArchiveFile(ctx, id).Execute()

Delete Archive File



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
	resp, r, err := apiClient.ArchivesAPI.DeleteArchiveFile(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.DeleteArchiveFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteArchiveFile`: DeleteArchiveFile200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.DeleteArchiveFile`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteArchiveFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DeleteArchiveFile200Response**](DeleteArchiveFile200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteArchiveFileLink

> DeleteArchiveFileLink200Response DeleteArchiveFileLink(ctx, id, linkId).Execute()

Delete an Archive File Link



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
	linkId := int64(789) // int64 | The ID of the archive file link.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.DeleteArchiveFileLink(context.Background(), id, linkId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.DeleteArchiveFileLink``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteArchiveFileLink`: DeleteArchiveFileLink200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.DeleteArchiveFileLink`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 
**linkId** | **int64** | The ID of the archive file link. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteArchiveFileLinkRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DeleteArchiveFileLink200Response**](DeleteArchiveFileLink200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchiveBucket

> GetArchiveBucket200Response GetArchiveBucket(ctx, id).Execute()

Get a Specific Archive Bucket



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
	resp, r, err := apiClient.ArchivesAPI.GetArchiveBucket(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchiveBucket``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetArchiveBucket`: GetArchiveBucket200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.GetArchiveBucket`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetArchiveBucketRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetArchiveBucket200Response**](GetArchiveBucket200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchiveFile

> GetArchiveFile(ctx, bucket, filepath).Execute()

Download an Archive File



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
	bucket := "bucket_example" // string | Bucket name
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ArchivesAPI.GetArchiveFile(context.Background(), bucket, filepath).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchiveFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**bucket** | **string** | Bucket name | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiGetArchiveFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchiveFileDetail

> GetArchiveFileDetail200Response GetArchiveFileDetail(ctx, id).Execute()

Get Archive File Details



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
	resp, r, err := apiClient.ArchivesAPI.GetArchiveFileDetail(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchiveFileDetail``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetArchiveFileDetail`: GetArchiveFileDetail200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.GetArchiveFileDetail`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetArchiveFileDetailRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetArchiveFileDetail200Response**](GetArchiveFileDetail200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchiveFileLinks

> GetArchiveFileLinks200Response GetArchiveFileLinks(ctx, id).Execute()

Get Archive File Links



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
	resp, r, err := apiClient.ArchivesAPI.GetArchiveFileLinks(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchiveFileLinks``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetArchiveFileLinks`: GetArchiveFileLinks200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.GetArchiveFileLinks`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetArchiveFileLinksRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**GetArchiveFileLinks200Response**](GetArchiveFileLinks200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchivePublicFile

> GetArchivePublicFile(ctx, bucket, filepath).Execute()

Download a Public Archive File



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
	bucket := "bucket_example" // string | Bucket name
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ArchivesAPI.GetArchivePublicFile(context.Background(), bucket, filepath).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchivePublicFile``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**bucket** | **string** | Bucket name | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiGetArchivePublicFileRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetArchivePublicFileLink

> GetArchivePublicFileLink(ctx).S(s).Execute()

Download an Archive File Link



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
	s := "45a214fce9a546b9" // string | The secret access token for the archive file being downloaded.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ArchivesAPI.GetArchivePublicFileLink(context.Background()).S(s).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.GetArchivePublicFileLink``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetArchivePublicFileLinkRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **s** | **string** | The secret access token for the archive file being downloaded. | 

### Return type

 (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListArchiveBuckets

> ListArchiveBuckets200Response ListArchiveBuckets(ctx).Name(name).Phrase(phrase).Execute()

Get All Archive Buckets



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
	name := "example" // string | Filter by name (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.ListArchiveBuckets(context.Background()).Name(name).Phrase(phrase).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.ListArchiveBuckets``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListArchiveBuckets`: ListArchiveBuckets200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.ListArchiveBuckets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListArchiveBucketsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Filter by name | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 

### Return type

[**ListArchiveBuckets200Response**](ListArchiveBuckets200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListArchiveFiles

> ListArchiveFiles200Response ListArchiveFiles(ctx, bucket, filepath).Name(name).Phrase(phrase).FullTree(fullTree).Execute()

Get All Archive Files



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
	bucket := "bucket_example" // string | Bucket name
	filepath := "/config/environments/" // string | The path to to search for files under. Default is the root directory /. (default to "/")
	name := "test-config.yaml" // string | If specified will return an exact match on file name (optional)
	phrase := "test-%.yaml" // string | Search phrase for partial matches on file name, wildcard may be specified as %, eg. example-% This also searches for files under sub directories too.  (optional)
	fullTree := true // bool | Include files under sub directories too. This is always true when searching with phrase. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.ListArchiveFiles(context.Background(), bucket, filepath).Name(name).Phrase(phrase).FullTree(fullTree).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.ListArchiveFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListArchiveFiles`: ListArchiveFiles200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.ListArchiveFiles`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**bucket** | **string** | Bucket name | 
**filepath** | **string** | The path to to search for files under. Default is the root directory /. | [default to &quot;/&quot;]

### Other Parameters

Other parameters are passed through a pointer to a apiListArchiveFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **name** | **string** | If specified will return an exact match on file name | 
 **phrase** | **string** | Search phrase for partial matches on file name, wildcard may be specified as %, eg. example-% This also searches for files under sub directories too.  | 
 **fullTree** | **bool** | Include files under sub directories too. This is always true when searching with phrase. | [default to false]

### Return type

[**ListArchiveFiles200Response**](ListArchiveFiles200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateArchiveBucket

> UpdateArchiveBucket200Response UpdateArchiveBucket(ctx, id).UpdateArchiveBucketRequest(updateArchiveBucketRequest).Execute()

Update an Archive Bucket



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
	updateArchiveBucketRequest := *openapiclient.NewUpdateArchiveBucketRequest() // UpdateArchiveBucketRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArchivesAPI.UpdateArchiveBucket(context.Background(), id).UpdateArchiveBucketRequest(updateArchiveBucketRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArchivesAPI.UpdateArchiveBucket``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateArchiveBucket`: UpdateArchiveBucket200Response
	fmt.Fprintf(os.Stdout, "Response from `ArchivesAPI.UpdateArchiveBucket`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64** | Morpheus ID of the Object being referenced | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateArchiveBucketRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateArchiveBucketRequest** | [**UpdateArchiveBucketRequest**](UpdateArchiveBucketRequest.md) |  | 

### Return type

[**UpdateArchiveBucket200Response**](UpdateArchiveBucket200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

