# \SetupAPI

All URIs are relative to *https://CHANGEME*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Setup**](SetupAPI.md#Setup) | **Post** /api/setup | Setup appliance



## Setup

> Setup200Response Setup(ctx).SetupRequest(setupRequest).Execute()

Setup appliance



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
	setupRequest := *openapiclient.NewSetupRequest("ApplianceName_example", "AccountName_example", "Username_example", "Email_example", "Password_example") // SetupRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SetupAPI.Setup(context.Background()).SetupRequest(setupRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SetupAPI.Setup``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Setup`: Setup200Response
	fmt.Fprintf(os.Stdout, "Response from `SetupAPI.Setup`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSetupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **setupRequest** | [**SetupRequest**](SetupRequest.md) |  | 

### Return type

[**Setup200Response**](Setup200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

