# AddCloudResourcePoolRequestResourcePool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of Resource Pool | 
**DefaultPool** | Pointer to **bool** | Set as the Default Pool | [optional] [default to false]
**DefaultImage** | Pointer to **bool** | Set as the Default Image Target | [optional] [default to false]
**Active** | Pointer to **bool** | Activate &#x60;true&#x60; or disable &#x60;false&#x60; the datastore | [optional] [default to true]
**Visibility** | Pointer to **string** | Setting &#x60;private&#x60; or &#x60;public&#x60; | [optional] [default to "private"]
**DisplayName** | Pointer to **string** | Optional Display Name (VMware only) | [optional] 
**Inventory** | Pointer to **bool** | Enable &#x60;True&#x60; or disable &#x60;False&#x60; inventory sync for resource pool during cloud refresh | [optional] [default to true]
**Config** | [**AddCloudResourcePoolRequestResourcePoolConfig**](AddCloudResourcePoolRequestResourcePoolConfig.md) |  | 
**TenantPermissions** | Pointer to [**AddCloudResourcePoolRequestResourcePoolTenantPermissions**](AddCloudResourcePoolRequestResourcePoolTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**AddCloudResourcePoolRequestResourcePoolResourcePermissions**](AddCloudResourcePoolRequestResourcePoolResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddCloudResourcePoolRequestResourcePool{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


