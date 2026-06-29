# CreateNetworkRouterRequestNetworkRouter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name | 
**Type** | [**CreateNetworkRouterRequestNetworkRouterType**](CreateNetworkRouterRequestNetworkRouterType.md) |  | 
**Site** | [**CreateNetworkRouterRequestNetworkRouterSite**](CreateNetworkRouterRequestNetworkRouterSite.md) |  | 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the network router (true, false). Default is on | [optional] 
**EnableBgp** | Pointer to **bool** |  | [optional] 
**Zone** | Pointer to [**CreateNetworkRouterRequestNetworkRouterZone**](CreateNetworkRouterRequestNetworkRouterZone.md) |  | [optional] 
**NetworkServer** | Pointer to [**CreateNetworkRouterRequestNetworkRouterNetworkServer**](CreateNetworkRouterRequestNetworkRouterNetworkServer.md) |  | [optional] 
**Config** | Pointer to [**CreateNetworkRouterRequestNetworkRouterConfig**](CreateNetworkRouterRequestNetworkRouterConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkRouterRequestNetworkRouter{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


