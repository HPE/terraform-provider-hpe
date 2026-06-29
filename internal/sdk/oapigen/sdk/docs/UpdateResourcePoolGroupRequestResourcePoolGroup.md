# UpdateResourcePoolGroupRequestResourcePoolGroup

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Mode** | Pointer to **string** | Pool selection mode. Valid values are &#x60;roundrobin&#x60; or &#x60;availablecapacity&#x60;. | [optional] 
**Pools** | Pointer to **[]int64** |  | [optional] 
**Tenants** | Pointer to [**[]UpdateResourcePoolGroupRequestResourcePoolGroupTenantsInner**](UpdateResourcePoolGroupRequestResourcePoolGroupTenantsInner.md) |  | [optional] 
**ResourcePermission** | Pointer to [**UpdateResourcePoolGroupRequestResourcePoolGroupResourcePermission**](UpdateResourcePoolGroupRequestResourcePoolGroupResourcePermission.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateResourcePoolGroupRequestResourcePoolGroup{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


