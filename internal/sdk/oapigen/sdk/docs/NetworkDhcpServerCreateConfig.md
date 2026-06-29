# NetworkDhcpServerCreateConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EdgeCluster** | Pointer to **string** | Edge Cluster | [optional] 
**PreferredEdgeNode1** | Pointer to **string** | Active Edge Node Options obtained by calling option source with :optionSource &#x3D; nsxtEdgeNodes and networkServerId param | [optional] 
**PreferredEdgeNode2** | Pointer to **string** | Standby Edge Node Options obtained by calling option source with optionSource &#x3D; nsxtEdgeNodes and networkServerId param | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkDhcpServerCreateConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


