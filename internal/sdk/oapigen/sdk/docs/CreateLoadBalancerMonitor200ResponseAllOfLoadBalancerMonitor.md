# CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitor

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**LoadBalancer** | Pointer to [**CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorLoadBalancer**](CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorLoadBalancer.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**MonitorType** | Pointer to **string** |  | [optional] 
**MonitorInterval** | Pointer to **int64** |  | [optional] 
**MonitorTimeout** | Pointer to **int64** |  | [optional] 
**SendData** | Pointer to **NullableString** |  | [optional] 
**SendVersion** | Pointer to **string** |  | [optional] 
**SendType** | Pointer to **string** |  | [optional] 
**ReceiveData** | Pointer to **NullableString** |  | [optional] 
**ReceiveCode** | Pointer to **string** |  | [optional] 
**DisabledData** | Pointer to **NullableString** |  | [optional] 
**MonitorUsername** | Pointer to **NullableString** |  | [optional] 
**MonitorPassword** | Pointer to **NullableString** |  | [optional] 
**MonitorPasswordHash** | Pointer to **NullableString** |  | [optional] 
**MonitorDestination** | Pointer to **string** |  | [optional] 
**MonitorReverse** | Pointer to **bool** |  | [optional] 
**MonitorTransparent** | Pointer to **bool** |  | [optional] 
**MonitorAdaptive** | Pointer to **bool** |  | [optional] 
**AliasAddress** | Pointer to **NullableString** |  | [optional] 
**AliasPort** | Pointer to **int64** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**MonitorSource** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**StatusMessage** | Pointer to **NullableString** |  | [optional] 
**StatusDate** | Pointer to **NullableTime** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**MaxRetry** | Pointer to **int64** |  | [optional] 
**FallCount** | Pointer to **int64** |  | [optional] 
**RiseCount** | Pointer to **int64** |  | [optional] 
**DataLength** | Pointer to **NullableInt64** |  | [optional] 
**Config** | Pointer to [**CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfig**](CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorConfig.md) |  | [optional] 
**CreatedBy** | Pointer to [**CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorCreatedBy**](CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitorCreatedBy.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateLoadBalancerMonitor200ResponseAllOfLoadBalancerMonitor{
    // Set fields directly
}
```

### Code (Nullable)

Use the Nullable wrapper methods:
- `obj.Code.IsSet()` — check if set
- `obj.Code.Get()` — get the inner value (returns pointer)
- `obj.Code.Set(&val)` — set the value
- `obj.Code.Unset()` — clear the value
### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### SendData (Nullable)

Use the Nullable wrapper methods:
- `obj.SendData.IsSet()` — check if set
- `obj.SendData.Get()` — get the inner value (returns pointer)
- `obj.SendData.Set(&val)` — set the value
- `obj.SendData.Unset()` — clear the value
### ReceiveData (Nullable)

Use the Nullable wrapper methods:
- `obj.ReceiveData.IsSet()` — check if set
- `obj.ReceiveData.Get()` — get the inner value (returns pointer)
- `obj.ReceiveData.Set(&val)` — set the value
- `obj.ReceiveData.Unset()` — clear the value
### DisabledData (Nullable)

Use the Nullable wrapper methods:
- `obj.DisabledData.IsSet()` — check if set
- `obj.DisabledData.Get()` — get the inner value (returns pointer)
- `obj.DisabledData.Set(&val)` — set the value
- `obj.DisabledData.Unset()` — clear the value
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
### MonitorPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.MonitorPasswordHash.IsSet()` — check if set
- `obj.MonitorPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.MonitorPasswordHash.Set(&val)` — set the value
- `obj.MonitorPasswordHash.Unset()` — clear the value
### AliasAddress (Nullable)

Use the Nullable wrapper methods:
- `obj.AliasAddress.IsSet()` — check if set
- `obj.AliasAddress.Get()` — get the inner value (returns pointer)
- `obj.AliasAddress.Set(&val)` — set the value
- `obj.AliasAddress.Unset()` — clear the value
### StatusMessage (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusMessage.IsSet()` — check if set
- `obj.StatusMessage.Get()` — get the inner value (returns pointer)
- `obj.StatusMessage.Set(&val)` — set the value
- `obj.StatusMessage.Unset()` — clear the value
### StatusDate (Nullable)

Use the Nullable wrapper methods:
- `obj.StatusDate.IsSet()` — check if set
- `obj.StatusDate.Get()` — get the inner value (returns pointer)
- `obj.StatusDate.Set(&val)` — set the value
- `obj.StatusDate.Unset()` — clear the value
### DataLength (Nullable)

Use the Nullable wrapper methods:
- `obj.DataLength.IsSet()` — check if set
- `obj.DataLength.Get()` — get the inner value (returns pointer)
- `obj.DataLength.Set(&val)` — set the value
- `obj.DataLength.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


