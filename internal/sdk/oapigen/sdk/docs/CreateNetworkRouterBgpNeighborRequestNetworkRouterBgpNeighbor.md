# CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IpAddress** | **string** | IP address of the BGP neighbor | 
**Description** | Pointer to **string** | Description of the BGP neighbor | [optional] 
**ForwardingAddress** | Pointer to **string** | Forwarding address (NSX-V distributed router) | [optional] 
**ProtocolAddress** | Pointer to **string** | Protocol address (NSX-V distributed router) | [optional] 
**RemoteAs** | Pointer to **string** | Remote AS number | [optional] 
**Weight** | Pointer to **int64** | Route weight | [optional] [default to 60]
**KeepAlive** | Pointer to **int64** | Keep-alive interval in seconds | [optional] [default to 60]
**HoldDown** | Pointer to **int64** | Hold-down timer in seconds | [optional] [default to 180]
**Password** | Pointer to **string** | BGP session password | [optional] 
**RouteFilteringType** | Pointer to **string** | Address family for route filtering (e.g. IPV4, IPV6) | [optional] 
**RouteFilteringIn** | Pointer to **string** | Inbound route filter name | [optional] 
**RouteFilteringOut** | Pointer to **string** | Outbound route filter name | [optional] 
**BfdEnabled** | Pointer to [**CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled**](CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborBfdEnabled.md) |  | [optional] 
**BfdInterval** | Pointer to **int64** | BFD interval in milliseconds | [optional] 
**BfdMultiple** | Pointer to **int64** | BFD multiplier | [optional] 
**AllowAsIn** | Pointer to [**CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn**](CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborAllowAsIn.md) |  | [optional] 
**HopLimit** | Pointer to **int64** | Maximum hop limit | [optional] [default to 1]
**RestartMode** | Pointer to **string** | Graceful restart mode (e.g. HELPER_ONLY, GRACEFUL_RESTART, DISABLE) | [optional] 
**Config** | Pointer to [**CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig**](CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighborConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


