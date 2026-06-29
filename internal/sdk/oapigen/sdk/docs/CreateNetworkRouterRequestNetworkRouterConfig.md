# CreateNetworkRouterRequestNetworkRouterConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**HaMode** | Pointer to **string** |  | [optional] 
**EdgeCluster** | Pointer to **string** |  | [optional] 
**FailOver** | Pointer to **string** |  | [optional] 
**IpManagementType** | Pointer to **string** |  | [optional] 
**IpServerId** | Pointer to **string** |  | [optional] 
**TIER0STATIC** | Pointer to **string** |  | [optional] 
**TIER0NAT** | Pointer to **string** |  | [optional] 
**TIER0IPSECLOCALIP** | Pointer to **string** |  | [optional] 
**TIER0DNSFORWARDERIP** | Pointer to **string** |  | [optional] 
**TIER0SERVICEINTERFACE** | Pointer to **string** |  | [optional] 
**TIER0EXTERNALINTERFACE** | Pointer to **string** |  | [optional] 
**TIER0LOOPBACKINTERFACE** | Pointer to **string** |  | [optional] 
**TIER0SEGMENT** | Pointer to **string** |  | [optional] 
**TIER1DNSFORWARDERIP** | Pointer to **string** |  | [optional] 
**TIER1STATIC** | Pointer to **string** |  | [optional] 
**TIER1LBVIP** | Pointer to **string** |  | [optional] 
**TIER1NAT** | Pointer to **string** |  | [optional] 
**TIER1LBSNAT** | Pointer to **string** |  | [optional] 
**TIER1IPSECLOCALENDPOINT** | Pointer to **string** |  | [optional] 
**TIER1SERVICEINTERFACE** | Pointer to **string** |  | [optional] 
**TIER1SEGMENT** | Pointer to **string** |  | [optional] 
**LOCAL_AS_NUM** | Pointer to **string** |  | [optional] 
**ECMP** | Pointer to **string** |  | [optional] 
**MULTIPATH_RELAX** | Pointer to **string** |  | [optional] 
**RESTART_MODE** | Pointer to **string** |  | [optional] 
**RESTART_TIME** | Pointer to **int64** |  | [optional] 
**STALE_ROUTE_TIME** | Pointer to **int64** |  | [optional] 
**INTER_SR_IBGP** | Pointer to **string** |  | [optional] 
**Tier0Gateway** | Pointer to **string** |  | [optional] 
**TIER1CONNECTED** | Pointer to **string** |  | [optional] 
**TIER1STATICROUTES** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkRouterRequestNetworkRouterConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


