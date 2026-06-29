# AddClusterNamespaceRequestNamespace

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Namespace name | 
**Description** | Pointer to **string** | Namespace description | [optional] 
**Active** | Pointer to **bool** | Namespace active | [optional] [default to false]
**ResourcePermissions** | Pointer to [**AddClusterNamespaceRequestNamespaceResourcePermissions**](AddClusterNamespaceRequestNamespaceResourcePermissions.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddClusterNamespaceRequestNamespace{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


