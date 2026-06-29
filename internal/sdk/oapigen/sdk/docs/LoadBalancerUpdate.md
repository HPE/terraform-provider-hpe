# LoadBalancerUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Enabled** | Pointer to **bool** | Activate (true) or disable (false) | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object with parameters that vary by load balancer type. | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "public"]
**Tenants** | Pointer to [**[]LoadBalancerUpdateTenantsInner**](LoadBalancerUpdateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermission** | Pointer to [**LoadBalancerUpdateResourcePermission**](LoadBalancerUpdateResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancerUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


