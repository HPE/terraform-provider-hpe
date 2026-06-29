# NetworkRoutersCreate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name | 
**Type** | [**NetworkRoutersCreateType**](NetworkRoutersCreateType.md) |  | 
**Site** | [**NetworkRoutersCreateSite**](NetworkRoutersCreateSite.md) |  | 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the network router (true, false). Default is on | [optional] 
**EnableBgp** | Pointer to **bool** |  | [optional] 
**Zone** | Pointer to [**NetworkRoutersCreateZone**](NetworkRoutersCreateZone.md) |  | [optional] 
**NetworkServer** | Pointer to [**NetworkRoutersCreateNetworkServer**](NetworkRoutersCreateNetworkServer.md) |  | [optional] 
**Config** | Pointer to [**NetworkRoutersCreateConfig**](NetworkRoutersCreateConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkRoutersCreate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


