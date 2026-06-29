# UpdateLoadBalancerMonitorRequestLoadBalancerMonitor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**MonitorType** | Pointer to **string** |  | [optional] 
**MonitorInterval** | Pointer to **int64** |  | [optional] 
**MonitorTimeout** | Pointer to **int64** |  | [optional] 
**SendData** | Pointer to **NullableString** |  | [optional] 
**SendVersion** | Pointer to **NullableString** |  | [optional] 
**SendType** | Pointer to **NullableString** |  | [optional] 
**ReceiveData** | Pointer to **NullableString** |  | [optional] 
**ReceiveCode** | Pointer to **NullableString** |  | [optional] 
**MonitorUsername** | Pointer to **NullableString** |  | [optional] 
**MonitorPassword** | Pointer to **NullableString** |  | [optional] 
**MonitorDestination** | Pointer to **NullableString** |  | [optional] 
**FallCount** | Pointer to **int64** |  | [optional] 
**RiseCount** | Pointer to **int64** |  | [optional] 
**AliasPort** | Pointer to **int64** |  | [optional] 
**DataLength** | Pointer to **int64** |  | [optional] 
**MaxRetry** | Pointer to **int64** |  | [optional] 
**ExtraConfig** | Pointer to **NullableString** |  | [optional] 
**Partition** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to [**UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig**](UpdateLoadBalancerMonitorRequestLoadBalancerMonitorConfig.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateLoadBalancerMonitorRequestLoadBalancerMonitor{
    // Set fields directly
}
```

### SendData (Nullable)

Use the Nullable wrapper methods:
- `obj.SendData.IsSet()` — check if set
- `obj.SendData.Get()` — get the inner value (returns pointer)
- `obj.SendData.Set(&val)` — set the value
- `obj.SendData.Unset()` — clear the value
### SendVersion (Nullable)

Use the Nullable wrapper methods:
- `obj.SendVersion.IsSet()` — check if set
- `obj.SendVersion.Get()` — get the inner value (returns pointer)
- `obj.SendVersion.Set(&val)` — set the value
- `obj.SendVersion.Unset()` — clear the value
### SendType (Nullable)

Use the Nullable wrapper methods:
- `obj.SendType.IsSet()` — check if set
- `obj.SendType.Get()` — get the inner value (returns pointer)
- `obj.SendType.Set(&val)` — set the value
- `obj.SendType.Unset()` — clear the value
### ReceiveData (Nullable)

Use the Nullable wrapper methods:
- `obj.ReceiveData.IsSet()` — check if set
- `obj.ReceiveData.Get()` — get the inner value (returns pointer)
- `obj.ReceiveData.Set(&val)` — set the value
- `obj.ReceiveData.Unset()` — clear the value
### ReceiveCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ReceiveCode.IsSet()` — check if set
- `obj.ReceiveCode.Get()` — get the inner value (returns pointer)
- `obj.ReceiveCode.Set(&val)` — set the value
- `obj.ReceiveCode.Unset()` — clear the value
### MonitorUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorUsername.IsSet()` — check if set
- `obj.MonitorUsername.Get()` — get the inner value (returns pointer)
- `obj.MonitorUsername.Set(&val)` — set the value
- `obj.MonitorUsername.Unset()` — clear the value
### MonitorPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorPassword.IsSet()` — check if set
- `obj.MonitorPassword.Get()` — get the inner value (returns pointer)
- `obj.MonitorPassword.Set(&val)` — set the value
- `obj.MonitorPassword.Unset()` — clear the value
### MonitorDestination (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorDestination.IsSet()` — check if set
- `obj.MonitorDestination.Get()` — get the inner value (returns pointer)
- `obj.MonitorDestination.Set(&val)` — set the value
- `obj.MonitorDestination.Unset()` — clear the value
### ExtraConfig (Nullable)

Use the Nullable wrapper methods:
- `obj.ExtraConfig.IsSet()` — check if set
- `obj.ExtraConfig.Get()` — get the inner value (returns pointer)
- `obj.ExtraConfig.Set(&val)` — set the value
- `obj.ExtraConfig.Unset()` — clear the value
### Partition (Nullable)

Use the Nullable wrapper methods:
- `obj.Partition.IsSet()` — check if set
- `obj.Partition.Get()` — get the inner value (returns pointer)
- `obj.Partition.Set(&val)` — set the value
- `obj.Partition.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


