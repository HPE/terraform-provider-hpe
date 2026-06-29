# NetworkRoutersUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Type** | Pointer to [**NetworkRoutersUpdateType**](NetworkRoutersUpdateType.md) |  | [optional] 
**Site** | Pointer to [**NetworkRoutersUpdateSite**](NetworkRoutersUpdateSite.md) |  | [optional] 
**Enabled** | Pointer to **bool** | Can be used to enable / disable the network router (true, false). Default is on | [optional] 
**Zone** | Pointer to [**NetworkRoutersUpdateZone**](NetworkRoutersUpdateZone.md) |  | [optional] 
**NetworkServer** | Pointer to [**NetworkRoutersUpdateNetworkServer**](NetworkRoutersUpdateNetworkServer.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NetworkRoutersUpdate{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


