# UpdateLoadBalancerPoolRequestLoadBalancerPool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**VipBalance** | Pointer to **string** | Balance Algorithm | [optional] 
**MinActive** | Pointer to **int64** | Min Active Members | [optional] 
**Port** | Pointer to **int64** | Port number | [optional] 
**VipSticky** | Pointer to **string** | Session persistence mode | [optional] 
**VipClientIpMode** | Pointer to **string** | VIP client IP mode | [optional] 
**Partition** | Pointer to **string** | Partition | [optional] 
**Config** | Pointer to [**UpdateLoadBalancerPoolRequestLoadBalancerPoolConfig**](UpdateLoadBalancerPoolRequestLoadBalancerPoolConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateLoadBalancerPoolRequestLoadBalancerPool{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


