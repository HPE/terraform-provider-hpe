# ZoneDatastoreCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**DatastoreType** | Pointer to **string** | The code of the datastoreType | [optional] 
**StorageServer** | Pointer to [**ZoneDatastoreCreateStorageServer**](ZoneDatastoreCreateStorageServer.md) |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**DefaultStore** | Pointer to **bool** |  | [optional] 
**Tenants** | Pointer to [**[]ZoneDatastoreCreateTenantsInner**](ZoneDatastoreCreateTenantsInner.md) |  | [optional] 
**ResourcePermissions** | Pointer to [**ZoneDatastoreCreateResourcePermissions**](ZoneDatastoreCreateResourcePermissions.md) |  | [optional] 
**Datastores** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ZoneDatastoreCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


