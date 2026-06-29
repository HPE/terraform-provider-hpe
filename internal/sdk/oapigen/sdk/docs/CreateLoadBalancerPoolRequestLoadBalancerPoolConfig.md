# CreateLoadBalancerPoolRequestLoadBalancerPoolConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActiveMonitorPaths** | Pointer to **int64** | The ID of the active health monitor (NetworkLoadBalancerMonitor). The Options API &#x60;/api/options/nsxt/nsxtLBPoolActiveMonitor?loadBalancerId&#x3D;{id}&#x60; can be used to see which options are available. | [optional] 
**PassiveMonitorPath** | Pointer to **int64** | The ID of the passive health monitor (NetworkLoadBalancerMonitor). The Options API &#x60;/api/options/nsxt/nsxtLBPoolPassiveMonitor?loadBalancerId&#x3D;{id}&#x60; can be used to see which options are available. | [optional] 
**SnatTranslationType** | Pointer to **string** | SNAT translation type. Determines how source NAT is applied to pool traffic. | [optional] 
**SnatIpAddresses** | Pointer to **[]string** | List of SNAT IP addresses. Required when snatTranslationType is LBSnatIpPool. | [optional] 
**TcpMultiplexing** | Pointer to **bool** | Whether TCP multiplexing is enabled for the pool. | [optional] 
**TcpMultiplexingNumber** | Pointer to **int64** | Maximum number of TCP multiplexing connections. Defaults to 6. | [optional] 
**MemberGroup** | Pointer to [**NSXTLoadBalancerPoolConfigObjectMemberGroup**](NSXTLoadBalancerPoolConfigObjectMemberGroup.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerPoolRequestLoadBalancerPoolConfig{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


