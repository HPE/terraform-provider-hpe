# UpdateCloudResourcePoolRequestResourcePool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** | Activate &#x60;true&#x60; or disable &#x60;false&#x60; the datastore | [optional] 
**Visibility** | Pointer to **string** | Setting &#x60;private&#x60; or &#x60;public&#x60; | [optional] [default to "private"]
**DisplayName** | Pointer to **string** | Optional Display Name (VMware only) | [optional] 
**Inventory** | Pointer to **bool** | Enable &#x60;True&#x60; or disable &#x60;False&#x60; inventory sync for resource pool during cloud refresh | [optional] 
**TenantPermissions** | Pointer to [**UpdateCloudResourcePoolRequestResourcePoolTenantPermissions**](UpdateCloudResourcePoolRequestResourcePoolTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**UpdateCloudResourcePoolRequestResourcePoolResourcePermissions**](UpdateCloudResourcePoolRequestResourcePoolResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateCloudResourcePoolRequestResourcePool{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


