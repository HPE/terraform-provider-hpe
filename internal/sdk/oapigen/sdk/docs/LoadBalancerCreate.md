# LoadBalancerCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** | Load Balancer Type Code | [optional] 
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**NetworkServerId** | Pointer to **int64** | Network Server ID | [optional] 
**Site** | Pointer to [**LoadBalancerCreateSite**](LoadBalancerCreateSite.md) |  | [optional] 
**Zone** | Pointer to [**LoadBalancerCreateZone**](LoadBalancerCreateZone.md) |  | [optional] 
**Config** | Pointer to [**LoadBalancerCreateConfig**](LoadBalancerCreateConfig.md) |  | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "public"]
**Tenants** | Pointer to [**[]LoadBalancerCreateTenantsInner**](LoadBalancerCreateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 
**ResourcePermissions** | Pointer to [**LoadBalancerCreateResourcePermissions**](LoadBalancerCreateResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &LoadBalancerCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


