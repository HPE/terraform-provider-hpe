# UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SourceAddresses** | Pointer to **[]string** | Source IP addresses for the BGP session | [optional] 
**RouterId** | Pointer to **string** | The router identifier (auto-populated for edge routers) | [optional] 
**Interface** | Pointer to **string** | The interface name for the BGP session (distributed routers only) | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


