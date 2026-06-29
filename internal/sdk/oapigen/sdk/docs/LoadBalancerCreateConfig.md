# LoadBalancerCreateConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Plan** | Pointer to [**HAProxyLoadBalancerConfigObject1Plan**](HAProxyLoadBalancerConfigObject1Plan.md) |  | [optional] 
**Pool** | Pointer to [**HAProxyLoadBalancerConfigObject1Pool**](HAProxyLoadBalancerConfigObject1Pool.md) |  | [optional] 
**AdminState** | Pointer to **bool** | If true, the admin state is enabled. | [optional] 
**Size** | Pointer to **string** | Load balancer size. | [optional] 
**Tier1** | Pointer to **string** | Tier 1 gateway provider ID. | [optional] 
**Loglevel** | Pointer to **string** | Log level. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancerCreateConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


