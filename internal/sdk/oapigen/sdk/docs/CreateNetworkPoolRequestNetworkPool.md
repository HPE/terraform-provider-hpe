# CreateNetworkPoolRequestNetworkPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Type** | Pointer to [**CreateNetworkPoolRequestNetworkPoolType**](CreateNetworkPoolRequestNetworkPoolType.md) |  | [optional] 
**IpRanges** | Pointer to [**[]CreateNetworkPoolRequestNetworkPoolIpRangesInner**](CreateNetworkPoolRequestNetworkPoolIpRangesInner.md) | Array of IP range objects. Type &#39;morpheus&#39; expects startAddress and endAddress. Type &#39;morpheusipv6&#39; expects a cidrIPv6. | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object with parameters that vary by pool type. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkPoolRequestNetworkPool{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


