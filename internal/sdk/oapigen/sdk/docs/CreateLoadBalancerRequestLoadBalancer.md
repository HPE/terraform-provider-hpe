# CreateLoadBalancerRequestLoadBalancer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** | Load Balancer Type Code | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**NetworkServerId** | Pointer to **int64** | Network Server ID | [optional] 
**Site** | Pointer to [**CreateLoadBalancerRequestLoadBalancerSite**](CreateLoadBalancerRequestLoadBalancerSite.md) |  | [optional] 
**Zone** | Pointer to [**CreateLoadBalancerRequestLoadBalancerZone**](CreateLoadBalancerRequestLoadBalancerZone.md) |  | [optional] 
**Config** | Pointer to [**CreateLoadBalancerRequestLoadBalancerConfig**](CreateLoadBalancerRequestLoadBalancerConfig.md) |  | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "public"]
**Tenants** | Pointer to [**[]CreateLoadBalancerRequestLoadBalancerTenantsInner**](CreateLoadBalancerRequestLoadBalancerTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**CreateLoadBalancerRequestLoadBalancerResourcePermissions**](CreateLoadBalancerRequestLoadBalancerResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerRequestLoadBalancer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


