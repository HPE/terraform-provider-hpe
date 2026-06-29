# \CypherAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddCypherKey**](CypherAPI.md#AddCypherKey) | **Post** /api/cypher/{cypherPath} | Write a Cypher
[**GetCypherKey**](CypherAPI.md#GetCypherKey) | **Get** /api/cypher/{cypherPath} | Read or Create a Cypher Key
[**ListCypherKeys**](CypherAPI.md#ListCypherKeys) | **Get** /api/cypher | List Cypher Keys
[**RemoveCypher**](CypherAPI.md#RemoveCypher) | **Delete** /api/cypher/{cypherPath} | Delete a Cypher



## AddCypherKey

> AddCypherKey200Response AddCypherKey(ctx, cypherPath).Ttl(ttl).Value(value).Type_(type_).AddCypherKeyRequest(addCypherKeyRequest).Execute()

Write a Cypher



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
	cypherPath := "cypherPath_example" // string | The key includes a mount prefix separated by a /. For example, the key secret/foo uses the secret mount.  Available Mounts  <table>   <tr>     <th>Mount</th>     <th>Description</th>     <th>Example</th>   </tr>   <tr>     <td>password</td>     <td>Generates a secure password of specified character length in the key pattern (or 15) with symbols, numbers, upper case, and lower case letters (i.e. password/15/mypass generates a 15 character password).</td>     <td>password/15/mypass</td>   </tr>   <tr>     <td>tfvars</td>     <td>This is a module to store a tfvars file for terraform.</td>     <td>tfvars/mytfvar</td>   </tr>   <tr>     <td>secret</td>     <td>This is the standard secret module that stores a key/value in encrypted form. Capable of storing entire JSON object or a String.</td>     <td>secret/foo</td>   </tr>   <tr>     <td>uuid</td>     <td>Returns a new UUID by key name when requested and stores the generated UUID by key name for a given lease timeout period.</td>     <td>uuid/autoMac1</td>   </tr>   <tr>     <td>key</td>     <td>Generates a Base 64 encoded AES Key of specified bit length in the key pattern (i.e. key/128/mykey generates a 128-bit key)</td>     <td>key/128/mykey</td>   </tr> </table> 
	ttl := openapiclient.addCypherKey_ttl_parameter{Int64: new(int64)} // AddCypherKeyTtlParameter | Time to Live. The lease duration in seconds, or a human readable format eg. '15m', 8h, '7d'.  0 means no expiry. (optional)
	value := "{"foo":"bar", "zip":"zap"}" // string | The secret value to be stored. Only required for certain mounts. Some mounts generate their own value and do not require a value to be passed. eg. `uuid`, `key` and `password`. (optional)
	type_ := "type__example" // string | The type of data being stored, `string` or `object`. The data type depends on the cypher mount being used. Most mounts use `string` as their data type, but `secret` uses `object` by default. You can store a string instead by passing `type=string`. This means the `data` value returned by the API will be a string instead of an object. (optional)
	addCypherKeyRequest := *openapiclient.NewAddCypherKeyRequest() // AddCypherKeyRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CypherAPI.AddCypherKey(context.Background(), cypherPath).Ttl(ttl).Value(value).Type_(type_).AddCypherKeyRequest(addCypherKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CypherAPI.AddCypherKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddCypherKey`: AddCypherKey200Response
	fmt.Fprintf(os.Stdout, "Response from `CypherAPI.AddCypherKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cypherPath** | **string** | The key includes a mount prefix separated by a /. For example, the key secret/foo uses the secret mount.  Available Mounts  &lt;table&gt;   &lt;tr&gt;     &lt;th&gt;Mount&lt;/th&gt;     &lt;th&gt;Description&lt;/th&gt;     &lt;th&gt;Example&lt;/th&gt;   &lt;/tr&gt;   &lt;tr&gt;     &lt;td&gt;password&lt;/td&gt;     &lt;td&gt;Generates a secure password of specified character length in the key pattern (or 15) with symbols, numbers, upper case, and lower case letters (i.e. password/15/mypass generates a 15 character password).&lt;/td&gt;     &lt;td&gt;password/15/mypass&lt;/td&gt;   &lt;/tr&gt;   &lt;tr&gt;     &lt;td&gt;tfvars&lt;/td&gt;     &lt;td&gt;This is a module to store a tfvars file for terraform.&lt;/td&gt;     &lt;td&gt;tfvars/mytfvar&lt;/td&gt;   &lt;/tr&gt;   &lt;tr&gt;     &lt;td&gt;secret&lt;/td&gt;     &lt;td&gt;This is the standard secret module that stores a key/value in encrypted form. Capable of storing entire JSON object or a String.&lt;/td&gt;     &lt;td&gt;secret/foo&lt;/td&gt;   &lt;/tr&gt;   &lt;tr&gt;     &lt;td&gt;uuid&lt;/td&gt;     &lt;td&gt;Returns a new UUID by key name when requested and stores the generated UUID by key name for a given lease timeout period.&lt;/td&gt;     &lt;td&gt;uuid/autoMac1&lt;/td&gt;   &lt;/tr&gt;   &lt;tr&gt;     &lt;td&gt;key&lt;/td&gt;     &lt;td&gt;Generates a Base 64 encoded AES Key of specified bit length in the key pattern (i.e. key/128/mykey generates a 128-bit key)&lt;/td&gt;     &lt;td&gt;key/128/mykey&lt;/td&gt;   &lt;/tr&gt; &lt;/table&gt;  | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddCypherKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ttl** | [**AddCypherKeyTtlParameter**](AddCypherKeyTtlParameter.md) | Time to Live. The lease duration in seconds, or a human readable format eg. &#39;15m&#39;, 8h, &#39;7d&#39;.  0 means no expiry. | 
 **value** | **string** | The secret value to be stored. Only required for certain mounts. Some mounts generate their own value and do not require a value to be passed. eg. &#x60;uuid&#x60;, &#x60;key&#x60; and &#x60;password&#x60;. | 
 **type_** | **string** | The type of data being stored, &#x60;string&#x60; or &#x60;object&#x60;. The data type depends on the cypher mount being used. Most mounts use &#x60;string&#x60; as their data type, but &#x60;secret&#x60; uses &#x60;object&#x60; by default. You can store a string instead by passing &#x60;type&#x3D;string&#x60;. This means the &#x60;data&#x60; value returned by the API will be a string instead of an object. | 
 **addCypherKeyRequest** | [**AddCypherKeyRequest**](AddCypherKeyRequest.md) |  | 

### Return type

[**AddCypherKey200Response**](AddCypherKey200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCypherKey

> GetCypherKey200Response GetCypherKey(ctx, cypherPath).LeaseToken(leaseToken).Sort(sort).Direction(direction).Execute()

Read or Create a Cypher Key



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
	cypherPath := "password/15/mypass" // string | The cypher key including the mount prefix.
	leaseToken := "5bd808dc-82bb-4974-b699-5f1cdd3cc503" // string | An execution lease token. (optional)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CypherAPI.GetCypherKey(context.Background(), cypherPath).LeaseToken(leaseToken).Sort(sort).Direction(direction).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CypherAPI.GetCypherKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCypherKey`: GetCypherKey200Response
	fmt.Fprintf(os.Stdout, "Response from `CypherAPI.GetCypherKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cypherPath** | **string** | The cypher key including the mount prefix. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetCypherKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **leaseToken** | **string** | An execution lease token. | 
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]

### Return type

[**GetCypherKey200Response**](GetCypherKey200Response.md)

### Authorization

[cypherAuth-XVToken](../README.md#cypherAuth-XVToken), [cypherAuth-XMLease](../README.md#cypherAuth-XMLease), [cypherAuth-XCToken](../README.md#cypherAuth-XCToken), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListCypherKeys

> ListCypherKeys200Response ListCypherKeys(ctx).LeaseToken(leaseToken).List(list).Key(key).Phrase(phrase).Max(max).Offset(offset).Sort(sort).Direction(direction).Execute()

List Cypher Keys



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
	leaseToken := "5bd808dc-82bb-4974-b699-5f1cdd3cc503" // string | An execution lease token. (optional)
	list := true // bool | This endpoint is available via the http method LIST. The GET method can be used to list keys as well, by passing the query parameter list=true. (optional) (default to false)
	key := "secret/test" // string | If specified will return Cypher items with the matching Key value (optional)
	phrase := "phrase_example" // string | Search phrase for partial matches on name or description (optional)
	max := int64(789) // int64 | Maximum number of records to return (optional) (default to 25)
	offset := int64(789) // int64 | Offset records, the number of records to skip, for paginating requests (optional) (default to 0)
	sort := "sort_example" // string | Sort order, the name of the property to sort by (optional) (default to "name")
	direction := "asc" // string | Sort direction, use 'desc' to reverse sort (optional) (default to "asc")

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CypherAPI.ListCypherKeys(context.Background()).LeaseToken(leaseToken).List(list).Key(key).Phrase(phrase).Max(max).Offset(offset).Sort(sort).Direction(direction).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CypherAPI.ListCypherKeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListCypherKeys`: ListCypherKeys200Response
	fmt.Fprintf(os.Stdout, "Response from `CypherAPI.ListCypherKeys`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListCypherKeysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **leaseToken** | **string** | An execution lease token. | 
 **list** | **bool** | This endpoint is available via the http method LIST. The GET method can be used to list keys as well, by passing the query parameter list&#x3D;true. | [default to false]
 **key** | **string** | If specified will return Cypher items with the matching Key value | 
 **phrase** | **string** | Search phrase for partial matches on name or description | 
 **max** | **int64** | Maximum number of records to return | [default to 25]
 **offset** | **int64** | Offset records, the number of records to skip, for paginating requests | [default to 0]
 **sort** | **string** | Sort order, the name of the property to sort by | [default to &quot;name&quot;]
 **direction** | **string** | Sort direction, use &#39;desc&#39; to reverse sort | [default to &quot;asc&quot;]

### Return type

[**ListCypherKeys200Response**](ListCypherKeys200Response.md)

### Authorization

[cypherAuth-XVToken](../README.md#cypherAuth-XVToken), [cypherAuth-XMLease](../README.md#cypherAuth-XMLease), [cypherAuth-XCToken](../README.md#cypherAuth-XCToken), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveCypher

> RemoveExecuteSchedules200Response RemoveCypher(ctx, cypherPath).Execute()

Delete a Cypher



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
	cypherPath := "password/15/mypass" // string | The cypher key including the mount prefix.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CypherAPI.RemoveCypher(context.Background(), cypherPath).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CypherAPI.RemoveCypher``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveCypher`: RemoveExecuteSchedules200Response
	fmt.Fprintf(os.Stdout, "Response from `CypherAPI.RemoveCypher`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**cypherPath** | **string** | The cypher key including the mount prefix. | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveCypherRequest struct via the builder pattern


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

