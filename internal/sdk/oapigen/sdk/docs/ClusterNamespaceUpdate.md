# ClusterNamespaceUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Namespace name | [optional] 
**Description** | Pointer to **string** | Namespace description | [optional] 
**Active** | Pointer to **bool** | Namespace active | [optional] [default to false]
**Permissions** | Pointer to [**ClusterNamespaceUpdatePermissions**](ClusterNamespaceUpdatePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterNamespaceUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


