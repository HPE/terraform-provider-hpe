# CreateLoadBalancerPoolNodeRequestLoadBalancerNode

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**IpAddress** | Pointer to **string** | IP Address | [optional] 
**Port** | Pointer to **int32** | Port number | [optional] 
**Weight** | Pointer to **int32** | Weight applied balance algoritm | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object with parameters that vary by type. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerPoolNodeRequestLoadBalancerNode{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


