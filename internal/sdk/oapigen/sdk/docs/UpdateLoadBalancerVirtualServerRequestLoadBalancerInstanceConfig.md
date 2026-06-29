# UpdateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ApplicationProfile** | Pointer to **int64** | The Load Balancer Application Profile ID. The Options API &#x60;/api/options/nsxt/nsxtLBVirtualServerApplicationProfile?loadBalancerId&#x3D;42&amp;loadBalancerInstance.vipProtocol&#x3D;tcp&#x60; can be used to see which options are available. | [optional] 
**Pool** | Pointer to **string** | The backend server pool ID (&#x60;NetworkLoadBalancerPool&#x60;). The Options API &#x60;/api/options/nsxt/nsxtLBPool?loadBalancerId&#x3D;42&#x60; can be used to see which options are available. | [optional] 
**Persistence** | Pointer to **NullableString** | Session persistence mode. The available values depend on the virtual server protocol. For HTTP: &#x60;SOURCE_IP&#x60;, &#x60;COOKIE&#x60;, or empty string (disabled). For TCP/UDP: &#x60;SOURCE_IP&#x60; or empty string (disabled). The Options API &#x60;/api/options/nsxt/nsxtLBPersistence?loadBalancerId&#x3D;42&amp;loadBalancerInstance.vipProtocol&#x3D;tcp&#x60; can be used to see which options are available. | [optional] 
**PersistenceProfile** | Pointer to **int64** | The ID of the persistence profile to use. Required when &#x60;persistence&#x60; is set to a non-empty value (&#x60;SOURCE_IP&#x60; or &#x60;COOKIE&#x60;). The Options API &#x60;/api/options/nsxt/nsxtLBPersistenceProfile?loadBalancerId&#x3D;42&amp;config.persistence&#x3D;SOURCE_IP&#x60; can be used to see which options are available. | [optional] 
**SslClientProfile** | Pointer to **int64** | The SSL client profile ID. Only applicable when &#x60;sslCert&#x60; is set to a non-zero value. The Options API &#x60;/api/options/nsxt/nsxtLBClientSSlProfiles?loadBalancerId&#x3D;42&#x60; can be used to see which options are available. | [optional] 
**SslServerProfile** | Pointer to **int64** | The SSL server profile ID. Only applicable when &#x60;sslServerCert&#x60; is set to a non-zero value. The Options API &#x60;/api/options/nsxt/nsxtLBServerSSlProfiles?loadBalancerId&#x3D;42&#x60; can be used to see which options are available. | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateLoadBalancerVirtualServerRequestLoadBalancerInstanceConfig{
    // Set fields directly
}
```

### Persistence (Nullable)

Use the Nullable wrapper methods:
- `obj.Persistence.IsSet()` — check if set
- `obj.Persistence.Get()` — get the inner value (returns pointer)
- `obj.Persistence.Set(&val)` — set the value
- `obj.Persistence.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


