# UpdateDatastoresRequestDatastore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** | Activate &#x60;true&#x60; or disable &#x60;false&#x60; the datastore | [optional] 
**HeartbeatTarget** | Pointer to **bool** | Heartbeat Target | [optional] 
**Visibility** | Pointer to **string** | Setting &#x60;private&#x60; or &#x60;public&#x60; | [optional] [default to "private"]
**TenantPermissions** | Pointer to [**UpdateDatastoresRequestDatastoreTenantPermissions**](UpdateDatastoresRequestDatastoreTenantPermissions.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**UpdateDatastoresRequestDatastoreResourcePermissions**](UpdateDatastoresRequestDatastoreResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateDatastoresRequestDatastore{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


