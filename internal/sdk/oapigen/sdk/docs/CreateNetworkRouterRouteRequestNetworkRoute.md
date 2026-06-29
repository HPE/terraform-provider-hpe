# CreateNetworkRouterRouteRequestNetworkRoute

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Route name | [optional] 
**Description** | Pointer to **string** | Route description | [optional] 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the route (true, false). Default is off | [optional] [default to false]
**DefaultRoute** | Pointer to **bool** | Can be used to set as default route (true, false). Default is off | [optional] [default to false]
**Source** | **string** | Source address or range | 
**Destination** | **string** | Destination address or range | 
**NetworkMtu** | Pointer to **float32** | MTU | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkRouterRouteRequestNetworkRoute{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


